package ssh

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"github.com/vortex14/gotyphoon/elements/models/awaitable"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
	"net"
	"os"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	singleton.Singleton
	awaitable.Object

	Ip string
	Login string
	Password string

	client *ssh.Client
	session *ssh.Session
	sftpClient *sftp.Client

	LOG interfaces.LoggerInterface

}

func (s *SSH) CopyFileFromHost(srcPath string, pathTarget string) error {
	if utils.NotNill(s.client, s.sftpClient) { err, _ := s.CreateNewSFTPClient(); if err != nil { return err} }
	// Open the source file
	srcFile, errO := os.Open(srcPath)
	if errO != nil {
		return errO
	}
	defer srcFile.Close()
	// Create the destination file

	dstFile, err := s.sftpClient.Create(pathTarget)
	if err != nil {
		return err
	}
	defer dstFile.Close()


	fI, _ := srcFile.Stat()


	bar := progressbar.DefaultBytes(
		fI.Size(),
		"uploading",
	)

	proxyReader := progressbar.NewReader(srcFile, bar)
	bar.Reset()
	bar.Add(10000)
	// write to file
	println("Start ")

	if  _, err := dstFile.ReadFrom(proxyReader.Reader); err!= nil {
		return err
	}

	bar.Close()
	//bar.Finish()
	println("End !")
	return nil
}

func (s *SSH) CreateNewSFTPClient() (error, *sftp.Client) {

	if s.client     == nil { err := s.initClient(); if err != nil { return err, nil} }
	if s.sftpClient != nil { return nil, s.sftpClient }

	sftpClient, err := sftp.NewClient(s.client)
	if err != nil {
		return err, nil
	}
	s.sftpClient = sftpClient
	return nil, sftpClient
}

func (s *SSH) initSession() error {
	session, err := s.client.NewSession()
	if err != nil {
		color.Red(Errors.ErrorSshCloseSession.Error(), "  >  ",err.Error())
		return err
	}
	s.session = session
	return nil
}

func (s *SSH) initClient() error {
	var errC error
	s.Construct(func() {
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
			errC = err
		}
		s.client = client
		err = s.initSession()
		if err != nil {
			errC = err
		}
	})

	return errC
}

func (s *SSH) closeSession()  {
	err := s.session.Close()
	if err != nil {
		color.Red(Errors.ErrorSshCloseSession.Error(), "  >  ",err.Error())
	}
}

func (s *SSH) Close()  {
	defer s.closeSession()
	err := s.client.Close()
	if err != nil {
		color.Red(Errors.ErrorSshCloseClient.Error(), "  >  ",err.Error())
	}
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
