package lease

import "net/smtp"

func (l *Lease) Notify() error {
	// may need to create a user that manages the system.
	// would need a password and email address for that person.
	// that would probably need to be someone who manages the tool rather than a regular admin.
	// then use that for sending notifications.
	from := " admin@zebra.project-safari.io"
	pwd := "admin123"

	person := l.Status.UsedBy

	receiver := []string{
		person,
	}

	host := "smtp.gmail.com"
	port := "587"
	addr := host + ":" + port

	subject := "Lease request fulfilled. "
	body := "\nThe lease request for " + l.Type + " with lease ID: " + l.ID + " is ready to be used."

	notification := []byte(subject + body)

	auth := smtp.PlainAuth("", from, pwd, host)

	err := smtp.SendMail(addr, auth, from, receiver, notification)
	if err != nil {
		panic(err)
	}

	return nil
}
