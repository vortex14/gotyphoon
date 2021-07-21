package ssh

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
)

type SSH struct {
	Ip string
	Login string
	Password string
	client *ssh.Client

}

func (s *SSH) TestConnection() {
	config := &ssh.ClientConfig{
		User: s.Login,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}
	address := fmt.Sprintf("%s:22", s.Ip)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}

	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		color.Red("Failed to create session: ", err.Error())
		os.Exit(1)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("df -h"); err != nil {
		color.Red("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())



}
