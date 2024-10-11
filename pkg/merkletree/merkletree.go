package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"os"
)

const NUM_CHUNKS int64 = 8

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
	idx    int // index of this node within its tier
}

func ProveMember(mt *MerkleTree, f *os.File, idx, run int64) (bool, error) {
	h, err := hashFileChunk(f, idx, run)
	if err != nil {
		return false, err
	}

	return ProveDigest(mt, h.Sum(nil)), nil
}

// Receives a leaf digest and returns whether it's a member of the tree or not
func ProveDigest(mt *MerkleTree, hash []byte) bool {
	idx := findLeafIndex(mt, hash)
	if idx == -1 {
		return false
	}

	var li, ri, pi int
	t := 0

	for t < mt.t-1 {
		if idx%2 == 0 {
			li = idx
			ri = idx + 1
		} else {
			li = idx - 1
			ri = idx
		}
		pi = li / 2

		hd := hashChildrensData(mt.hashes[t][li], mt.hashes[t][ri])
		if !bytes.Equal(hd, mt.hashes[t+1][pi].Data) {
			return false
		}

		idx = pi
		t++
	}

	return true
}

func findLeafIndex(mt *MerkleTree, target []byte) int {
	for i := 0; i < mt.n; i++ {
		if bytes.Equal(target, mt.hashes[0][i].Data) {
			return i
		}
	}
	return -1
}

func InitTreeFromFile(f *os.File, cnum int64) (*MerkleTree, error) {
	if cnum <= 0 {
		return nil, errors.New("invalid leaf count")
	}

	fsize, err := getFileSize(f, cnum)
	if err != nil {
		return nil, err
	}

	cidxs := genChunkIndexes(fsize, cnum)
	digests, err := genDigestsFromFile(f, cidxs, cnum)
	if err != nil {
		return nil, err
	}

	mt, err := initTreeFromDigests(digests)
	if err != nil {
		return nil, err
	}

	return mt, nil
}

func getTierFromN(n int64) int {
	ret := 0
	for n > 0 {
		n = (n >> 1)
		ret++
	}
	return ret
}

func PrintTree(mt *MerkleTree) {
	depth, width := mt.t, mt.n

	fmt.Println("{")
	fmt.Printf("\t\"root\": \"%x\",\n", MerkleRoot(mt))
	fmt.Printf("\t\"t\": %d,\n", depth)
	fmt.Printf("\t\"n\": %d,\n", width)
	fmt.Printf("\t\"hashes\": [\n")

	for t := 0; t < depth; t, width = t+1, width/2 {
		fmt.Println("\t\t[")

		for i := 0; i < width; i++ {
			fmt.Printf("\t\t\t\"%x\"", mt.hashes[t][i].Data)
			if i != width-1 {
				fmt.Println(",")
			}
		}

		fmt.Print("\n\t\t]")
		if t != depth-1 {
			fmt.Println(",")
		}
	}

	fmt.Println("\n\t]")
	fmt.Println("}")
}

func initTreeFromDigests(digests [][]byte) (*MerkleTree, error) {
	nodeCount := len(digests)
	if nodeCount > 0 && nodeCount&(nodeCount-1) != 0 {
		return nil, errors.New("no. of digests is not a power of 2")
	}

	mt := initTreeAndLeaves(digests, nodeCount)
	nodeCount /= 2
	t := 1

	for nodeCount > 0 {
		mt.hashes[t] = make([]*Node, nodeCount)
		for i := 0; i < nodeCount; i++ {
			mt.hashes[t][i] = createParent(i, mt.hashes[t-1][2*i], mt.hashes[t-1][2*i+1])
			mt.hashes[t-1][2*i].Parent = mt.hashes[t][i]
			mt.hashes[t-1][2*i+1].Parent = mt.hashes[t][i]
		}
		nodeCount /= 2
		t++
	}

	return mt, nil
}

func initTreeAndLeaves(digests [][]byte, n int) *MerkleTree {
	t := getTierFromN(int64(n))

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

	return mt
}

func createParent(idx int, left, right *Node) *Node {
	d := hashChildrensData(left, right)
	return &Node{
		Data:   d,
		leaf:   false,
		Left:   left,
		Right:  right,
		Parent: nil,
		idx:    idx,
	}
}

func hashChildrensData(left, right *Node) []byte {
	// NOTE: we must convert to string else it SHAs the byte data
	conc := fmt.Sprintf("%x", left.Data) + fmt.Sprintf("%x", right.Data)
	hash := sha256.Sum256([]byte(conc))
	return hash[:]
}

func getFileSize(f *os.File, cnum int64) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}

	if fi.Size() < cnum {
		return 0, errors.New("file is too small to be chunked")
	}

	return fi.Size(), nil
}

func genChunkIndexes(fsize, cnum int64) []int64 {
	csize := (fsize / cnum)
	cidxs := make([]int64, cnum+1)

	for i := int64(0); i < cnum; i++ {
		cidxs[i] = int64(i) * csize
	}

	cidxs[cnum] = fsize
	return cidxs
}

// Write 'length' bytes, starting at the 'start'-th byte, in file 'f' into a hash object
func hashFileChunk(f *os.File, start, length int64) (hash.Hash, error) {
	h := sha256.New()
	read := int64(0)

	for read < length {
		size := min(4096, length-read)
		buf := make([]byte, size)

		n, err := f.ReadAt(buf, start+read)
		if err != nil {
			return nil, err
		}

		read += int64(n)
		h.Write(buf)
	}

	return h, nil
}

// From a given file, with cnum chunks at indexes cidxs, return cnum sha256sums
func genDigestsFromFile(f *os.File, cidxs []int64, cnum int64) ([][]byte, error) {
	hashes := make([]hash.Hash, cnum)
	for i := int64(0); i < cnum; i++ {
		hashes[i] = sha256.New()
	}

	for i := 0; i < len(cidxs)-1; i++ {
		// TODO(?): multithread
		h, err := hashFileChunk(f, cidxs[i], cidxs[i+1]-cidxs[i])
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
