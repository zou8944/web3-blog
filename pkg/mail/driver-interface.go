package mail

type Driver interface {
	Send(sender, recipient, subject, content string) bool
	ReceiveMessageAsString() ([]string, error)
}
