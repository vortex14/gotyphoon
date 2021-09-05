package interfaces

import (
	"context"
)

type Controller func(context context.Context, logger LoggerInterface)