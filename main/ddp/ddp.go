package main

import (
	"github.com/gopackage/ddp"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"time"
)

func init() {
	log.InitD()
}

type Observer struct {
	LOG interfaces.LoggerInterface
}

func (o *Observer) CollectionUpdate(collection, operation, id string, doc ddp.Update) {
	o.LOG.Debug("update collection ", collection, operation, id)
	o.LOG.Debug(doc)
}

func main() {
	// Turn up logging
	logger := log.New(log.D{"ddp": true})

	logger.Info("start !")

	// Assumes Meteor running in development mode normally (no custom port etc.).
	client := ddp.NewClient("ws://localhost:3000/websocket", "http://localhost/")
	defer client.Close()

	// Connect to the server
	err := client.Connect()
	if err != nil {
		logger.Error(">>>>", err.Error())
		return
	}
	logger.Info("connected")

	// Login our user - Meteor.loginWithPassword implements logins using the `login` method. The DDP library
	// provides the implementation for the data the method call expects in ddp.NewUsernameLogin and ddp.NewEmailLogin
	// respectively.
	//login, err := client.Call("login", ddp.NewUsernameLogin(user, pass))
	//if err != nil {
	//	log.WithError(err).Fatal("failed login")
	//} else {
	//	log.WithField("response", login).Info("logged in")
	//}

	// We send a parameter to the `tasks` subscription to demonstrate what that looks like but the tutorial doesn't
	// use this parameter.
	err = client.Sub("links", "abc")
	if err != nil {
		logger.Error("could not subscribe")
		return
	}
	logger = log.Patch(logger, log.D{"version": client.Version(), "session": client.Session()})
	// We know client.Sub is synchronous and will only respond after we connect.

	links := client.CollectionByName("links")
	observer := &Observer{
		LOG: logger,
	}
	links.AddUpdateListener(observer)
	//logger.Info(fmt.Sprintf("collection count :%d", len(tasks.FindAll())))

	time.Sleep(5 * time.Second)

	// Insert a task using a meteor method call
	//log.Info("sending RPC method call to create a task")
	//response, err := client.Call("tasks.insert", "hello " + time.Now().String())
	//if err != nil {
	//	log.WithError(err).Fatal("task create failed")
	//} else {
	//	log.WithField("response", response).Info("created task")
	//}

	// Monitor activity over time. If you create/remove tasks via the Meteor web UI, you will see the tasks collection
	// size change to match.
	for {
		logger.Info("stats", client.Stats())
		logger.Info("count", len(links.FindAll()))

		time.Sleep(10 * time.Second)
	}
}
