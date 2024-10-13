package merkletree_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/masroof-maindak/darakht/pkg/merkletree"
)

const (
	cnum  int64  = 16
	fname string = "sample-*.txt"

	equalityErr string = "Expected merkle trees to be equal"
	memberErr   string = "Expected byte range to exist in Merkle tree"
	rootErr     string = "Invalid Merkle Root"
)

var (
	f *os.File = nil
)

func TestMain(m *testing.M) {
	var err error
	f, err = os.CreateTemp("", fname)
	if err != nil {
		log.Println("Error creating test file")
		os.Exit(1)
	}

	f.WriteString("abcdefghijklmnopqrstuvwxyz123456789")

	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Error closing temp file:", err)
		}
		os.Remove(f.Name())
	}()

	os.Exit(m.Run())
}

func TestConstructionAndMerkleRoot(t *testing.T) {
	mt, err := merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		t.Error(err)
	}

	root := fmt.Sprintf("%x", mt.MerkleRoot())
	if root != "857a51b0311986c84ed67794c70fa4509c0e744aa69cda1774514d02dbbad7cb" {
		t.Error(rootErr)
	}
}

func TestEquals(t *testing.T) {
	var mt1, mt2 *merkletree.MerkleTree
	var err error

	mt1, err = merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		t.Error(err)
	}
	mt2 = mt1
	if !mt1.Equals(mt2) {
		t.Error(equalityErr)
	}

	mt2, err = merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		t.Error(err)
	}
	if !mt1.Equals(mt2) {
		t.Error(equalityErr)
	}
}

func TestProveMember(t *testing.T) {
	mt, err := merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		t.Error(err)
	}

	exists, err := mt.ProveMember(f, 4, 2)
	if err != nil {
		t.Error(err)
	}

	if !exists {
		t.Error(memberErr)
	}
}
