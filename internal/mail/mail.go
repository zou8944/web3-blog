package mail

import (
	"github.com/emersion/go-message/mail"
	"github.com/pkg/errors"
	"io"
	"regexp"
	"strings"
	"time"
)

type Email struct {
	Date        time.Time
	From        *mail.Address
	To          []*mail.Address
	Subject     string
	Body        string
	attachments []*mail.Part
}

type BlogEmail struct {
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

func parseEmail(r io.Reader) (*Email, error) {
	m, err := mail.CreateReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "Parse Email fail")
	}

	dat, datErr := m.Header.Date()
	sub, subErr := m.Header.Subject()
	frm, frmErr := m.Header.AddressList("From")
	to, toErr := m.Header.AddressList("To")

	if datErr != nil {
		return nil, errors.Wrap(datErr, "Parse Date fail")
	}
	if subErr != nil {
		return nil, errors.Wrap(subErr, "Parse Subject fail")
	}
	if frmErr != nil {
		return nil, errors.Wrap(frmErr, "Parse From Header fail")
	}
	if toErr != nil {
		return nil, errors.Wrap(toErr, "Parse To Header fail")
	}

	var attachments []*mail.Part
	var body string
	for {
		p, err := m.NextPart()
		if err == io.EOF {
			break
		}
		switch h := p.Header.(type) {
		case *mail.AttachmentHeader:
			attachments = append(attachments, p)
		case *mail.InlineHeader:
			contentType, _, err := h.ContentType()
			if err != nil {
				return nil, errors.Wrap(err, "Parse Content Type fail")
			}
			if contentType == "text/plain" {
				var bytes []byte
				if bytes, err = io.ReadAll(p.Body); err != nil {
					return nil, errors.Wrap(err, "Parse Content Body fail")
				}
				body = string(bytes)
			}
		}
	}

	return &Email{
		Date:        dat,
		From:        frm[0],
		To:          to,
		Subject:     sub,
		Body:        body,
		attachments: attachments,
	}, nil
}

func convert2Blog(e *Email) (*BlogEmail, error) {
	var b BlogEmail

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

func convertHeader(d *BlogEmail, s *Email) error {
	d.Date = s.Date
	d.SendFrom = s.From.Address
	d.UserName = strings.Split(s.To[0].Address, "@")[0]
	return nil
}

func convertSubject(b *BlogEmail, s *Email) error {
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

func convertBody(b *BlogEmail, s *Email) error {
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
