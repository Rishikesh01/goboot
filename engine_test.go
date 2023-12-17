package goboot

import (
	"testing"
)

func TestBoot(t *testing.T) {
	s := Default()
	v1 := s.Group("/v1")
	v1.GET(":hello", func(ctx *Context) {
		ctx.String(200, "hello world")
	})

	v1.GET("hello", func(ctx *Context) {
		ctx.String(200, "hi")
	})
	s.GET("hello", func(ctx *Context) {
		ctx.String(200, "not v1 path")
	})

	s.Run(":8080")
}

func TestSplitPath(t *testing.T) {
	// listOfInvalidPath := []string{"/home//hello", "/something/ /something", "/seom/ s/something"}
}
