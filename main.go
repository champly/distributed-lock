package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func createScript() *redis.Script {
	script := redis.NewScript(` local key = tostring(KEYS[1])
		local value = tostring(KEYS[2])
		local expireTime = tonumber(ARGV[1])

		local result = redis.call("setnx", key, value)
		if result == 0 then
			return 0
		end

		redis.call("expire", key, expireTime)
		return 1
	`)

	return script
}

func createScriptDel() *redis.Script {
	script := redis.NewScript(`
		local key = tostring(KEYS[1])

		return redis.call("del", key)
	`)
	return script
}

func scriptCacheToCluster(c *redis.ClusterClient, script *redis.Script) string {
	var ret string

	c.ForEachMaster(func(m *redis.Client) error {
		if result, err := script.Load(m).Result(); err != nil {
			panic("缓存脚本到主节点失败" + err.Error())
		} else {
			ret = result
		}
		return nil
	})

	return ret
}

var (
	redisCluster *redis.ClusterClient
	scriptLock   string
	scriptDel    string
	lockKey      = "lock_key{key}"
	expireTime   = 10
)

func getLock(index string) bool {
	ret := redisCluster.EvalSha(scriptLock, []string{
		lockKey,
		"192.168.50.48-" + index + "-{key}",
	}, expireTime)

	if result, err := ret.Result(); err != nil {
		return false
	} else {
		return strings.EqualFold("1", fmt.Sprintf("%v", result))
	}
}

func delLock() bool {
	ret := redisCluster.EvalSha(scriptDel, []string{
		lockKey,
	})

	if result, err := ret.Result(); err != nil {
		return false
	} else {
		fmt.Println(result)
		return true
	}
}

func main() {
	// golang 在redis中执行lua脚本
	redisCluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"192.168.50.48:7000",
			"192.168.50.48:7001",
			"192.168.50.48:7002",
			"192.168.50.48:7003",
			"192.168.50.48:7004",
			"192.168.50.48:7005",
		},
	})

	scriptLock = scriptCacheToCluster(redisCluster, createScript())
	scriptDel = scriptCacheToCluster(redisCluster, createScriptDel())

	wait := sync.WaitGroup{}
	count := 3
	wait.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			for {
				if getLock(fmt.Sprintf("%d", i)) {
					fmt.Printf("goroutine %d get lock\n", i)
					if i != 2 {
						time.Sleep(time.Second * 1)
						fmt.Printf("goroutine %d exec succ, del lock key\n", i)
						delLock()
					} else {
						fmt.Printf("goroutine %d exec panic\n", i)
					}
					wait.Done()
					return
				}

				fmt.Printf("goroutine %d not get the lock, wait next time!\n", i)
				time.Sleep(time.Second * 1)
			}
		}(i)
	}

	wait.Wait()
}
