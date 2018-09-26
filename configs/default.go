package configs

import (
	"os"
	"strings"
	"time"
)

func getDebugFromEnv() bool {
	debug := false
	debugStr := os.Getenv("AWSUDO_DEBUG")
	if strings.ToLower(debugStr) == "true" {
		debug = true
	}

	return debug
}

var DEBUG bool = getDebugFromEnv()

var DefaultConfPath string = "~/.awsudo/config.toml"
var ReqTimeout time.Duration = 60
var SocketFile string = "/var/tmp/awsudo.sock"
var RetryInterval = []int{1, 3, 5}
var DefaultLogPath string = "/tmp/awsudo.log"
var AppName string = "awsudo"

/*
 * Config File Settings
 */
var DefaultSessionDuration int64 = 3600
var DefaultAgentExpiration int64 = 3600
