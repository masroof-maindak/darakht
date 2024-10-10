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
	// Constructing Merkle Tree
	//
	mt, err := merkletree.InitTreeFromFile(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating merkletree")
		return
	}
	merkletree.PrintTree(mt)

	//
	// Proving membership of the hash of "cd"
	//
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file for verification")
		return
	}
	defer f.Close()

	exists := merkletree.Verify(mt, f, 4, 2)

	if exists {
		fmt.Println("Digest exists")
	} else {
		fmt.Println("Digest does not exist")
	}
}
