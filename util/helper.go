package util

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func SendMail(emailBody string) {
	// Set up authentication information.

	smtpServer := "smtp.mailgun.org"
	auth := smtp.PlainAuth(
		"",
		"postmaster@app19530200.mailgun.org",
		"44xf4meg4du1",
		smtpServer,
	)

	from := mail.Address{"postmaster", "postmaster@app19530200.mailgun.org"}
	to := mail.Address{"Jasdeep", "jasdeepm@gmail.com"}
	title := "test email jasdeep"

	body := emailBody

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer+":25",
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	//[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func SendCustomMail(emailBody, fromUser, toUser, emailTitle string) {
	// Set up authentication information.

	smtpServer := "smtp.mailgun.org"
	auth := smtp.PlainAuth(
		"",
		"postmaster@app19530200.mailgun.org",
		"44xf4meg4du1",
		smtpServer,
	)

	from := mail.Address{"", fromUser}
	to := mail.Address{"", toUser}
	title := emailTitle

	body := emailBody

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer+":25",
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	//[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func GenUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// TODO: verify the two lines implement RFC 4122 correctly
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid), nil
}

func ReadFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
