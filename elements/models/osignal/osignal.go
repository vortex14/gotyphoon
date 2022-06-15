package osignal

import (
	"os"
	"os/signal"

	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type OSignal struct {
	singleton.Singleton
	buff chan os.Signal
	LOG interfaces.LoggerInterface
	Callback func(logger interfaces.LoggerInterface, sig os.Signal)
}

func (s *OSignal) Wait() {
	s.LOG = log.New(log.D{"signal": "syscall.SIGINT"})
	s.buff = make(chan os.Signal, 1)

	signal.Notify(s.buff, os.Interrupt)
	s.LOG.Debug("waiting for CTRL+C ...")
	sig := <-s.buff

	s.Destruct(func() {
		s.Callback(s.LOG, sig)
	})
}
