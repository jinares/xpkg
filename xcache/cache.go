package xcache

type (
	Cache interface {
		Get(key string) (string, error)
		MGet(key ...string) (map[string]string, error)
		Set(key, val string) error
		Incr(key string, val int) (int64, error)
		Expire(key string, expire int) error
	}
)
