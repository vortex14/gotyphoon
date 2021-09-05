package forms

import (
	"context"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Controller func(ctx context.Context, logger interfaces.LoggerInterface)
