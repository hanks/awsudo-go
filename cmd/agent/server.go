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

	credentials *creds.Creds
}

func newServer(socket string, expire int64) *server {
	s := &server{
		SocketFile:  socket,
		Expiration:  expire,
		sigChan:     make(chan os.Signal, 1),
		stopChan:    make(chan bool, 1),
		connChan:    make(chan net.Conn),
		credentials: new(creds.Creds),
	}

	return s
}

func (s *server) handleServerFunc(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, BUFF_SIZE)
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
		if msgStr == GET_CREDS_FLAG {
			if s.credentials.AccessKeyID == "" {
				log.Println("No stored credentials yet, please set it firstly.")
				conn.Write([]byte(NO_CREDS_FLAG))
			} else {
				log.Println("Credential is already existed, send to client directly.")
				data, err := s.credentials.Encode()
				if err != nil {
					conn.Write([]byte(ENCODE_ERROR))
				} else {
					conn.Write(data)
				}
			}
		}

		if strings.HasPrefix(msgStr, SET_CREDS_FLAG) {
			log.Println("Store credentials to memory for reuse.")
			idx := len(SET_CREDS_FLAG) + len(DELIMITER)
			err := s.credentials.Decode([]byte(msgStr[idx:]))
			if err != nil {
				conn.Write([]byte(DECODE_ERROR))
			}
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
		log.Println("Server is running, just reuse it.")
		return
	}

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
