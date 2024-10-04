package main

import (
	"Merkle/pkg/merkletree"

	"fmt"
)

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Println("Usage: ", os.Args[0], "<filename>")
	// 	return
	// }

	// fpath := os.Args[1]
	fpath := "../../misc/sample-35b.txt"

	mt, err := merkletree.Init_tree_from_file(fpath)
	if err != nil {
		fmt.Println("idiot")
		return
	}

	merkletree.Print_tree(mt)
}
