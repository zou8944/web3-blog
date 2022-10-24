package mail

import "sync"

type Mailer struct {
	Driver
}

var once sync.Once
var singleton *Mailer

func DefaultMailer() *Mailer {
	once.Do(func() {
		singleton = &Mailer{
			Driver: NewAWSMailer(),
		}
	})
	return singleton
}

func (s *Mailer) Send(sender, recipient, subject, content string) bool {
	return s.Driver.Send(sender, recipient, subject, content)
}

func (s *Mailer) ReceiveMessageAsString() ([]string, error) {
	return s.Driver.ReceiveMessageAsString()
}
