package agent

// Const for simple protocols use in agent communications
const (
	BuffSize     = 1024
	DELIMITER    = "#"
	GetCredsFlag = "GetCred"
	SetCredsFlag = "SetCred"
	NoCredsFlag  = "NoCred"
	EncodeError  = "Error by encode"
	DecodeError  = "Error by decode"
	BadRequest   = "Bad request"
)
