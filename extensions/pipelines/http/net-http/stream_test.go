package net_http

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	Fake "github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

func TestHttpFullBody(t *testing.T) {

	Convey("basic download", t, func() {
		task := Fake.CreateDefaultTask()
		//decoder := charmap.Windows1251.NewDecoder()
		task.SetFetcherUrl("https://4kwallpapers.com/images/walls/thumbs_3t/10017.png")

		eF, body := FetchData(task)
		//enc, e := decoder.String(*body)
		//if e != nil {
		//	panic(e)
		//}
		So(eF, ShouldBeNil)

		e := ioutil.WriteFile("stream.png", []byte(*body), os.ModePerm)

		So(e, ShouldBeNil)
	})

}

func TestReader(t *testing.T) {
	reader := strings.NewReader("Clear is better than clever")
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Println(string(p[:n]))
	}

}

func TestBufferSize(t *testing.T) {

	Convey("test buffer", t, func() {
		data := make([]byte, 2048)
		So(len(data), ShouldEqual, 2048)
		So(cap(data), ShouldEqual, 2048)
	})
}

func TestHttpFixedBuffer(t *testing.T) {

	Convey("stream download", t, func() {
		task := Fake.CreateDefaultTask()
		//decoder := charmap.Windows1251.NewDecoder()
		task.SetFetcherUrl("https://4kwallpapers.com/images/walls/thumbs_3t/10017.png")
		request, _ := NewRequest(task)
		_, client := GetHttpClientTransport(task)

		LOG := log.New(map[string]interface{}{"test": "test"})

		response, err := client.Do(request)
		So(err, ShouldBeNil)

		reader := bufio.NewReader(response.Body)

		f, ef := os.Create("stream.png")

		So(ef, ShouldBeNil)

		buff := make([]byte, 4096)
		count, dataDump := 0, 0

		for {
			count += 1
			bound, erF := reader.Read(buff)

			dataDump += bound
			//enc, e := decoder.String(string(buff))

			if erF == io.EOF {
				_ = f.Close()
				break
			}

			LOG.Debug(count, ">>> LEN: ", len(buff), "CAP: ", cap(buff), "   BO : ", bound, " SUM: ", dataDump)
			_, _ = f.Write(buff[:bound])

		}

	})

}

func TestHttpAccumulateBuffer(t *testing.T) {
	Convey("stream download", t, func() {
		task := Fake.CreateDefaultTask()
		//decoder := charmap.Windows1251.NewDecoder()
		task.SetFetcherUrl("https://4kwallpapers.com/images/walls/thumbs_3t/10017.png")
		request, _ := NewRequest(task)
		_, client := GetHttpClientTransport(task)

		LOG := log.New(map[string]interface{}{"test": "test"})

		response, err := client.Do(request)
		So(err, ShouldBeNil)

		reader := bufio.NewReader(response.Body)

		f, ef := os.Create("stream-2.jpeg")

		So(ef, ShouldBeNil)

		buff := make([]byte, 1024)

		count, dataDump := 0, 0
		finalBuffer := make([]byte, 0)
		for {
			count += 1
			bound, erF := reader.Read(buff)

			dataDump += len(buff)

			if erF == io.EOF {
				LOG.Debug(count, " >> .....", dataDump, "  ", cap(buff))
				_ = f.Close()
				break
			}

			finalBuffer = append(finalBuffer, buff[:bound]...)

		}

		_, _ = f.Write(finalBuffer)
	})
}
