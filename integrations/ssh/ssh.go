package ssh

import (
	"bytes"
	"fmt"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"net"
	"os"

	"github.com/fatih/color"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"

	"github.com/vortex14/gotyphoon/elements/models/awaitable"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
	progressFile "github.com/vortex14/gotyphoon/extensions/models/progress-file"
)

func init()  {
	log.InitD()
}

type Options struct {
	Ip string
	Login string
	Password string
}

type SSH struct {
	singleton.Singleton
	awaitable.Object
	Options

	client *ssh.Client
	session *ssh.Session
	sftpClient *sftp.Client

	LOG interfaces.LoggerInterface

}

func (s *SSH) CopyFileFromHost(srcPath string, pathTarget string) error {
	if utils.NotNill(s.client, s.sftpClient) { err, _ := s.CreateNewSFTPClient(); if err != nil { return err} }

	logger := log.New(log.D{"name": "uploaderFile"})

	sftpUploadFile := &progressFile.File{
		MetaInfo: &label.MetaInfo{
			Description: fmt.Sprintf("copying %s", srcPath),
		},
		Path: srcPath,
		OnFinish: func(f *os.File) {
			logger.Debug("test done !", f.Name())
		},
	}


	dstFile, err := s.sftpClient.Create(pathTarget)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if  _, err := dstFile.ReadFrom(sftpUploadFile); err!= nil {
		return err
	}

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

func (s *SSH) TestConnection() bool {
	var status bool
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
		return status
	}

	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		color.Red("Failed to create session: ", err.Error())
		return status
	}
	defer session.Close()



	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run("df -h"); err != nil {
		color.Red("Failed to run: " + err.Error())
		return status
	}
	//fmt.Println(b.String())
	status = true
	return status

}
