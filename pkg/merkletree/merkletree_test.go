package merkletree_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/masroof-maindak/darakht/pkg/merkletree"
)

const testCnum = 16

var testF *os.File

func TestMain(m *testing.M) {
	var err error
	testF, err = os.CreateTemp("", "sample-*.txt")
	if err != nil {
		log.Println("Error creating test file")
		os.Exit(1)
	}

	testF.WriteString("abcdefghijklmnopqrstuvwxyz123456789")

	defer func() {
		if err := testF.Close(); err != nil {
			log.Println("Error closing temp file:", err)
		}
		os.Remove(testF.Name())
	}()

	os.Exit(m.Run())
}

func TestMerkleRoot(t *testing.T) {
	mt, err := merkletree.InitTreeFromFile(testF, testCnum)
	if err != nil {
		t.Error(err)
	}

	root := fmt.Sprintf("%x", merkletree.MerkleRoot(mt))
	if root != "857a51b0311986c84ed67794c70fa4509c0e744aa69cda1774514d02dbbad7cb" {
		t.Error("Invalid Merkle Root!")
	}
}

func TestProveMember(t *testing.T) {
	mt, err := merkletree.InitTreeFromFile(testF, testCnum)
	if err != nil {
		t.Error(err)
	}

	exists, err := merkletree.ProveMember(mt, testF, 4, 2)
	if err != nil {
		t.Error(err)
	}

	if !exists {
		t.Error("Expected byte range to exist in Merkle tree")
	}
}
