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
			groupPath:       "/",
			basePath:        "/",
			path:            "/",
			method:          "*",
			isChainNotEmpty: false,
			subNodes:        map[string]*routingTreeNode{},
			handlerChain:    make([]Handler, 10),
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
}

func (s *Engine) handlerRequest(ctx *Context) {
	// break path into segments
	urlPath := splitPath(ctx.Request.URL.Path)
	node := s.rootNode

	for i := 0; i < len(urlPath); i++ {
		child := node.getNodeByPath(urlPath[i])
		if child != nil {
			node = child
			continue
		}
		for _, apply := range node.handlerChain {
			apply(ctx)
		}
		return
	}

	for i := 0; i < len(node.handlerChain); i++ {
		node.handlerChain[i](ctx)
	}
	return
}

func (s *Engine) GET(path string, handler ...Handler) {
	s.rootNode.addRoute(path, http.MethodGet, handler)
}
