package redis

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Service struct {
	Project interfaces.Project
	client *redis.Client
	Config *interfaces.ServiceRedis
	//*interfaces.ServiceRedis
}

func (s *Service) initClient()  {
	if s.client == nil {
		redisString := fmt.Sprintf("%s:%d", s.Config.GetHost(), s.Config.GetPort())
		//color.Yellow("init Redis Service %s", redisString)
		rdb := redis.NewClient(&redis.Options{
			Addr:     redisString, // use default Addr
			Password: "",               // no password set
			DB:       0,                // use default DB
		})
		s.client = rdb
		conn := s.connect()
		if conn {
			color.Green("Redis client success %s", redisString)
		}
	}

}

func (s *Service) connect() bool {
	s.initClient()
	var ctx = context.Background()
	defer ctx.Done()
	status := false
	_, err := s.client.Ping(ctx).Result()
	if err != nil {
		color.Red("%s", err)
	} else {
		status = true
	}
	return status
}

func (s *Service) Set(key string, value interface{}) error {
	var ctx = context.Background()
	defer ctx.Done()
	status := s.client.Set(ctx, key, value, 0)
	return status.Err()
}

func (s *Service) Init()  {
	if !s.Ping() {
		color.Red("Redis connection failed. %s:%s", s.Config.GetHost(), s.Config.GetPort())
	}
}

func (s *Service) Ping() bool {
	return s.connect()
}