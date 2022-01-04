package email

import (
	"fmt"
	"net/url"

	"gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to Lenslocked.com"
	resetSubject   = "Instructions for resetting your password!"
	resetBaseURL   = "https://www.lenslocked.com"
)

const welcomeText = `Hi there, ho there, hey there!

Welcome to this super awesome site, lenslocked.com! I hope you have a great time
exploring!

Best,
m
`

const welcomeHTML = `Hi there, ho there, hey there!<br/>
<br/>
Welcome to this super awesome site, <a href="https://www.lenslocked.com">LensLocked.com</a>!<br/>
<br/>I hope you have a great time
exploring!
<br/>
Best,
<br/>
m`

const resetTextTmpl = `Hi there!

It appears you have requested a password reset! If this was you, please follow the link below to update your password:

%s

If you are asked for a token please use the following value:

%s

If you didn't request a password reset you can safely ignore this message and your account will be unchanged!

Best,
Da support team!
`

const resetHTMLTmpl = `Hi there!<br/>
<br/>
It appears you have requested a password reset! If this was you, please follow the link below to update your password:<br/>
<br/>
<a href="%s">%s</a><br/>
<br/>
If you are asked for a token please use the following value:<br/>
<br/>
%s<br/>
<br/>
If you didn't request a password reset you can safely ignore this message and your account will be unchanged!<br/>
<br/>
Best,<br/>
Da support team!<br/>
`

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		c.mg = mg
	}
}

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		from: "support@lenslocked.com",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

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
