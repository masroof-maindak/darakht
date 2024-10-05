package main

import (
	"Merkle/pkg/merkletree"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "<filename>")
		return
	}
	fpath := os.Args[1]

	//
	// Constructing merkletree
	//
	mt, err := merkletree.Init_tree_from_file(fpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating merkletree\n")
		return
	}
	merkletree.Print_tree(mt)

	//
	// Proving membership of the hash of "cd"
	//
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file for verification\n")
		return
	}
	defer f.Close()

	exists := merkletree.Verify(mt, f, 1, 3)

	if exists {
		fmt.Println("Digest exists")
	} else {
		fmt.Println("Digest does not exist")
	}
}
