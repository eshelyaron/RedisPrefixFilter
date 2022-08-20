package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"os"
)

const luaDirPath = "scripts"

var (
	//bfInfoScript    = newScript("bf_info")
	bfMAddScript    = newScript("bf_madd")
	bfExistsScript  = newScript("bf_exists")
	bfReserveScript = newScript("bf_reserve")
)

type redisScript struct {
	name   string
	script *redis.Script
}

func newScript(name string) *redisScript {
	path := fmt.Sprintf("%s/%s.lua", luaDirPath, name)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("could not open redis script %v", name)
		return nil
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("could not read redis script %v", name)
		return nil
	}
	script, err := redis.NewScript(string(data)), nil
	if err != nil {
		log.Fatalf("could not load redis scrbf_exists.lua\nbf_info.luaipt %v", name)
		return nil
	}
	return &redisScript{name: name, script: script}
}

func (r *redisScript) exec(ctx context.Context, client redis.UniversalClient, keys []string, args ...interface{}) (interface{}, error) {
	res, err := r.script.Run(ctx, client, keys, args...).Result()
	return res, err
}
