package key

import "fmt"

func BadTestTwo() {
	fmt.Println("Nothing here")
}

/*
func GetUserKey() string {
	fmt.Println("What is your private key. Please enter to authenticate: ")

	var key string

	fmt.Scanln(&key)

	return key
}

type Ssh struct {
	Config *ssh.ClientConfig
	Server string
}

func NewClient(user string, host string, port int, privKey string, PrivPass string) (*sshClient, error) {
	pemBytes, err := ioutil.ReadFile(privKey)
	if err != nil {
		return nil, fmt.Errorf("Private key failed.")
	}

	signer, err := signerFromPem(pemBytes, []byte(PrivPass))
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},

		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client := &Ssh{
		Config: config,
		Server: fmt.Sprintf("%v:%v", host, port),
	}

	return client, nil
}

func (s *Ssh) RunSSh(cmd *cobra.Command) (string, error) {

}
*/
