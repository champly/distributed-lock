package client

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type Cluster struct {
	getSha string
	lock   sync.Mutex
	client *redis.ClusterClient
}

func NewCluster(addrs []string) *Cluster {
	return &Cluster{
		client: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: addrs,
		}),
	}
}

func (c *Cluster) Get(key, value string, expireTime int) (bool, error) {
	if c.client == nil {
		return false, fmt.Errorf("redis cluster mode: client is nil")
	}
	if c.getSha == "" {
		sha, err := c.loadScript(getScript)
		if err != nil {
			return false, err
		}
		c.getSha = sha
	}
	result := c.client.EvalSha(c.getSha, []string{
		prefix + key,
		prefix + value,
	}, expireTime)
	r, err := result.Result()
	if err != nil {
		return false, err
	}
	return strings.EqualFold("1", fmt.Sprintf("%v", r)), nil
}

func (c *Cluster) Del(key string) (bool, error) {
	result := c.client.Del(prefix + key)
	r, err := result.Result()
	if err != nil {
		return false, err
	}
	return strings.EqualFold("1", fmt.Sprintf("%v", r)), nil
}

func (c *Cluster) Delay(key string, expireTime int) (bool, error) {
	return c.client.Expire(prefix+key, time.Second*time.Duration(expireTime)).Result()
}

func (c *Cluster) loadScript(str string) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.getSha != "" {
		return c.getSha, nil
	}

	script := redis.NewScript(str)
	var sha string
	err := c.client.ForEachMaster(func(m *redis.Client) error {
		result, err := script.Load(m).Result()
		if err != nil {
			return fmt.Errorf("redis cluster mode: load script error:%s", err.Error())
		}
		sha = result
		return nil
	})

	return sha, err
}
