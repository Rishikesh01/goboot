package goboot

import (
	"testing"
)

func TestBoot(t *testing.T) {
	s := Default()

	s.GET("/hello/home", func(ctx *Context) {
		ctx.String(200, "home")
	})

	s.GET("/v1/:ho", func(ctx *Context) {
		ctx.String(200, "h")
	})

	s.Run(":8080")
}

func TestSplitPath(t *testing.T) {
	// listOfInvalidPath := []string{"/home//hello", "/something/ /something", "/seom/ s/something"}
}
