package constants

const (
	Network                 = "tcp"
	Port                    = ":8080"
	Delimiter               = '\n'
	NotificationSound       = "\a"
	ServerName              = "server"
	ServerInterruptCall     = 1
	StatusNotFound          = -1
	StatusAcceptable        = "ACCEPTABLE"
	StatusNotAcceptable     = "NOT_ACCEPTABLE"
	StatusInterrupted       = "INTERRUPTED"
	StatusRecipientNotFound = "RECIPIENT_NOT_FOUND"
	MsgNotAcceptable        = "MSG_CORRUPTED\n"
	MsgInterrupted          = "INTERRUPTED\n"
	MsgRecipientNotFound    = "RECIPIENT_NOT_FOUND\n"
)
