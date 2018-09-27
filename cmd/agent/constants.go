package agent

// Const for simple protocols use in agent communications
const (
	BuffSize     = 1024
	DELIMITER    = ":"
	GetCredsFlag = "GetCreds"
	SetCredsFlag = "SetCreds"
	NoCredsFlag  = "NoCreds"
	EncodeError  = "Error by encode"
	DecodeError  = "Error by decode"
)
