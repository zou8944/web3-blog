package eml

import (
	"github.com/emersion/go-message/mail"
	"github.com/pkg/errors"
	"io"
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

func Parse(r io.Reader) (*Email, error) {
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
