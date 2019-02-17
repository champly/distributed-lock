package redis

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// docker search grokzen/redis-cluster
// docker run --rm --name redis-cluster -p 7000:7000 -p 7001:7001 -p 7002:7002 -p 7003:7003 -p 7004:7004 -p 7005:7005 -e "IP=0.0.0.0" grokzen/redis-cluster
func TestGet(t *testing.T) {
	client := NewLockerCluster([]string{
		"192.168.50.48:7000",
		"192.168.50.48:7001",
		"192.168.50.48:7002",
		"192.168.50.48:7003",
		"192.168.50.48:7004",
		"192.168.50.48:7005",
	})

	wait := sync.WaitGroup{}
	count := 3
	wait.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			for {
				if b, err := client.Get("test_key", time.Now().Format("2006-01-02 15:04:05"), 10); err == nil && b {
					fmt.Printf("goroutine %d get lock\n", i)
					if i != 2 {
						time.Sleep(time.Second * 1)
						fmt.Printf("goroutine %d exec succ, del lock key\n", i)
						client.Del("test_key")
					} else {
						fmt.Printf("goroutine %d exec 8 second\n", i)
						time.Sleep(time.Second * 8)
						client.Delay("test_key", 10)
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
