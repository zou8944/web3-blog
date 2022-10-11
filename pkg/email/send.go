package email

type ResponseTemplate struct {
	Subject string
	Body    string
}

var (
	BadFormatTemplate = ResponseTemplate{
		Subject: "Bad format email",
		Body: `
The format of your email is incorrect, the correct format should be:

Subject: [action] #tag1 #tag2 #tagn title
	action should be one of [create, update, delete]
	there must be at least one tag: #private or #public, it will be treated as visibility. other tag will be treated as real tag.
	title will be treated as your blog title.
Body: raw markdown text, or plain text.
	body will be treated as your blog content.
`,
	}
)

func SendEmailTemplate(address string, template ResponseTemplate) {

}

func SendEmail(address, subject, content string) {

}
