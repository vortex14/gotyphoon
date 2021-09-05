package http

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/data/fake"
	net_http "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"testing"
)

func TestHTTPRequest(t *testing.T) {

	Convey("create a new task for request", t, func() {
		task := fake.CreateDefaultTask()
		Convey("set url for request", func() {

			task.SetFetcherUrl("http://localhost:12666/fake/image")

			Convey("fetch data", func() {

				err, data := net_http.FetchData(task)


					Convey("checking data", func() {

						So(err, ShouldBeNil)

						So(data, ShouldNotBeEmpty)

					})
			})
		})

	})
}
