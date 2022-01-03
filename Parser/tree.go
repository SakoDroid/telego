package Parser

import (
	"strings"

	objs "github.com/SakoDroid/telego/objects"
)

//A special tree element containing handlers.
type TreeNode struct {
	father      *TreeNode
	right, left *TreeNode
	data        *handler
}

//A special binary tree for storing handlers. Right node hase a value that does not match the it's father regex and the left node matches it's father regex.
type HandlerTree struct {
	root *TreeNode
}

//Adds a new handler to the tree.
func (tr *HandlerTree) AddHandler(hdl *handler) {
	tn := TreeNode{data: hdl}
	tr.addNode(&tn)
}

//Gets the proper handler for the given text.
func (tr *HandlerTree) GetHandler(msg *objs.Message) *handler {
	tn := tr.findTheNodeRegex(msg.Text, msg.Chat.Type)
	if tn != nil {
		return tn.data
	} else {
		return nil
	}
}

func (tr *HandlerTree) findTheNodeRegex(text, chatType string) *TreeNode {
	node := tr.root
	for {
		if node == nil {
			break
		}
		if node.data.regex.Match([]byte(text)) {
			if node.left != nil {
				node = node.left
			} else {
				break
			}
		} else {
			if node.right != nil {
				node = node.right
			} else {
				node = node.father
				break
			}
		}
	}
	return tr.checkForChatTypes(node, chatType)
}

func (tr *HandlerTree) checkForChatTypes(currentNode *TreeNode, chatType string) *TreeNode {
	for {
		if currentNode == nil {
			break
		}
		if strings.Contains(currentNode.data.chatType, chatType) || strings.Contains(currentNode.data.chatType, "all") {
			break
		} else {
			currentNode = currentNode.father
		}
	}
	return currentNode
}

//Finds the perfect location for this handler.
func (tr *HandlerTree) addNode(tn *TreeNode) {
	var fatherNode *TreeNode
	node := tr.root
	dir := 0
	for {
		//this is the spot
		if node == nil {
			if fatherNode != nil {
				tn.father = fatherNode
				if dir == 0 {
					fatherNode.left = tn
				} else {
					fatherNode.right = tn
				}
			} else {
				//Root node
				tr.root = tn
			}
			break
		} else {
			if tr.checkRegex(node, tn) {
				dir = 0
				fatherNode = node
				node = fatherNode.left
			} else {
				dir = 1
				fatherNode = node
				node = fatherNode.right
			}
		}
	}
}

func (tr *HandlerTree) checkRegex(father, child *TreeNode) bool {
	return father.data.regex.Match(
		[]byte(child.data.regex.String()),
	)
}
