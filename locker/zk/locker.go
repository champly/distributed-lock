package zk

import (
	"fmt"
	"strings"

	"github.com/samuel/go-zookeeper/zk"
)

type Locker struct {
	client *zk.Conn
}

func NewLockerWithClient(cli *zk.Conn) *Locker {
	return &Locker{client: cli}
}

func (l *Locker) createSeqNode(path, value string) (bool, error) {
	exists, _, err := l.client.Exists(path)
	if exists {
		result, err := l.client.Create(path, []byte(value), zk.FlagSequence, zk.WorldACL(zk.PermAll))
		if err != nil {
			return false, err
		}
		fmt.Println(result)
		return true, nil
	}

	// loop create zk node
	paths := []string{}
	pp := strings.Trim(path, "/")
	pslice := strings.Split(pp, "/")
	fmt.Println(pslice)
	for i := 1; i < len(pslice); i++ {
		fmt.Println("/" + strings.Join(pslice[:i], "/"))
		paths = append(paths, "/"+strings.Join(pslice[:i], "/"))
	}

	for _, p := range paths {
		b, _, err := l.client.Exists(p)
		if err != nil {
			return false, err
		}
		if !b {
			l.client.Create(p, []byte(""), 0, zk.WorldACL(zk.PermAll))
		}
	}

	result, err := l.client.Create(path, []byte(value), zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return false, err
	}
	fmt.Println(result)
	return true, nil
}

func (l *Locker) Get(path, value string) (bool, error) {
	if l.client == nil {
		return false, fmt.Errorf("zk conn socker is nil")
	}
	return l.createSeqNode(path, value)
}

func (l *Locker) GetChildren(path string) []string {
	result, _, _ := l.client.Children(path)
	return result
}
