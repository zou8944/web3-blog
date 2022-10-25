package mail

type ReceivedMail struct {
	Content  string
	Notifier Notifier
}

// Notifier Notify the email source that the email has been handle successfully
// if the email is from sqs, then corresponding message will be deleted
type Notifier interface {
	Notify()
}

type Driver interface {
	Send(sender, recipient, subject, content string) bool
	Receive() ([]ReceivedMail, error)
}
