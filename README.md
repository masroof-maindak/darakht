# Merkle Tree

This library provides a somewhat robust & fairly performant* implementation to manage Merkle Trees in Go. The user-facing API provides functionality to construct a Merkle Tree from a provided file of arbitrary size, and to verify a content block's existence in said tree.

Further improvements involve allowing the user to define their own hash function to be used during the construction of the tree; accomplishing this would be trivial by virtue of using function pointers.

Another possible optimisation, sorting, would allow us to return false  from the Verify API sooner, but has been deliberately overlooked as sorting an arbitrarily-sized file is not viable in all scenarios, and is a non-complex addition anyway; we merely need a copy of the input file where the chunks are sorted instead of being in their original positions.

On known, or upper-bound file sizes, e.g the number of transactions in a bitcoin block, this is accomplish-able, but would be antithetical to the goal of this project; i.e providing the means to construct Merkle Tree and verify the existence of a data block from an arbitrary-sized file.

\**Constructing a Merkle Tree from a 1GB file takes ~4 seconds on an i5-4278U.*

### Usage

```bash
# install
go get github.com/masroof-maindak/darakht
```

```Go
// import
import "github.com/masroof-maindak/darakht/pkg/merkletree"
```

### Development

```bash
# Build
make build

# Generate random file
make rng
./rngF <fileName> <IntendedFileSize (MBs)>

# Run
./merkle <fileName> # OR: `make` to run on the sample file

# Run tests
make test
```
