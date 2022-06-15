package sftp

import "github.com/vortex14/gotyphoon/integrations/ssh"

func Constructor(options ssh.Options) *Agent {

	return &Agent{
		SSHOptions: options,
	}

}
