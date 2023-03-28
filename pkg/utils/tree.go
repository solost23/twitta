package utils

type TreeNode interface {
	GetPID() string
	GetID() string
	AppendChildrenNode(node TreeNode)
}

func BuildTree(array []TreeNode) TreeNode {
	maxLen := len(array)
	var treeNode TreeNode = nil

	for i := 0; i < maxLen; i++ {
		count := 0
		for j := 0; j < maxLen; j++ {
			if array[j].GetID() == array[i].GetPID() {
				count++
				array[j].AppendChildrenNode(array[i])
				break
			}
		}
		if count <= 0 {
			treeNode = array[i]
		}
	}
	return treeNode
}

func BuildTrees(array []TreeNode) []TreeNode {
	maxLen := len(array)
	treeNodes := make([]TreeNode, 0)

	for i := 0; i < maxLen; i++ {
		count := 0
		for j := 0; j < maxLen; j++ {
			if array[j].GetID() == array[i].GetPID() {
				count++
				array[j].AppendChildrenNode(array[i])
				break
			}
		}
		if count <= 0 {
			treeNodes = append(treeNodes, array[i])
		}
	}
	return treeNodes
}
