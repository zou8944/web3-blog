package mail

import (
	"blog-web3/pkg/mail/eml"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"time"
)

type BlogMail struct {
	Action   string
	UserName string
	SendFrom string
	Title    string
	Content  string
	Tags     []string
	Visible  string
	Date     time.Time
}

const (
	Public  = "public"
	Private = "private"

	Create = "create"
	Update = "update"
	Delete = "delete"
)

func convert2Blog(e *eml.Email) (*BlogMail, error) {
	var b BlogMail

	if err := errors.WithStack(convertHeader(&b, e)); err != nil {
		return nil, err
	}
	if err := errors.WithStack(convertSubject(&b, e)); err != nil {
		return nil, err
	}
	if err := errors.WithStack(convertBody(&b, e)); err != nil {
		return nil, err
	}

	return &b, nil
}

func convertHeader(d *BlogMail, s *eml.Email) error {
	d.Date = s.Date
	d.SendFrom = s.From.Address
	d.UserName = strings.Split(s.To[0].Address, "@")[0]
	return nil
}

func convertSubject(b *BlogMail, s *eml.Email) error {
	// subject format: [action] #tag1 #tag2 title
	regex := regexp.MustCompile(`^\s*(\[\w+]\s+)*((#\S+\s)+)*([^#]+)\s*$`)
	segments := regex.FindAllStringSubmatch(s.Subject, -1)
	if segments == nil {
		return newFormatError("subject format error")
	}
	actionStr := segments[0][1]
	tagsStr := segments[0][2]
	titleStr := segments[0][4]

	var action string
	var tags []string
	var visible string
	var title string

	// action
	if actionStr == "" {
		action = Create
	} else {
		actionStr = strings.TrimSpace(actionStr)
		action = actionStr[1 : len(actionStr)-1]
		if action != Create && action != Update && action != Delete {
			return newFormatError("action must be one of [create, update, delete]")
		}
	}

	// tags and visible
	tags = strings.Fields(tagsStr)
	visibleIndex := -1
	for i := 0; i < len(tags); i++ {
		tags[i] = strings.TrimSpace(tags[i])
		tags[i] = strings.TrimPrefix(tags[i], "#")
		if tags[i] == Public || tags[i] == Private {
			if visible == "" {
				visible = tags[i]
				visibleIndex = i
			} else {
				return newFormatError("visible tag should be only one")
			}
		}
	}
	if visible == "" {
		visible = Public
	}
	if visibleIndex >= 0 {
		tags = append(tags[:visibleIndex], tags[visibleIndex+1:]...)
	}

	// title
	title = strings.TrimSpace(titleStr)
	if title == "" {
		return newFormatError("title must not be empty")
	}

	b.Title = title
	b.Visible = visible
	b.Tags = tags
	b.Action = action

	return nil
}

func convertBody(b *BlogMail, s *eml.Email) error {
	b.Content = strings.TrimSpace(s.Body)
	if b.Content == "" {
		return newFormatError("body must not be empty")
	}
	return nil
}

var formatErrorType formatError

type formatError struct {
	errString string
}

func newFormatError(errString string) formatError {
	return formatError{errString: errString}
}

func (f formatError) Error() string {
	return f.errString
}
