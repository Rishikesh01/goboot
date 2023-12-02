package goboot

import (
	"net/http"
	"testing"
)

func TestBoot(t *testing.T) {
	s := &Server{}

	s.addRoute("/hello", http.MethodGet, HandlerChain{func(ctx *Context) {
		s := make(map[string]string)
		s["h"] = "ellow"
		ctx.JSON(200, s)
	}})

	s.run()
}
