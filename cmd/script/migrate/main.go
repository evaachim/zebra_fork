package main

import (
	"crypto/tls"

	"net/http"

	"github.com/project-safari/zebra/cmd/script/migration"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	migration.PostIt()
}
