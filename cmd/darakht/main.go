package main

import (
	"fmt"
	"log"
	"os"

	"github.com/masroof-maindak/darakht/pkg/merkletree"
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
		log.Println("Error generating merkletree: ", err)
		return
	}
	merkletree.PrintTree(mt)

	//
	// Proving membership of the hash of "ef"
	//
	f, err := os.Open(fpath)
	if err != nil {
		log.Println("Error opening file for verification: ", err)
		return
	}
	defer f.Close()

	exists, err := merkletree.Verify(mt, f, 4, 2)
	if err != nil {
		log.Println("Error verifying block: ", err)
		return
	}

	if exists {
		fmt.Println("Digest exists")
	} else {
		fmt.Println("Digest does not exist")
	}
}
