package goboot

import (
	"testing"
)

func TestBoot(t *testing.T) {
	s := Default()

	s.GET("/hello", func(ctx *Context) {
		s := make(map[string]string)
		s["h"] = "ellow"
		ctx.JSON(200, s)
	})

	s.Run(":8080")
}

func TestSplitPath(t *testing.T) {
	// listOfInvalidPath := []string{"/home//hello", "/something/ /something", "/seom/ s/something"}
}
