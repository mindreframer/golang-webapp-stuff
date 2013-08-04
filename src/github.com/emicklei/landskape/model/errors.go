package model

import "github.com/emicklei/hopwatch"

func Debug(args ...interface{}) {
	hopwatch.CallerOffset(3).Display(args...).Break()
}
func Display(args ...interface{}) {
	hopwatch.CallerOffset(3).Display(args...)
}

const (
	MISMATCH_ID    = 4001 // first special BadRequest(400)
	MISMATCH_SCOPE = 4002 // first special BadRequest(400)
)
