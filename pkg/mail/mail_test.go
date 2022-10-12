package mail

import (
	"github.com/emersion/go-message/mail"
	"github.com/pkg/errors"
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
	email, err := parseEmail(illegalEmail)
	if err != nil {
		t.Fatal(err)
	}
	if email.Date.After(time.Now()) {
		t.Fatal("Date parse error")
	}
	if email.From.Address != "zou894475@gmail.com" {
		t.Fatal("Sender parse error")
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

func TestConvertHeader(t *testing.T) {
	var d BlogEmail
	var s Email

	s.Date = time.Now()
	s.From = &mail.Address{
		Address: "gd@from.com",
	}
	s.To = []*mail.Address{
		{
			Address: "gd@to.com",
		},
	}
	if err := convertHeader(&d, &s); err != nil {
		t.Fatal(err)
	}
	if d.Date != s.Date {
		t.Fatal("date convert error")
	}
	if d.SendFrom != s.From.Address {
		t.Fatal("from address convert error")
	}
	if d.UserName != "gd" {
		t.Fatal("to address convert error")
	}
}

func TestParseSubject(t *testing.T) {
	var d BlogEmail
	var s Email
	// no action
	s.Subject = "#private title"
	if err := convertSubject(&d, &s); err != nil {
		t.Fatal(err)
	}
	if d.Action != Create {
		t.Fatal("default action should be create")
	}
	// action not in [create, update, delete]
	s.Subject = "[query] #private title"
	if err := convertSubject(&d, &s); err == nil {
		t.Fatal("should return error for 'query'")
	} else if _, ok := errors.Cause(err).(formatError); !ok {
		t.Fatal("should return format error for 'query'")
	}
	// no visible tag
	s.Subject = "title"
	if err := convertSubject(&d, &s); err != nil {
		t.Fatal(err)
	}
	// multiple visible tag
	s.Subject = "#private #public title"
	if err := convertSubject(&d, &s); err == nil {
		t.Fatal("should return error for multiple visible tag")
	} else if _, ok := errors.Cause(err).(formatError); !ok {
		t.Fatal("should return format error for  multiple visible tag")
	}
	// no title
	s.Subject = "#private "
	if err := convertSubject(&d, &s); err == nil {
		t.Fatal("should return error for empty title")
	} else if _, ok := errors.Cause(err).(formatError); !ok {
		t.Fatal("should return format error for empty title")
	}
	// normal case
	s.Subject = "title"
	if err := convertSubject(&d, &s); err != nil {
		t.Fatal(err)
	}
	s.Subject = "[create] #private #hot title"
	if err := convertSubject(&d, &s); err != nil {
		t.Fatal(err)
	}
}

func TestParseBody(t *testing.T) {
	var d BlogEmail
	var s Email

	// body too short
	s.Body = ""
	err := convertBody(&d, &s)

	if err == nil {
		t.Fatal("body too short should be reject")
	}
	if _, ok := errors.Cause(err).(formatError); !ok {
		t.Fatal("body too short should be reject")
	}
	// normal case
	s.Body = "this is a formal content"
	if err := convertBody(&d, &s); err != nil {
		t.Fatal(err)
	}
}
