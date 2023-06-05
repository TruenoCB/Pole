package poleweb

import "testing"

func TestTree(t *testing.T) {
	tree := &treeNode{
		name: "/",
	}
	tree.Put("/asdf/user/a")
	tree.Put("/user/a")
	tree.Get("/user")
	tree.Get("/user/a")
}
