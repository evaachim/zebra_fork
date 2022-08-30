package lease

import "net/smtp"

func (l *Lease) Notify() error {
	from := "admin@zebra.io"
	pwd := "admin123"

	person := l.Status.UsedBy

	to := "<" + person + ">"

	receiver := []string{to}

	host := "smtp.gmail.com"
	port := "587"
	addr := host + ":" + port

	subject := "Lease request fulfilled."
	body := "The lease request for " + l.Type + " for lease ID: " + l.ID + " is ready to be used."

	notification := []byte(subject + body)

	auth := smtp.PlainAuth("", from, pwd, host)

	err := smtp.SendMail(addr, auth, from, receiver, notification)

	if err != nil {
		panic(err)
	}

	return nil
}
