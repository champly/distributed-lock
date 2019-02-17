package client

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type Alone struct {
	getSha string
	lock   sync.Mutex
	client *redis.Client
}

func NewAlone(addr string) *Alone {
	return &Alone{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (a *Alone) Get(key, value string, expireTime int) (bool, error) {
	if a.client == nil {
		return false, fmt.Errorf("redis alone mode: client is nil")
	}
	if a.getSha == "" {
		sha, err := a.loadScript(getScript)
		if err != nil {
			return false, err
		}
		a.getSha = sha
	}
	result := a.client.EvalSha(a.getSha, []string{
		prefix + key,
		prefix + value,
	}, expireTime)
	r, err := result.Result()
	if err != nil {
		return false, err
	}
	return strings.EqualFold("1", fmt.Sprintf("%v", r)), nil
}

func (a *Alone) Del(key string) (bool, error) {
	result := a.client.Del(prefix + key)
	r, err := result.Result()
	if err != nil {
		return false, err
	}
	return strings.EqualFold("1", fmt.Sprintf("%v", r)), nil
}

func (a *Alone) Delay(key string, expireTime int) (bool, error) {
	return a.client.Expire(prefix+key, time.Second*time.Duration(expireTime)).Result()
}

func (a *Alone) loadScript(str string) (string, error) {
	a.lock.Lock()
	a.lock.Unlock()
	if a.getSha != "" {
		return a.getSha, nil
	}

	script := redis.NewScript(str)
	ret, err := script.Load(a.client).Result()
	if err != nil {
		return "", fmt.Errorf("redis alone mode: load script error:%s", err.Error())
	}
	return ret, nil
}
