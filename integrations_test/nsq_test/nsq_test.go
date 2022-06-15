package nsq_test

import (
	"github.com/vortex14/gotyphoon/integrations/nsq"
	"github.com/vortex14/gotyphoon/interfaces"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	typhoon "github.com/vortex14/gotyphoon"
)

func TestMessagePubV11(t *testing.T) {
	Convey("connect to NSQ", t, func() {
		pathProject, _ := os.Getwd()
		project := &typhoon.Project{
			ConfigFile: "config-versions/v1.1/config.local.yaml",
			Path:       pathProject,
		}

		nsqService := nsq.Service{Project: project}
		status := nsqService.Ping()
		So(status, ShouldBeTrue)

	})

	testMessage := "{\"test\":{\"m\":1}}"

	Convey("Pub message to NSQ", t, func() {
		pathProject, _ := os.Getwd()
		project := &typhoon.Project{
			ConfigFile: "config-versions/v1.1/config.local.yaml",
			Path:       pathProject,
		}

		nsqService := nsq.Service{Project: project, Options: interfaces.MessageBrokerOptions{EnabledProducer: true, Active: true}}

		testQueue := &interfaces.Queue{
			Channel:  "test-channel",
			Topic:    "test-topic",
			Writable: true,
		}

		testQueue.SetGroupName("test-group")

		nsqService.InitQueue(testQueue)
		err := nsqService.Pub("test-group", "test-topic", testMessage)
		So(err, ShouldBeNil)
	})

	Convey("Test NSQ reader", t, func() {
		pathProject, _ := os.Getwd()
		project := &typhoon.Project{
			ConfigFile: "config-versions/v1.1/config.local.yaml",
			Path:       pathProject,
		}

		nsqService := nsq.Service{
			Project: project,
			Options: interfaces.MessageBrokerOptions{EnabledProducer: true, Active: true},
		}

		status := nsqService.Ping()

		So(status, ShouldBeTrue)

		testQueue := &interfaces.Queue{
			Channel:  "test-channel",
			Topic:    "test-topic",
			Readable: true,
		}

		testQueue.SetGroupName("test-group")

		nsqService.InitQueue(testQueue)

		consumer := nsqService.InitConsumer(testQueue)

		statusR := make(chan bool, 1)
		MessageD := make(chan []byte, 1)

		go func() {
			println("Run NSQ reader")
			var count int
			count = 0
			for msg := range consumer.Messages() {
				msg.Finish()
				count += 1
				statusR <- true
				MessageD <- msg.Body
				println("received a new message")
				break
			}
		}()

		Convey("Awaiting for test message", func() {
			<-statusR

			data := <-MessageD

			So(string(data), ShouldEqual, testMessage)

			consumer.Stop()
		})

	})

}
