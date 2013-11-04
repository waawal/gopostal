package main

import (
	"github.com/darkhelmet/env"
	"github.com/garyburd/redigo/redis"
	"github.com/hoisie/web"
	"io/ioutil"
)

var r = RedisConnection()

func main() {
	host := env.StringDefault("GOPOSTAL_HOST", "0.0.0.0:9999")
	web.Post("/(.*)", publisher)
	web.Run(host)
}

func publisher(ctx *web.Context, queue string) string {
	ctx.SetHeader("Server", "go.postal", true)
	body, err := ioutil.ReadAll(ctx.Request.Body)
	check(err)
	ctx.Request.Body.Close()
	go publish(queue, string(body))
	return "OK"

}

func publish(queue string, data string) {
	_, err := r.Do("PUBLISH", queue, data)
	check(err)
}

func RedisConnection() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	check(err)
	// defer c.Close()
	return c
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
