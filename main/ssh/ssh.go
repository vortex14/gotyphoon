package main

import (
	"fmt"
	sftpAgent "github.com/vortex14/gotyphoon/extensions/agents/sftp"
	"github.com/vortex14/gotyphoon/integrations/ssh"
)

func main()  {
	//req, _ := http.NewRequest("GET", "https://dl.google.com/go/go1.14.2.src.tar.gz", nil)
	//resp, _ := http.DefaultClient.Do(req)
	//defer resp.Body.Close()

	agent := sftpAgent.Constructor(ssh.Options{

	})

	agent.Watch(sftpAgent.Mirror{
		HostPath:        fmt.Sprintf("/Users/vortex/ci-agent"),
		DestinationPath: "/var/test/cli_amd64",
	})



	agent.Await()

	//_ = agent.CopyFileFromHost(sftpAgent.Mirror{
	//	HostPath: "/Users/vortex/ci-agent/cli_amd64",
	//	DestinationPath: "/var/agent/cli_amd64",
	//})


}
