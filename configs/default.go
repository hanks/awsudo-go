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

// DEBUG is to print verbose log info with line number
var DEBUG = getDebugFromEnv()

// DefaultConfPath is the path of awsudo config
var DefaultConfPath = "~/.awsudo/config.toml"

// ReqTimeout is the duration of request api timeout
var ReqTimeout time.Duration = 60

// SocketFile is the path of awsudo server socket file
var SocketFile = "/var/tmp/awsudo.sock"

// RetryInterval is the back-off array used in request retry
var RetryInterval = []int{1, 3, 5}

// DefaultLogPath is the log file path of awsudo
var DefaultLogPath = "/tmp/awsudo.log"

// AppName is the name of this app
var AppName = "awsudo"

// MaxSleepUnit is the max value of sleep unit (Millisecond),
// will use random range [min, max] to make the requests not to be sent at the same time
// to trigger the OKta API rate limit error
var MaxSleepUnit = 6000

// MinSleepUnit is the min value of sleep unit (Millisecond)
var MinSleepUnit = 100

/*
 * Config File Settings
 */

// DefaultSessionDuration is session duration of aws assume role seesion
var DefaultSessionDuration int64 = 3600

// DefaultAgentExpiration is expiration duration of awsudo agent server
// DefaultAgentExpiration should be smaller than DefaultSessionDuration, to avoid
// to cache expired credentials
var DefaultAgentExpiration int64 = 3300
