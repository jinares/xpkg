package xboomer

import (
	"github.com/jinares/xpkg/xtools"
	"github.com/myzhan/boomer"
	"time"
)

type (
	BoomerHandler func() (interface{}, error)
)

func Action(name string, fn BoomerHandler) *boomer.Task {
	return &boomer.Task{
		Weight: 100,
		Fn: func() {
			start := time.Now()
			ret, err := fn()

			elapsed := time.Since(start)
			if err == nil {
				boomer.RecordSuccess(
					"func", name,
					elapsed.Milliseconds(),
					int64(len(xtools.JSONToStr(ret))),
				)
				return

			}

			boomer.RecordFailure(
				"func", name,
				elapsed.Milliseconds(),
				err.Error(),
			)
		},
		Name: name,
	}
}
func Run(tasks ...*boomer.Task) {
	boomer.Run(tasks...)
}
