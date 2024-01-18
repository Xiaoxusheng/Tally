package global

type Mutex interface {
	Lock(key, id string) bool
	Unlock(key string) bool
}
