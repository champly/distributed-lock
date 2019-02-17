package client

import "github.com/go-redis/redis"

type IClient interface {
	Get(key, value string, expireTime int) (bool, error)
	Del(key string) (bool, error)
	Delay(key string, expireTime int) (bool, error)
}

func NewWithAloneClient(client *redis.Client) *Alone {
	return &Alone{client: client}
}

func NewWithClusterCluster(client *redis.ClusterClient) *Cluster {
	return &Cluster{client: client}
}
