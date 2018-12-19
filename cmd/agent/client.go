package agent

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
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

func (c *client) buildReq(cmds []string) string {
	return strings.Join(cmds[:], DELIMITER)
}

func (c *client) handleClientFunc(conn net.Conn) {
	defer conn.Close()

	// simple protocol to fetch credentials
	// ask credentials existed or not
	//    if existed, get them directly from server
	//    if not existed, ask user account and pass to fetch credentials and then store into server
	//       here we will fetch credentials for all aws accounts at the first time for future usages.

	_, roleARN := c.Config.GetARN(c.RoleName)
	if roleARN == "" {
		log.Fatalf("Can not find role name: %v", c.RoleName)
	}

	// GetCred#roleARN
	req := c.buildReq([]string{GetCredsFlag, roleARN})
	_, err := conn.Write([]byte(req))
	if err != nil {
		log.Fatalf("Client error, can not send %s successfully, %v", req, err)
	}

	resp := make([]byte, BuffSize)
	length, err := conn.Read(resp)
	if err != nil {
		log.Fatalf("Server error, read response error, %v", err)
	}

	resp = resp[:length]
	cred := new(creds.Creds)

	if string(resp) == NoCredsFlag {
		scanner := bufio.NewScanner(os.Stdin)
		user, pass := utils.AskUserInput(scanner)

		c.fetchCredentialsByGroup(conn, cred, user, pass)
	} else if strings.HasPrefix(string(resp), "Error") {
		log.Fatalf("Server error, %s, %v", req, err)
	} else {
		// set credential to env var from cache
		cred.Decode(resp)
		cred.SetEnv()
	}

	log.Printf("Set credentials env var ok")
}

func (c *client) fetchCredentialsByGroup(conn net.Conn, cred *creds.Creds, user string, pass string) {
	var wg sync.WaitGroup

	for _, role := range c.Config.Roles {
		wg.Add(1)
		go func(roleName string, roleARN string) {
			defer wg.Done()

			// to simulate random.range(min, max)
			sleepUnit := rand.Intn(configs.MaxSleepUnit-configs.MinSleepUnit) + configs.MinSleepUnit
			log.Printf("Role %s's random sleep unit for api request is %d", roleName, sleepUnit)
			time.Sleep(time.Duration(sleepUnit) * time.Millisecond)

			cred = creds.FetchCreds(user, pass, roleName, c.Config)
			encoded, _ := cred.Encode()

			// SetCred#roleARN#value
			req := c.buildReq([]string{SetCredsFlag, roleARN, string(encoded)})
			_, err := conn.Write([]byte(req))
			if err != nil {
				log.Fatalf("Client error, can not send %s successfully, %v", req, err)
			}

			if c.RoleName == roleName {
				// set credential to env var from newly created cred instance if it is the user's request
				cred.SetEnv()
			}
		}(role.Name, role.ARN)
	}

	wg.Wait()
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

	// wait for server running
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
