package locker

type ILocker interface {
	Get(key, value string, expireTime int) (bool, error)
	Del(key string) (bool, error)
	Delay(key string, expireTime int) (bool, error)
}
