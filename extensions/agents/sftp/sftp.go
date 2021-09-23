package sftp

import (
	"github.com/fsnotify/fsnotify"
	. "github.com/vortex14/gotyphoon/extensions/models/watcher"

	"github.com/vortex14/gotyphoon/elements/models/awaitable"
	"github.com/vortex14/gotyphoon/elements/models/label"
	. "github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/integrations/ssh"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type Mirror struct {
	HostPath        string
	DestinationPath string
	ExcludeMatch    []string
}

type Agent struct {
	awaitable.Object
	*label.MetaInfo
	Singleton

	ssh *ssh.SSH
	SSHOptions ssh.Options
	LOG interfaces.LoggerInterface
}

func (a *Agent) init() {
	a.Construct(func() {
		a.ssh = &ssh.SSH{
			Options: ssh.Options{
				Ip: a.SSHOptions.Ip,
				Password: a.SSHOptions.Password,
				Login:a.SSHOptions.Login,
			},
		}
		a.LOG = log.New(log.D{"agent": "SFTP-AGENT"})
		a.LOG.Debug("init agent")
	})

}

func (a *Agent) PingSSH() bool {
	a.init()
	return a.ssh.TestConnection()
}

func (a *Agent) Watch(options Mirror)  {
	a.init()
	if !a.PingSSH() { a.LOG.Error(Errors.ErrorSshConnection.Error()); return }

	w := Watcher{
		ExcludeMatch: options.ExcludeMatch,
		Path: options.HostPath,
		Callback: func(logger interfaces.LoggerInterface, event *fsnotify.Event) {
			logger.Debug(event.Name)

			err := a.CopyFileFromHost(Mirror{HostPath: event.Name, DestinationPath: "/var/test00000100101010.txt"})
			if err != nil {
				return 
			}
		},
	}
	a.Add()
	go w.Watch()
}

func (a *Agent) CopyFileFromHost(options Mirror) error {
	a.init()
	return a.ssh.CopyFileFromHost(options.HostPath, options.DestinationPath)
}