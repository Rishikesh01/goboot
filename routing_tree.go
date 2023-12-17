package goboot

import (
	"strings"
)

type routingTreeNode struct {
	fullPath    string
	hasWildCard bool
	wildCards   []wildCard
	method      []methodHandler
	subNodes    map[string]*routingTreeNode
}

type methodHandler struct {
	method       string
	handlerChain HandlerChain
}

type wildCard struct {
	name     string
	startPos int
}

func (r *routingTreeNode) addRoute(path string, method string, handlerChain HandlerChain) {
	r.insertNode(path, method, handlerChain)
}

func (r *routingTreeNode) insertNode(fullPath string, method string, handlerChain HandlerChain) {
	node := r
	path := splitPath(fullPath)
	wildCards := []wildCard{}
	for i, subPath := range path {
		isWildCard := r.isWildCard(subPath)
		child := node.getNodeByPath(subPath)
		if isWildCard && r.hasWildCard && child == nil {
			panic("existing wildCard")
		}
		if isWildCard {
			r.hasWildCard = true
			wildCards = append(wildCards, wildCard{name: subPath[1:], startPos: i})
		}
		if child == nil && len(path)-1 != i {
			child = &routingTreeNode{
				subNodes: make(map[string]*routingTreeNode),
			}
		}
		if len(path)-1 == i && child == nil {
			child = &routingTreeNode{
				fullPath:  fullPath,
				wildCards: wildCards,
				method:    []methodHandler{},
				subNodes:  make(map[string]*routingTreeNode),
			}
			child.method = append(child.method, methodHandler{method: method, handlerChain: handlerChain})
		} else if len(path)-1 == i && child != nil {
			for _, httpMethod := range child.method {
				if method == httpMethod.method {
					panic("error for same HTTP method already exists for the path")
				}
			}
			child.method = append(child.method, methodHandler{method: method, handlerChain: handlerChain})
		}

		node.subNodes[subPath] = child
		node = child
	}
}

func (r *routingTreeNode) getNodeByPath(path string) *routingTreeNode {
	if child, ok := r.subNodes[path]; ok {
		return child
	}

	return nil
}

func formatAndValidatePath(path string) string {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	path = strings.TrimSpace(path)
	if path[0] == '/' {
		path = path[1:]
	}

	return path
}

func splitPath(path string) []string {
	return strings.Split(formatAndValidatePath(path), "/")
}

func (r *routingTreeNode) isWildCard(path string) bool {
	return []byte(path)[0] == ':' || []byte(path)[0] == '*'
}
