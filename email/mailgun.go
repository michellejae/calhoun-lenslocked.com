package email

import (
	"fmt"

	"gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to Lenslocked.com"
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

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
