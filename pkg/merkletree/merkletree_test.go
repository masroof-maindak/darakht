package merkletree_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/masroof-maindak/darakht/pkg/merkletree"
)

const fpath = "sample.txt"
const testCnum = 16

func createTmpFile() error {
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("abcdefghijklmnopqrstuvwxyz123456789")
	if err != nil {
		return err
	}

	return nil
}

func TestMerkleRoot(t *testing.T) {
	if err := createTmpFile(); err != nil {
		t.Error(err)
	}
	defer os.Remove(fpath)

	mt, err := merkletree.InitTreeFromFile(fpath, testCnum)
	if err != nil {
		t.Error(err)
	}

	root := fmt.Sprintf("%x", merkletree.MerkleRoot(mt))
	if root != "857a51b0311986c84ed67794c70fa4509c0e744aa69cda1774514d02dbbad7cb" {
		t.Error("Invalid Merkle Root!")
	}
}

func TestVerify(t *testing.T) {
	if err := createTmpFile(); err != nil {
		t.Error(err)
	}
	defer os.Remove(fpath)

	mt, err := merkletree.InitTreeFromFile(fpath, testCnum)
	if err != nil {
		t.Error(err)
	}

	f, err := os.Open(fpath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	exists, err := merkletree.ProveMember(mt, f, 4, 2)
	if err != nil {
		t.Error(err)
	}

	if !exists {
		t.Error("Expected byte range to exist in Merkle tree")
	}
}
