package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/pkg/errors"
)

type (
	tEmail struct {
		From    string
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	tSmtpLoginAuth struct {
		username string
		password string
	}
)

func smtpSend(email *tEmail) error {
	fmt.Println("connecting to SMTP: ", fmt.Sprintf("%s:%d", config.Server, config.Port))
	var (
		conn net.Conn
		err  error
	)
	// Set TLS version
	requiredTLSmin, err := setTLSversion(config.TLSmin)
	if err != nil {
		return fmt.Errorf("SET TLSmin: %s", err.Error())
	}
	requiredTLSmax, err := setTLSversion(config.TLSmax)
	if err != nil {
		return fmt.Errorf("SET TLSman: %s", err.Error())
	}
	// Open server connection
	switch config.TLS {
	case "", "StartTLS":
		conn, err = net.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", config.Server, config.Port),
		)
	case "TLS":
		conn, err = tls.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", config.Server, config.Port),
			&tls.Config{
				MinVersion:         requiredTLSmin,
				MaxVersion:         requiredTLSmax,
				InsecureSkipVerify: !config.TLSvalid,
			},
		)
	default:
		return fmt.Errorf("SMTP error: unsupported connection type")
	}
	if err != nil {
		return fmt.Errorf("SMTP error: %s", err.Error())
	}
	defer conn.Close()

	// Start SMTP session
	hostname, _ := os.Hostname()
	c, err := smtp.NewClient(conn, hostname)
	if err != nil {
		return fmt.Errorf("SMTP connection error: %s", err.Error())
	}
	defer c.Close()

	if config.TLS == "StartTLS" {
		if err := c.StartTLS(&tls.Config{
			MinVersion:         requiredTLSmin,
			MaxVersion:         requiredTLSmax,
			InsecureSkipVerify: !config.TLSvalid,
			ServerName:         config.Server,
		}); err != nil {
			return fmt.Errorf("SMTP c.StartTLS error: %s", err.Error())
		}
	}

	switch config.Auth {
	case "PLAIN":
		err = c.Auth(smtp.PlainAuth(
			"",
			config.User,
			config.Password,
			"",
		))
	case "LOGIN":
		err = c.Auth(smtpLoginAuth(
			config.User,
			config.Password,
		))
	}
	if err != nil {
		return fmt.Errorf("SMTP config.Auth error: %s", err.Error())
	}

	if err = c.Mail(smtpFormatPlainEmailAddress(email.From)); err != nil {
		return fmt.Errorf("SMTP c.Mail error: %s", err.Error())
	}
	if err = c.Rcpt(smtpFormatPlainEmailAddress(email.To)); err != nil {
		return fmt.Errorf("SMTP c.Rcpt error: %s", err.Error())
	}
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return fmt.Errorf("SMTP data error: %s", err.Error())
	}
	if _, err = fmt.Fprintf(wc, "%s", smtpComposeMimeMail(email)); err != nil {
		return fmt.Errorf("SMTP data write error: %s", err.Error())
	}
	if err = wc.Close(); err != nil {
		return fmt.Errorf("SMTP close error: %s", err.Error())
	}
	// Send the QUIT command and close the connection.
	if err = c.Quit(); err != nil {
		return fmt.Errorf("SMTP finish error: %s", err.Error())
	}
	return nil
}

func smtpFormatEmailAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}
	return e.String()
}

func smtpFormatPlainEmailAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}
	return e.Address
}

func smtpComposeMimeMail(email *tEmail) []byte {
	header := make(map[string]string)
	header["From"] = smtpFormatEmailAddress(email.From)
	header["To"] = smtpFormatEmailAddress(email.To)
	header["Subject"] = fmt.Sprintf("=?utf-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(email.Subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(email.Body))

	return []byte(message)
}

// Login Auth
func smtpLoginAuth(username, password string) smtp.Auth {
	return &tSmtpLoginAuth{username, password}
}

func (a *tSmtpLoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *tSmtpLoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func setTLSversion(ver string) (uint16, error) {
	switch ver {
	case "1.3":
		return tls.VersionTLS13, nil
	case "1.2":
		return tls.VersionTLS12, nil
	case "1.1":
		return tls.VersionTLS11, nil
	case "1.0":
		return tls.VersionTLS10, nil
	case "SSL":
		return tls.VersionSSL30, nil
	default:
		return 0, fmt.Errorf("unknown TLS version: %s", ver)
	}
}
