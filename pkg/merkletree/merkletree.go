package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"log"
	"os"
)

const NUM_CHUNKS int64 = 16

type MerkleTree struct {
	hashes [][]*Node // all hashes
	t      int       // no. of tiers
	n      int       // no. of leaves
}

type Node struct {
	Data   []byte
	Left   *Node
	Right  *Node
	Parent *Node
	leaf   bool
	idx    int // index of this node within its level
}

func Verify(mt *MerkleTree, f *os.File, idx, run int64) bool {
	h, err := hash_file_chunk(f, idx, run)
	if err != nil {
		return false
	}

	return Prove_digest(mt, h.Sum(nil))
}

// Receives a leaf digest and returns whether it's a member of the tree or not
func Prove_digest(mt *MerkleTree, hash []byte) bool {
	idx := find_leaf_index(mt, hash)
	if idx == -1 {
		return false
	}

	var li, ri, pi int
	t := 0

	// Compare our produced hash with the root hash
	for t < mt.t-1 {
		if idx%2 == 0 {
			li = idx
			ri = idx + 1
		} else {
			li = idx - 1
			ri = idx
		}
		pi = li / 2

		hd := hash_childrens_data(mt.hashes[t][li], mt.hashes[t][ri])
		if !bytes.Equal(hd, mt.hashes[t+1][pi].Data) {
			return false
		}

		idx = pi
		t++
	}

	return true
}

func find_leaf_index(mt *MerkleTree, target []byte) int {
	for i := 0; i < mt.n; i++ {
		if bytes.Equal(target, mt.hashes[0][i].Data) {
			return i
		}
	}
	return -1
}

func Init_tree_from_file(fpath string) (*MerkleTree, error) {
	// TODO(?): let user override this, or provide chunk size
	cnum := NUM_CHUNKS

	f, err := os.Open(fpath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close();

	fsize, err := get_file_size(f, cnum)
	if err != nil {
		return nil, err
	}

	cidxs := gen_chunk_indexes(fsize, cnum)
	digests, err := gen_digests_from_file(f, cidxs, cnum)
	if err != nil {
		return nil, err
	}

	mt, err := init_tree_from_digests(digests)
	if err != nil {
		return nil, err
	}

	return mt, nil
}

func get_tier_from_n(n int64) int {
	ret := 0
	for n > 0 {
		n = (n >> 1)
		ret++
	}
	return ret
}

func Print_tree(mt *MerkleTree) {
	for i := 0; i < mt.t; i++ {
		fmt.Printf("\nLevel %d\n\n", i)
		for j := 0; j < len(mt.hashes[i]); j++ {
			fmt.Printf("%x\n", mt.hashes[i][j].Data)
		}
	}
}

func init_tree_from_digests(digests [][]byte) (*MerkleTree, error) {
	nodeCount := len(digests)
	if nodeCount > 0 && nodeCount&(nodeCount-1) != 0 {
		return nil, errors.New("no. of digests is not a power of 2")
	}

	mt, err := init_tree_and_leaves(digests, nodeCount)
	if err != nil {
		return nil, err
	}

	nodeCount /= 2
	t := 1

	for nodeCount > 0 {
		mt.hashes[t] = make([]*Node, nodeCount)
		for i := 0; i < nodeCount; i++ {
			mt.hashes[t][i] = create_parent(i, mt.hashes[t-1][2*i], mt.hashes[t-1][2*i+1])
			mt.hashes[t-1][2*i].Parent = mt.hashes[t][i]
			mt.hashes[t-1][2*i+1].Parent = mt.hashes[t][i]
		}
		nodeCount /= 2
		t++
	}

	return mt, nil
}

func init_tree_and_leaves(digests [][]byte, n int) (*MerkleTree, error) {
	t := get_tier_from_n(int64(n))

	mt := &MerkleTree{
		hashes: make([][]*Node, t),
		t:      t,
		n:      n,
	}

	mt.hashes[0] = make([]*Node, n)
	for i := 0; i < n; i++ {
		mt.hashes[0][i] = &Node{
			Data:   digests[i],
			leaf:   true,
			Left:   nil,
			Right:  nil,
			Parent: nil,
			idx:    i,
		}
	}

	return mt, nil
}

func create_parent(idx int, left, right *Node) *Node {
	d := hash_childrens_data(left, right)
	return &Node{
		Data:   d,
		leaf:   false,
		Left:   left,
		Right:  right,
		Parent: nil,
		idx:    idx,
	}
}

func hash_childrens_data(left, right *Node) []byte {
	// NOTE: we must convert to string else it SHAs the byte data
	conc := fmt.Sprintf("%x", left.Data) + fmt.Sprintf("%x", right.Data)
	hash := sha256.Sum256([]byte(conc))
	return hash[:]
}

func get_file_size(f *os.File, cnum int64) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	if fi.Size() < cnum {
		log.Println("File is too small!")
		return 0, err
	}

	return fi.Size(), nil
}

func gen_chunk_indexes(fsize, cnum int64) []int64 {
	csize := (fsize / cnum)
	cidxs := make([]int64, cnum+1)

	for i := int64(0); i < cnum; i++ {
		cidxs[i] = int64(i) * csize
	}

	cidxs[cnum] = fsize
	return cidxs
}

// Write `target` bytes, starting at the `start`-th byte, in file `f` into a hash object
func hash_file_chunk(f *os.File, start, target int64) (hash.Hash, error) {
	h := sha256.New()
	read := int64(0)

	for read < target {
		size := min(4096, target-read)
		buf := make([]byte, size)

		n, err := f.ReadAt(buf, start+read)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		read += int64(n)
		h.Write(buf)
	}

	return h, nil
}

// From a given file, with cnum chunks at indexes cidxs, return cnum sha256sums
func gen_digests_from_file(f *os.File, cidxs []int64, cnum int64) ([][]byte, error) {
	hashes := make([]hash.Hash, cnum)
	for i := int64(0); i < cnum; i++ {
		hashes[i] = sha256.New()
	}

	for i := 0; i < len(cidxs)-1; i++ {
		// TODO(?): multithread
		h, err := hash_file_chunk(f, cidxs[i], cidxs[i+1]-cidxs[i])
		if err != nil {
			return nil, err
		}
		hashes[i] = h
	}

	sums := make([][]byte, cnum)
	for i, h := range hashes {
		sums[i] = h.Sum(nil)
	}

	return sums, nil
}

func MerkleRoot(m *MerkleTree) []byte {
	return m.hashes[m.t-1][0].Data
}
