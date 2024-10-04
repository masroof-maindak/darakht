package merkletree

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"os"
)

const (
	NUM_CHUNKS int64 = 16
)

type MerkleTree struct {
	hashes [][]*Node // all hashes
	t      int64     // no. of tiers
	n      int       // no. of leaves -- TODO: remove?
}

type Node struct {
	Data   []byte
	Left   *Node
	Right  *Node
	Parent *Node
	leaf   bool
	idx    int // index of this node within its level
}

func Init_tree_from_file(fpath string) (*MerkleTree, error) {
	f, fsize, err := open_file_and_get_size(fpath)
	if err != nil {
		return nil, err
	}

	// TODO: let user override this
	cnum := NUM_CHUNKS

	cidxs := gen_chunk_indexes(fsize, cnum)
	digests, err := gen_digests_from_file(f, cidxs, cnum)
	if err != nil {
		return nil, err
	}

	mt, err := Init_tree_from_digests(digests)
	if err != nil {
		return nil, err
	}

	return mt, nil
}

func get_tier_no_from_n(n int64) int64 {
	ret := int64(0)
	for n > 0 {
		n = (n >> 1)
		ret++
	}
	return ret
}

func Print_tree(mt *MerkleTree) {
	for i := 0; i < int(mt.t); i++ {
		fmt.Printf("\nLevel %d\n\n", i)
		for j := 0; j < len(mt.hashes[i]); j++ {
			fmt.Printf("%x\n", mt.hashes[i][j].Data)
		}
	}
}

func Init_tree_from_digests(digests [][]byte) (*MerkleTree, error) {
	// NOTE: Currently only works when there are 2^x chunks
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

func init_tree_and_leaves(digests [][]byte, numLeafs int) (*MerkleTree, error) {
	_t := get_tier_no_from_n(int64(numLeafs))

	mt := &MerkleTree{
		hashes: make([][]*Node, _t),
		t:      _t,
		n:      numLeafs,
	}

	mt.hashes[0] = make([]*Node, numLeafs)
	for i := 0; i < numLeafs; i++ {
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
	_data := hash_childrens_data(left, right)
	return &Node{
		Data:   _data,
		leaf:   false,
		Left:   left,
		Right:  right,
		Parent: nil,
		idx:    idx,
	}
}

func hash_childrens_data(left, right *Node) []byte {
	h1 := left.Data
	h2 := right.Data
	// NOTE: we must convert to string else it SHAs the byte data
	conc := fmt.Sprintf("%x", h1) + fmt.Sprintf("%x", h2)
	hash := sha256.Sum256([]byte(conc))
	return hash[:]
}

func open_file_and_get_size(fpath string) (*os.File, int64, error) {
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	fi, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	if fi.Size() < 32 {
		fmt.Println("File is too small!")
		return nil, 0, err
	}

	return f, fi.Size(), nil
}

func gen_chunk_indexes(fsize, cnum int64) []int64 {
	chunkSize := (fsize / cnum)
	chunkIndexes := make([]int64, cnum+1)

	for i := 0; i < int(cnum); i++ {
		chunkIndexes[i] = int64(i) * chunkSize
	}

	chunkIndexes[cnum] = fsize
	return chunkIndexes
}

// From a given file, with cnum chunks at indexes cidxs, return cnum sha256sums
func gen_digests_from_file(f *os.File, cidxs []int64, cnum int64) ([][]byte, error) {
	defer f.Close()
	hashes := make([]hash.Hash, cnum)
	for i := 0; i < int(cnum); i++ {
		hashes[i] = sha256.New()
	}

	for i := 0; i < len(cidxs)-1; i++ {
		readTarget := cidxs[i+1] - cidxs[i]
		read := int64(0)

		// TODO: multithread
		for read < readTarget {
			size := min(4096, readTarget-read)
			buf := make([]byte, size)

			n, err := f.ReadAt(buf, cidxs[i]+read)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			read += int64(n)
			hashes[i].Write(buf)
		}
		//
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
