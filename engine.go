package goboot

import (
	"net/http"
	"path"
	"sync"
)

type (
	Handler      func(*Context)
	HandlerChain []Handler
)

type Engine struct {
	RouteGroup
	pool     sync.Pool
	rootNode *routingTreeNode
}

func Default() *Engine {
	engine := &Engine{
		RouteGroup: RouteGroup{
			BasePath:              "/",
			CommonMiddleWareChain: nil,
		},
		rootNode: &routingTreeNode{
			fullPath: "/",
			method:   []methodHandler{},
			subNodes: map[string]*routingTreeNode{},
		},
	}
	engine.pool.New = func() any {
		return &Context{}
	}

	engine.RouteGroup.engine = engine

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

	for i := range urlPath {
		child := node.getNodeByPath(urlPath[i])

		if i == len(urlPath)-1 && child != nil {
			for j := range child.method {
				if child.method[j].method == ctx.Request.Method {
					for k := range child.method[j].handlerChain {
						child.method[j].handlerChain[k](ctx)
					}
				}
			}
			return
		}

		node = child
	}

	return
}

type RouteGroup struct {
	BasePath              string
	CommonMiddleWareChain HandlerChain
	engine                *Engine
}

func (r *RouteGroup) handleRoute(relativePath, method string, handler ...Handler) {
	r.engine.rootNode.addRoute(path.Join(r.BasePath, relativePath), method, append(r.CommonMiddleWareChain, handler...))
}

func (r *RouteGroup) Group(basePath string) *RouteGroup {
	return &RouteGroup{
		BasePath:              path.Join(r.BasePath, basePath),
		CommonMiddleWareChain: r.CommonMiddleWareChain,
		engine:                r.engine,
	}
}

func (r *RouteGroup) MiddleWare(handler ...Handler) {
	r.CommonMiddleWareChain = append(r.CommonMiddleWareChain, handler...)
}

func (s *RouteGroup) GET(path string, handler ...Handler) {
	s.handleRoute(path, http.MethodGet, handler...)
}

func (s *RouteGroup) POST(path string, handler ...Handler) {
	s.handleRoute(path, http.MethodPost, handler...)
}

func (s *RouteGroup) PATCH(path string, handler ...Handler) {
	s.handleRoute(path, http.MethodPatch, handler...)
}

func (s *RouteGroup) PUT(path string, handler ...Handler) {
	s.handleRoute(path, http.MethodPut, handler...)
}

func (s *RouteGroup) DELETE(path string, handler ...Handler) {
	s.handleRoute(path, http.MethodDelete, handler...)
}

func (s *RouteGroup) HEAD(path string, handler ...Handler) {
	s.handleRoute(path, http.MethodHead, handler...)
}
