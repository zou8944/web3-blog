package eml

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseEmail(t *testing.T) {
	illegalEmail, err := os.Open("mail_test/email1.eml")
	if err != nil {
		t.Fatal(err)
	}
	email, err := Parse(illegalEmail)
	if err != nil {
		t.Fatal(err)
	}
	if email.Date.After(time.Now()) {
		t.Fatal("Date parse error")
	}
	if email.From.Address != "zou894475@gmail.com" {
		t.Fatal("Mailer parse error")
	}
	if email.To[0].Address != "guodong@zou8944.com" {
		t.Fatal("Receiver parse error")
	}
	if email.Subject != "Q3总结" {
		t.Fatal("Subject parse error")
	}
	if !strings.Contains(email.Body, "## 完成事项") {
		t.Fatal("Email body parse error")
	}
	if len(email.attachments) == 0 {
		t.Fatal("Email attachment parse error")
	}
}
