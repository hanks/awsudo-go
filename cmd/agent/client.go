package agent

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hanks/awsudo-go/configs"
	"github.com/hanks/awsudo-go/pkg/creds"
	"github.com/hanks/awsudo-go/pkg/parser"
	"github.com/hanks/awsudo-go/utils"
)

type client struct {
	SocketFile string
	Config     *parser.Config
	RoleName   string
}

func (c *client) handleClientFunc(conn net.Conn) {
	defer conn.Close()

	// simple protocol to fetch credentials
	// ask credentials existed or not
	//    if existed, get them directly from server
	//    if not existed, ask user account and pass to fetch credentials and then store into server
	req := "GetCreds"
	_, err := conn.Write([]byte(req))
	if err != nil {
		log.Fatalf("Client error, can not send %s successfully, %v", req, err)
	}

	resp := make([]byte, BUFF_SIZE)
	length, err := conn.Read(resp)
	if err != nil {
		log.Fatalf("Server error, read response error, %v", err)
	}

	resp = resp[:length]
	credentials := new(creds.Creds)

	if string(resp) == NO_CREDS_FLAG {
		user, pass := utils.AskUserInput()
		credentials = creds.FetchCreds(user, pass, c.RoleName, c.Config)
		encoded, _ := credentials.Encode()
		req = fmt.Sprintf("%s%s%s", SET_CREDS_FLAG, DELIMITER, encoded)

		_, err = conn.Write([]byte(req))
		if err != nil {
			log.Fatalf("Client error, can not send %s successfully, %v", req, err)
		}
	} else if strings.HasPrefix(string(resp), "Error") {
		log.Fatalf("Server error, %s, %v", req, err)
	}

	credentials.Decode(resp)
	credentials.SetEnv()
	log.Printf("Set credentials env var ok")
}

func newClient(socket string, c *parser.Config, roleName string) *client {
	return &client{
		SocketFile: socket,
		Config:     c,
		RoleName:   roleName,
	}
}

func (c *client) run() {
	var conn net.Conn
	var err error

	for _, v := range configs.RetryInterval {
		conn, err = net.Dial("unix", c.SocketFile)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(v))
		log.Printf("Retry %d seconds later...", v)
	}

	if err != nil {
		log.Fatalf("Client error, can not connect server successfully, %v", err)
	}

	c.handleClientFunc(conn)
}
