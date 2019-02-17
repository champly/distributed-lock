package zk

import (
	"fmt"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func TestGet(t *testing.T) {
	client, _, _ := zk.Connect([]string{
		"192.168.50.48:2181",
	}, time.Second*1)
	locker := NewLockerWithClient(client)
	locker.Get("/test/lock/zk", time.Now().Format("2006-01-02 15:04:05"))
	children, _, _ := locker.client.Children("/test/lock")
	fmt.Println(children)
	for _, c := range children {
		_, s, _ := locker.client.Get("/test/lock/" + c)
		locker.client.Delete("/test/lock/"+c, s.Version)
	}
}
