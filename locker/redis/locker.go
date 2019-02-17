package redis

import "github.com/champly/distributed-lock/locker/redis/client"

type Locker struct {
	client client.IClient
}

func NewLockerAlone(addr string) *Locker {
	cli := client.NewAlone(addr)
	return &Locker{
		client: cli,
	}
}

func NewLockerCluster(addrs []string) *Locker {
	cli := client.NewCluster(addrs)
	return &Locker{
		client: cli,
	}
}

func NewLockerWithClient(client client.IClient) *Locker {
	return &Locker{
		client: client,
	}
}

func (r *Locker) Get(key, value string, expireTime int) (bool, error) {
	return r.client.Get(key, value, expireTime)
}

func (r *Locker) Del(key string) (bool, error) {
	return r.client.Del(key)
}

func (r *Locker) Delay(key string, expireTime int) (bool, error) {
	return r.client.Delay(key, expireTime)
}
