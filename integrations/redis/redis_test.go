package redis

import (
	"strconv"
	"testing"

	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"

	. "github.com/smartystreets/goconvey/convey"
)

var keys []string
var redisService *Service
var matchKey = "test:guest:*"
var testKey = "test"

func init() {
	log.InitD()
	keys = []string{"test:guest:1", "test:guest:2", "test:guest:3"}

	redisService = &Service{
		Config: &interfaces.ServiceRedis{
			Name: "Redis proxy data",
			Details: struct {
				Host     string      `yaml:"host"`
				Port     int         `yaml:"port"`
				Password interface{} `yaml:"password"`
			}(struct {
				Host     string
				Port     int
				Password interface{}
			}{Host: "localhost", Port: 6379}),
		},
	}
	redisService.Init()
}

func TestIncrement(t *testing.T) {
	Convey("check incr", t, func() {
		count := redisService.Incr(testKey)
		So(count > 0, ShouldBeTrue)
	})
}

func TestDecrement(t *testing.T) {
	Convey("set incr", t, func() {
		e := redisService.Set(testKey, 50)
		So(e, ShouldBeNil)
		v := redisService.Get(testKey)
		vint, err := strconv.Atoi(v)
		So(err, ShouldBeNil)
		So(vint, ShouldEqual, 50)
	})

	Convey("decr", t, func() {
		for i := 0; i < 50; i++ {
			_ = redisService.Decr(testKey)
		}

		v := redisService.Get(testKey)

		vint, err := strconv.Atoi(v)
		So(err, ShouldBeNil)
		So(vint, ShouldEqual, 0)
		err = redisService.Remove(testKey)
		So(err, ShouldBeNil)

	})
}

func TestGetList(t *testing.T) {

	Convey("check ping", t, func() {
		So(redisService.Ping(), ShouldBeTrue)
	})

	Convey("set list", t, func() {

		for _, v := range keys {
			e := redisService.Set(v, true)
			So(e, ShouldBeNil)
		}

		So(redisService.Count(matchKey), ShouldEqual, 3)

	})

	Convey("get list", t, func() {
		arr := redisService.GetList(matchKey)
		So(len(arr), ShouldEqual, 3)
	})

	Convey("remove list", t, func() {
		for _, v := range keys {
			e := redisService.Remove(v)
			So(e, ShouldBeNil)
		}
		So(redisService.Count(matchKey), ShouldEqual, 0)
	})
}
