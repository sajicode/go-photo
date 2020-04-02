package email

import (
	"fmt"
	"net/url"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to Shutters.com!"
	resetSubject   = "Instructions for resetting your password."
	resetBaseURL   = "https://www.lenslocked.com/reset"
)

const welcomeText = `Hi there!

Welcome to Shutters.com! We really hope you enjoy using
our application!

Best,
Saji
`

const welcomeHTML = `Hi there!<br/>
<br/>
Welcome to
<a href="https://www.github.com/sajicode">Shutters.com</a>! We really hope you enjoy using our application!<br/>
<br/>
Best,<br/>
Saji
`
const resetTextTmpl = `Hi there!

It appears that you have requested a password reset. If this was you, please follow the link below to update your password:

%s

If you are asked for a token, please use the following value:

%s

If you didn't request a password reset you can safely ignore this email and your account will not be changed.

Best,
LensLocked Support
`

const resetHTMLTmpl = `Hi there!<br/>
<br/>
It appears that you have requested a password reset. If this was you, please follow the link below to update your password:<br/>
<br/>
<a href="%s">%s</a><br/>
<br/>
If you are asked for a token, please use the following value:<br/>
<br/>
%s<br/>
<br/>
If you didn't request a password reset you can safely ignore this email and your account will not be changed.<br/>
<br/>
Best,<br/>
Shutters Support<br/>
`

// WithMailgun builds our mailgun credentials
func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		c.mg = mg
	}
}

// WithSender helps us set the sender for our email
func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

// ClientConfig function template
type ClientConfig func(*Client)

// NewClient creates an email client template
func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		// set a default from email address
		from: "support@shutters.com",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

// Client struct for our email
type Client struct {
	from string
	mg   mailgun.Mailgun
}

// Welcome sends the welcome email to users
func (c *Client) Welcome(toName, toEmail string) error {
	message := mailgun.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	message := mailgun.NewMessage(c.from, resetSubject, resetText, toEmail)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, resetUrl, token)
	message.SetHtml(resetHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
