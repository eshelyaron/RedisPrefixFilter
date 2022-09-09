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
	bfMAddScript    = newScript("bf_madd")
	bfAddScript     = newScript("bf_add")
	bfExistsScript  = newScript("bf_exists")
	bfMExistsScript = newScript("bf_mexists")
	bfReserveScript = newScript("bf_reserve")

	cfAddScript     = newScript("cf_add")
	cfExistsScript  = newScript("cf_exists")
	cfMExistsScript = newScript("cf_mexists")
	cfReserveScript = newScript("cf_reserve")

	pfAddScript     = newScript("pf_add")
	pfMAddScript    = newScript("pf_madd")
	pfExistsScript  = newScript("pf_exists")
	pfMExistsScript = newScript("pf_mexists")
	pfReserveScript = newScript("pf_reserve")
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
