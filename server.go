package goboot

import (
	"net/http"
	"strings"
)

type (
	Handler      func(*Context)
	HandlerChain []Handler
)

type Server struct {
	rootNode *routingTreeNode
	port     string
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	path := splitPath(req.URL.Path)[1:]
	node := s.rootNode
	for i, subPath := range path {
		child := node.getNodeByPath(subPath)
		if child != nil && i == len(path)-1 {
			if child.method == req.Method {
				for _, chain := range child.handlerChain {
					c := &Context{Writer: writer, Request: req}
					chain(c)
				}
			}
		}
	}
}

func (s *Server) run() {
	http.ListenAndServe(":8080", s)
}

type routingTreeNode struct {
	groupPath       string
	basePath        string
	path            string
	method          string
	isChainNotEmpty bool
	subNodes        map[string]*routingTreeNode
	handlerChain    HandlerChain
}

func (s *Server) addRoute(path string, method string, handlerChain HandlerChain) {
	if s.rootNode == nil {
		s.rootNode = &routingTreeNode{
			groupPath:       "/",
			basePath:        "/",
			path:            "/",
			method:          "*",
			isChainNotEmpty: false,
			subNodes:        make(map[string]*routingTreeNode),
		}
	}
	s.rootNode.insertNode(splitPath(path)[1:], method, handlerChain)
}

func (r *routingTreeNode) insertNode(path []string, method string, handlerChain HandlerChain) {
	node := r
	for i, subPath := range path {
		child := node.getNodeByPath(subPath)
		if child == nil {
			child = &routingTreeNode{
				groupPath: subPath,
				basePath:  subPath,
				path:      subPath,
				method:    "",
				subNodes:  make(map[string]*routingTreeNode),
			}
		}
		node.subNodes[subPath] = child
		node = child
		if len(path)-1 == i {
			node.isChainNotEmpty = true
			node.handlerChain = handlerChain
			node.method = method
		}
	}
}

func (r *routingTreeNode) getNodeByPath(path string) *routingTreeNode {
	if child, ok := r.subNodes[path]; ok {
		return child
	}

	return nil
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}
