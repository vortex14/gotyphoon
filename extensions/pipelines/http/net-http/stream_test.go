package net_http

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
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

func TestFixedSizeOfBuffer(t *testing.T) {

	Convey("test buffer", t, func() {
		data := make([]byte, 2048)
		So(len(data), ShouldEqual, 2048)
		So(cap(data), ShouldEqual, 2048)
	})
}

func TestNotFixedSizeOfBuffer(t *testing.T) {

	Convey("test buffer", t, func() {
		data := make([]byte, 0)
		So(len(data), ShouldEqual, 0)
		So(cap(data), ShouldEqual, 0)

		data = append(data, []byte("test")...)
		So(len(data), ShouldEqual, 4)
		So(cap(data), ShouldEqual, 8)

		data = append(data, []byte("toys-")...)
		So(len(data), ShouldEqual, 9)
		So(cap(data), ShouldEqual, 16)

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

		f, ef := os.Create("stream-1.jpeg")

		So(ef, ShouldBeNil)

		buff := make([]byte, 1024)

		count, dataDump := 0, 0
		finalBuffer := make([]byte, 0)

		for {
			count += 1
			bound, erF := reader.Read(buff)

			dataDump += bound

			if erF == io.EOF {
				LOG.Debug(count, " >> .....", dataDump, "  ", cap(buff))

				break
			}

			finalBuffer = append(finalBuffer, buff[:bound]...)

		}

		_, _ = f.Write(finalBuffer)
		_ = f.Close()
	})
}

func TestHttpChuckToChanel(t *testing.T) {
	Convey("stream download", t, func() {
		task := Fake.CreateDefaultTask()
		//decoder := charmap.Windows1251.NewDecoder()
		task.SetFetcherUrl("https://www.destination-asia.com/ebooks/New/Zip/Indonesia_Ebook.zip")
		task.SetFetcherTimeout(100000)

		request, _ := NewRequest(task)
		_, client := GetHttpClientTransport(task)

		LOG := log.New(map[string]interface{}{"test": "test"})

		response, err := client.Do(request)
		So(err, ShouldBeNil)

		reader := bufio.NewReader(response.Body)

		//So(ef, ShouldBeNil)

		count, dataDump := 0, 0

		ch := make(chan []byte)
		done := make(chan bool)

		allBuff := make([]byte, 0)

		f, _ := os.Create("data.zip")

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(input chan []byte, done chan bool) {
			receivedCount := 1
			for {
				select {
				case data := <-input:
					//println("<<< RECEIVED <<<< ", receivedCount)
					//println(hex.Dump(data.Bytes()))
					//println()
					//println()
					//allBuff = append(allBuff, data...)
					println("Wrote bytes: ", len(data))
					_, _ = f.Write(data)
					receivedCount++
					//println(data)
				case <-done:
					LOG.Debug("EOF received ", len(allBuff))
					//_, _ = f.Write(allBuff)

					wg.Done()
					break
				}
			}
		}(ch, done)

		for {
			buff := make([]byte, 1024*1024*1024)
			count += 1

			bound, erF := reader.Read(buff)

			dataDump += bound

			if erF == io.EOF {
				LOG.Debug(count, " >> .....", dataDump, "  ", cap(buff))
				done <- true
				break
			}

			println("SENDED >>>>> ", count, ": ", len(buff[:bound]))
			ch <- buff[:bound]

		}
		wg.Wait()
		_ = f.Close()

	})
}

func editBuff(data *[]byte) {
	*data = []byte("IT'S NOT MY TEXT")
}

func TestBufferChange(t *testing.T) {

	Convey("test buffer link", t, func() {
		test := []byte("MY TEXT TEST !")

		editBuff(&test)

		So(test, ShouldResemble, []byte("IT'S NOT MY TEXT"))
	})

}
