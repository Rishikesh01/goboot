package goboot

import (
	"net/http"
	"sync"
)

type (
	Handler      func(*Context)
	HandlerChain []Handler
)

type Engine struct {
	pool     sync.Pool
	rootNode *routingTreeNode
}

func Default() *Engine {
	engine := &Engine{
		rootNode: &routingTreeNode{
			fullPath:     "/",
			method:       []methodHandler{},
			subNodes:     map[string]*routingTreeNode{},
			handlerChain: make([]Handler, 10),
		},
	}
	engine.pool.New = func() any {
		return &Context{}
	}

	return engine
}

func (e *Engine) Run(port string) {
	http.ListenAndServe(port, e)
}

func (s *Engine) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	context := s.pool.Get().(*Context)
	context.Request = req
	context.Writer = writer
	s.handlerRequest(context)
	s.pool.Put(context)
}

func (s *Engine) handlerRequest(ctx *Context) {
	// break path into segments
	urlPath := splitPath(ctx.Request.URL.Path)
	node := s.rootNode

	for i := 0; i < len(urlPath); i++ {
		child := node.getNodeByPath(urlPath[i])

		if i == len(urlPath)-1 && child != nil {
			for i := range child.handlerChain {
				child.handlerChain[i](ctx)
			}
			return
		}

		node = child
	}

	return
}

func (s *Engine) GET(path string, handler ...Handler) {
	s.rootNode.addRoute(path, http.MethodGet, handler)
}

func (s *Engine) POST(path string, handler ...Handler) {
	s.rootNode.addRoute(path, http.MethodPost, handler)
}

func (s *Engine) PATCH(path string, handler ...Handler) {
	s.rootNode.addRoute(path, http.MethodPatch, handler)
}

func (s *Engine) PUT(path string, handler ...Handler) {
	s.rootNode.addRoute(path, http.MethodPut, handler)
}

func (s *Engine) DELETE(path string, handler ...Handler) {
	s.rootNode.addRoute(path, http.MethodDelete, handler)
}

func (s *Engine) HEAD(path string, handler ...Handler) {
	s.rootNode.addRoute(path, http.MethodHead, handler)
}
