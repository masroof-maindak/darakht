package merkletree

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
)

const (
	NUM_CHUNKS           int64  = 8
	SHA256_STRING_LEN    int    = 64
	INVALID_HASH_ERR     string = "Invalid hash found!"
	MISMATCHED_ROOTS_ERR string = "Roots in JSON file don't match"
)

type MerkleTreeTemp struct {
	Root   string     `json:"root"`
	T      int        `json:"t"`
	N      int        `json:"n"`
	Hashes [][]string `json:"hashes"`
}

type MerkleTree struct {
	hashes [][]*Node // all hashes
	t      int       // no. of tiers
	n      int       // no. of leaves
}

type Node struct {
	data  []byte
	left  *Node
	right *Node
	idx   int // index of this node within its tier
}

func (mt *MerkleTree) ProveMember(f *os.File, idx, run int64) (bool, error) {
	h, err := hashFileChunk(f, idx, run)
	if err != nil {
		return false, err
	}

	return mt.ProveDigest(h.Sum(nil)), nil
}

func (mt1 *MerkleTree) Equals(mt2 *MerkleTree) bool {
	if mt1 == mt2 {
		return true
	}

	if mt1 == nil || mt2 == nil {
		return false
	}

	if !(mt1.n == mt2.n && mt1.t == mt2.t) {
		return false
	}

	depth := mt1.t
	width := mt1.n

	for t := 0; t < depth; t, width = t+1, width/2 {
		for i := 0; i < width; i++ {
			n1 := mt1.hashes[t][i]
			n2 := mt2.hashes[t][i]
			if !(bytes.Equal(n1.data, n2.data) && n1.idx == n2.idx) {
				return false
			}
		}
	}

	return true
}

func (mt *MerkleTree) findLeafIndex(target []byte) int {
	for i := 0; i < mt.n; i++ {
		if bytes.Equal(target, mt.hashes[0][i].data) {
			return i
		}
	}
	return -1
}

// Receives a leaf digest and returns whether it's a member of the tree or not
func (mt *MerkleTree) ProveDigest(hash []byte) bool {
	idx := mt.findLeafIndex(hash)
	if idx == -1 {
		return false
	}

	for t := 0; t < mt.t-1; t++ {
		var li, ri, pi int

		if idx%2 == 0 {
			li = idx
			ri = idx + 1
		} else {
			li = idx - 1
			ri = idx
		}
		pi = li / 2

		hd := hashChildrensData(mt.hashes[t][li], mt.hashes[t][ri])
		if !bytes.Equal(hd, mt.hashes[t+1][pi].data) {
			return false
		}

		idx = pi
	}

	return true
}

// Checks if a JSON file holds a Merkle Tree. Destroys file cursor.
func deserialiseJSONFile(fJson *os.File) (*MerkleTreeTemp, error) {
	mtt := &MerkleTreeTemp{}

	// CHECK - can we do this without wasting the cursor? Does the prior even matter?
	fJson.Seek(0, io.SeekStart)
	// TODO(?): buffered reading
	b, err := io.ReadAll(fJson)
	if err != nil {
		return nil, err
	}

	if !json.Valid(b) {
		return nil, err
	}

	err = json.Unmarshal(b, &mtt)
	if err != nil {
		return nil, err
	}

	return mtt, nil
}

func validateTempRoots(jsonR1, jsonR2 string) ([]byte, error) {
	byteR1, err := hex.DecodeString(jsonR1)
	if err != nil {
		return nil, errors.New(INVALID_HASH_ERR)
	}

	byteR2, err := hex.DecodeString(jsonR2)
	if err != nil {
		return nil, errors.New(INVALID_HASH_ERR)
	}

	if !bytes.Equal(byteR1, byteR2) {
		return nil, errors.New(MISMATCHED_ROOTS_ERR)
	}

	return byteR1, nil
}

func validateMerkleTreeTemp(mtt *MerkleTreeTemp) bool {
	if mtt.T == 0 || mtt.N == 0 {
		return false
	}

	depth := mtt.T
	width := mtt.N

	if len(mtt.Root) != SHA256_STRING_LEN || len(mtt.Hashes) != depth {
		return false
	}

	for t := 0; t < depth; t, width = t+1, width/2 {
		if len(mtt.Hashes[t]) != width {
			return false
		}
		for i := 0; i < width; i++ {
			if len(mtt.Hashes[t][i]) != SHA256_STRING_LEN {
				return false
			}
		}
	}

	return true
}

func createMerkleTreeFromTemp(mtt *MerkleTreeTemp) (*MerkleTree, error) {
	if !validateMerkleTreeTemp(mtt) {
		return nil, errors.New("Invalid merkle tree deserialised!")
	}

	t, n, leaves := mtt.T, mtt.N, mtt.Hashes[0]

	_, err := validateTempRoots(mtt.Root, mtt.Hashes[t-1][0])
	if err != nil {
		return nil, err
	}

	digests := make([][]byte, int(n))
	for i, leaf := range leaves {
		strLeaf, err := hex.DecodeString(leaf)
		if err != nil {
			return nil, errors.New(INVALID_HASH_ERR)
		}
		digests[i] = []byte(strLeaf)
	}

	mt, err := createTreeFromDigests(digests)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i += 2 {
		if !mt.ProveDigest(digests[i]) {
			return nil, errors.New("Merkle Tree was tampered with!")
		}
	}

	return mt, nil
}

func NewMerkleTreeFromJSON(fJson *os.File) (*MerkleTree, error) {
	mtt, err := deserialiseJSONFile(fJson)
	if err != nil {
		return nil, err
	}

	mt, err := createMerkleTreeFromTemp(mtt)
	if err != nil {
		return nil, err
	}

	return mt, nil
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

func NewMerkleTreeFromFile(f *os.File, cnum int64) (*MerkleTree, error) {
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

	mt, err := createTreeFromDigests(digests)
	if err != nil {
		return nil, err
	}

	return mt, nil
}

// Write Merkle Tree to file; f must be accessed with os.Create
func (mt *MerkleTree) Serialise(f *os.File) error {
	depth, width := mt.t, mt.n

	if depth == 0 {
		return errors.New("Empty Merkle Tree: nothing to serialise")
	}

	w := bufio.NewWriter(f)

	fmt.Fprintf(w, "{\n")
	fmt.Fprintf(w, "\t\"root\": \"%x\",\n", mt.MerkleRoot())
	fmt.Fprintf(w, "\t\"t\": %d,\n", depth)
	fmt.Fprintf(w, "\t\"n\": %d,\n", width)
	fmt.Fprintf(w, "\t\"hashes\": [\n")

	for t := 0; t < depth; t, width = t+1, width/2 {
		fmt.Fprintf(w, "\t\t[\n")

		for i := 0; i < width; i++ {
			fmt.Fprintf(w, "\t\t\t\"%x\"", mt.hashes[t][i].data)
			if i != width-1 {
				fmt.Fprintf(w, ",\n")
			}
		}

		fmt.Fprintf(w, "\n\t\t]")
		if t != depth-1 {
			fmt.Fprintf(w, ",\n")
		}
	}

	fmt.Fprintf(w, "\n\t]\n")
	fmt.Fprintf(w, "}\n")

	return w.Flush()
}

// Print Merkle Tree to stdout
func (mt *MerkleTree) Print() error {
	return mt.Serialise(os.Stdout)
}

func hashChildrensData(left, right *Node) []byte {
	// We have to sha the hexadecimal strings of each digest
	conc := fmt.Sprintf("%x", left.data) + fmt.Sprintf("%x", right.data)
	hash := sha256.Sum256([]byte(conc))
	return hash[:]
}

func createParent(idx int, left, right *Node) *Node {
	d := hashChildrensData(left, right)
	return &Node{
		data:  d,
		left:  left,
		right: right,
		idx:   idx,
	}
}

// Create and return full tree object from digests
func createTreeFromDigests(digests [][]byte) (*MerkleTree, error) {
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
		}
		nodeCount /= 2
		t++
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

// Init tree object and copy leaves over
func initTreeAndLeaves(digests [][]byte, n int) *MerkleTree {
	_t := getTierFromN(int64(n))

	mt := &MerkleTree{
		hashes: make([][]*Node, _t),
		t:      _t,
		n:      n,
	}

	mt.hashes[0] = make([]*Node, n)
	for i := 0; i < n; i++ {
		mt.hashes[0][i] = &Node{
			data:  digests[i],
			left:  nil,
			right: nil,
			idx:   i,
		}
	}

	return mt
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

func (mt *MerkleTree) MerkleRoot() []byte {
	return mt.hashes[mt.t-1][0].data
}
