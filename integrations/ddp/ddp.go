package ddp

import (
	"github.com/gopackage/ddp"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type Collection struct {
	*label.MetaInfo
	Fn func(collection, operation, id string, doc ddp.Update)
}

type DDP struct {
	singleton.Singleton

	client       *ddp.Client
	EndpointWS   string // ws://localhost:3000/websocket
	EndpointHTTP string // http://localhost/
	LOG          interfaces.LoggerInterface

	Collections []*Collection
}

func (d *DDP) init() {
	d.Construct(func() {
		d.LOG = log.New(log.D{"connection": "ddp", "ws": d.EndpointWS, "http": d.EndpointHTTP})
		d.LOG.Debug("init")
		d.client = ddp.NewClient(d.EndpointWS, d.EndpointHTTP)

		err := d.client.Connect()
		if err != nil {
			d.LOG.Error(">>>>", err.Error())
			d.client = nil
			return
		}
	})
}

func (d *DDP) Connect() {
	if d.init(); d.client != nil {
		d.LOG.Debug("connected !")
	}

}

func (d *DDP) Close() {
	d.client.Close()
}
