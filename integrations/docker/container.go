package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Container struct {
	Client *client.Client
	Link *types.Container
}

func (c *Container) Exec(ctx context.Context, commands []string) (error, ExecResult) {
	exec, err := Exec(ctx, c.Client, c.Link.ID, commands)
	if err != nil { return err, ExecResult{} }

	return nil, exec
}
