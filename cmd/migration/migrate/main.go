// nolint: gosec, forcetypeassert // Using this script for a secure server with a secure jwt token.
package main

import (
	"crypto/tls"
	"net/http"

	"github.com/project-safari/zebra/cmd/migration"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	migration.PostIt()
}
