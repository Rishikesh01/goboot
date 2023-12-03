package goboot

import "strings"

type routingTreeNode struct {
	groupPath       string
	basePath        string
	path            string
	method          string
	isChainNotEmpty bool
	subNodes        map[string]*routingTreeNode
	handlerChain    HandlerChain
}

func (r *routingTreeNode) addRoute(path string, method string, handlerChain HandlerChain) {
	r.insertNode(splitPath(path), method, handlerChain)
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

func formatAndValidatePath(path string) string {
	path = strings.ReplaceAll(path, " ", "")
	if path[0] == '/' {
		path = path[1:]
	}

	return path
}

func splitPath(path string) []string {
	return strings.Split(formatAndValidatePath(path), "/")
}
