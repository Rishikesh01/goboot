package goboot

import (
	"testing"
)

func TestBoot(t *testing.T) {
	s := Default()

	s.GET("/v1", func(ctx *Context) {
		ctx.String(200, "hi")
	})
	s.GET("/v1/:ho", func(ctx *Context) {
		ctx.String(200, "h")
	})

	s.GET("/v1/:ho", func(ctx *Context) {
		ctx.String(200, "hi")
	})

	s.Run(":8080")
}

func TestSplitPath(t *testing.T) {
	// listOfInvalidPath := []string{"/home//hello", "/something/ /something", "/seom/ s/something"}
}
