package main

import (
	"flag"
	"log"
	"os"

	"github.com/masroof-maindak/darakht/pkg/merkletree"
)

func main() {
	fpath := flag.String("f", "", "path to file")
	cnum := flag.Int64("c", merkletree.NUM_CHUNKS, "number of chunks/leaves")

	flag.Parse()

	if *fpath == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*fpath)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	defer f.Close()

	constructAndPrint(f, *cnum)
	// proveMembership()
}

// Construct a Merkle Tree from a file and print it to stdout
func constructAndPrint(f *os.File, cnum int64) {
	mt, err := merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		log.Println("Merkle Tree generation failed:", err)
		return
	}

	if err = mt.Print(); err != nil {
		log.Println(err)
		return
	}
}

func proveMembership() {
	// 1. Reconstruct (and verify) tree from JSON
	// 2. Open file
	// 3. mt.ProveMember(f, 4, 2)

	// exists, err := mt.ProveMember(f, 4, 2)
	// if err != nil {
	// 	log.Println("Error verifying block:", err)
	// 	return
	// }
	//
	// if exists {
	// 	fmt.Println("true")
	// } else {
	// 	fmt.Println("false")
	// }
}
