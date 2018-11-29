package agent

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hanks/awsudo-go/pkg/creds"
)

type server struct {
	SocketFile string
	Expiration int64
	listener   net.Listener

	sigChan  chan os.Signal
	stopChan chan bool
	connChan chan net.Conn

	credentials map[string]*creds.Creds
}

func newServer(socket string, expire int64) *server {
	s := &server{
		SocketFile:  socket,
		Expiration:  expire,
		sigChan:     make(chan os.Signal, 1),
		stopChan:    make(chan bool, 1),
		connChan:    make(chan net.Conn),
		credentials: make(map[string]*creds.Creds),
	}

	return s
}

func (s *server) validate(msg string) bool {
	if len(msg) == 0 {
		return false
	}

	cmds := strings.Split(msg, DELIMITER)
	if len(cmds) <= 1 {
		return false
	}

	cmd := cmds[0]
	if cmd != GetCredsFlag && cmd != SetCredsFlag {
		return false
	}
	// [GetCred roleARN]
	if cmd == GetCredsFlag && len(cmds) != 2 {
		return false
	}
	// [SetCred roleARN cred]
	if cmd == SetCredsFlag && len(cmds) != 3 {
		return false
	}

	return true
}

func (s *server) handleServerFunc(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, BuffSize)
		length, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Server error, can not read buffer successfully, %v", err)
				log.Println("Please try again.")
				return
			}

			log.Printf("No more data is coming, conn %v is closed.", conn)
			return
		}
		buf = buf[:length]

		msgStr := string(buf)
		if !s.validate(msgStr) {
			conn.Write([]byte(BadRequest))
		}

		cmds := strings.Split(msgStr, DELIMITER)
		cmd := cmds[0]

		if cmd == GetCredsFlag {
			key := cmds[1]
			_, exist := s.credentials[key]
			if exist {
				log.Println("Credential is already existed, send to client directly.")
				data, err := s.credentials[key].Encode()
				if err != nil {
					conn.Write([]byte(EncodeError))
				} else {
					conn.Write(data)
				}
			} else {
				log.Println("No stored credentials yet, please set it firstly.")
				conn.Write([]byte(NoCredsFlag))
			}
		}

		if cmd == SetCredsFlag {
			log.Println("Store credentials to memory for reuse.")
			key := cmds[1]
			value := cmds[2]

			cred, err := creds.NewCred([]byte(value))
			if err != nil {
				conn.Write([]byte(DecodeError))
			}
			s.credentials[key] = cred
		}
	}
}

func (s *server) stop() {
	<-s.stopChan
	s.stopChan <- true

	s.listener.Close()
	os.Remove(s.SocketFile)

	os.Exit(0)
}

func (s *server) terminate() {
	signal := <-s.sigChan
	log.Printf("Caught signal %s: shutting down...", signal)
	// send loop stop flag to main loop
	s.stop()
}

func (s *server) accept() {
	for {
		s.stopChan <- false
		conn, err := s.listener.Accept()

		stop := <-s.stopChan
		if stop {
			os.Exit(1)
		}

		if err != nil {
			log.Fatalf("Server error, can not accept socket successfully, %v", err)
		}

		s.connChan <- conn
	}
}

func (s *server) run() {
	// skip when server is running
	if _, err := os.Stat(s.SocketFile); !os.IsNotExist(err) {
		log.Println("Server is already running, just reuse it.")
		return
	}
	log.Printf("Start agent server to handle new request, and will be expired after %d seconds", s.Expiration)

	log.Printf("Listen to socket file: %s", s.SocketFile)
	listener, err := net.Listen("unix", s.SocketFile)
	if err != nil {
		log.Fatalf("Server error, can not listen socket successfully, %v", err)
	}
	s.listener = listener

	signal.Notify(s.sigChan, syscall.SIGINT, syscall.SIGTERM)
	go s.terminate()
	go s.accept()

	for {
		select {
		case conn := <-s.connChan:
			log.Printf("Handle new conn, %v", conn)
			go s.handleServerFunc(conn)
		case <-time.After(time.Second * time.Duration(s.Expiration)):
			log.Printf("Agent server is expired after %d seconds. Please start it again.", s.Expiration)
			s.stop()
		}
	}
}
