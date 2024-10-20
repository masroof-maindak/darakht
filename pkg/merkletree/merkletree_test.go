package merkletree_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/masroof-maindak/darakht/pkg/merkletree"
)

const (
	cnum      int64  = 16
	fname     string = "sample-*.txt"
	fnameJson string = "sample-json-*.txt"

	equalityErr string = "Expected merkle trees to be equal"
	memberErr   string = "Expected byte range to exist in Merkle tree"
	rootErr     string = "Invalid Merkle Root"
)

var (
	f     *os.File = nil
	fJson *os.File = nil
)

func closeAndRemoveTempFile(_f *os.File) {
	if err := _f.Close(); err != nil {
		log.Println("Error closing temp file:", err)
	}

	if err := os.Remove(_f.Name()); err != nil {
		log.Println("Error removing temp file:", err)
	}
}

func TestMain(m *testing.M) {
	var err error
	f, err = os.CreateTemp("", fname)
	if err != nil {
		log.Println("Error creating test file")
		os.Exit(1)
	}
	defer closeAndRemoveTempFile(f)

	f.WriteString("abcdefghijklmnopqrstuvwxyz123456789")

	fJson, err = os.CreateTemp("", fnameJson)
	if err != nil {
		log.Println("Error creating test JSON file")
		os.Exit(2)
	}
	defer closeAndRemoveTempFile(fJson)

	os.Exit(m.Run())
}

func TestValidateJSONFile(t *testing.T) {
	mt1, err := merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		t.Error(err)
		return
	}

	if !rootMatch(mt1) {
		t.Error(rootErr)
		return
	}

	err = mt1.Serialise(fJson)
	if err != nil {
		t.Error(err)
		return
	}

	mt2, err := merkletree.NewMerkleTreeFromJSON(fJson)
	if err != nil {
		t.Error(err)
		return
	}

	if !rootMatch(mt2) || !mt1.Equals(mt2) {
		t.Error(rootErr)
		return
	}
}

func rootMatch(mt *merkletree.MerkleTree) bool {
	root := fmt.Sprintf("%x", mt.MerkleRoot())
	if root != "857a51b0311986c84ed67794c70fa4509c0e744aa69cda1774514d02dbbad7cb" {
		return false
	}
	return true
}

func TestConstructionAndMerkleRoot(t *testing.T) {
	mt, err := merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		t.Error(err)
		return
	}

	if !rootMatch(mt) {
		t.Error(rootErr)
		return
	}
}

func TestEquals(t *testing.T) {
	var mt1, mt2 *merkletree.MerkleTree
	var err error

	mt1, err = merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		t.Error(err)
		return
	}

	mt2 = mt1

	if !mt1.Equals(mt2) {
		t.Error(equalityErr)
		return
	}

	mt2, err = merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		t.Error(err)
		return
	}
	if !mt1.Equals(mt2) {
		t.Error(equalityErr)
		return
	}
}

func TestProveMember(t *testing.T) {
	mt, err := merkletree.NewMerkleTreeFromFile(f, cnum)
	if err != nil {
		t.Error(err)
		return
	}

	exists, err := mt.ProveMember(f, 4, 2)
	if err != nil {
		t.Error(err)
		return
	}

	if !exists {
		t.Error(memberErr)
		return
	}
}
