package xcache

type (
	ItemVal struct {
		Key    string
		Val    interface{}
		Expire int64
	}
)
