# Darakht

This library provides a somewhat robust & fairly performant* implementation to manage Merkle Trees in Go. The user-facing API provides functionality to construct a Merkle Tree from a provided file of arbitrary size, and to verify a content block's existence in said tree.

Currently, it does not provide support for sorting, which would allow us to return `false` from the Verify API sooner. This has been deliberately overlooked as sorting an arbitrarily-sized file is not viable in all scenarios. On known, or upper-bound file sizes (e.g the number of transactions in a bitcoin block), this is accomplish-able, but would be antithetical to the purpose of this project.

\* *Constructing a Merkle Tree from a 1GB file takes ~4 seconds on an i5-4278U.*

## Installation

### Library

```bash
go get github.com/masroof-maindak/darakht
```

### Executable

```bash
# [TODO]: go install?
```

## Usage

### Library

```Go
import "github.com/masroof-maindak/darakht/pkg/merkletree"

// -- Construction --
mt1, err := merkletree.InitTreeFromFile(fpath, cnum)
// [TODO]: mt2, err := merkletree.RecoverTreeFromJSON()

// -- Validation --
equal := merkletree.Equals(mt1, mt2)
// [TODO]: valid, err := merkletree.ValidateTreeJSON

// -- Proof of membership --
exists, err := merkletree.ProveMembershipFromFile(mt1, f, 4, 2)
// [TODO]: exists, err := merkletree.ProveMembershipFromJSON(fJson, f, 4, 2)
```

### Executable

```bash
# NOTE: none of these are functional at the moment

# Print the Merkle Tree of a file
darakht <file>

# Serialise a merkle tree w/ 16 leaves to `tree.json`
darakht <file> -c 16 > tree.json

# Validate whether `tree.json` holds a valid Merkle Tree
darakht -f=tree.json -validate

# Prove that bytes 4 - 14 from <file> belong in the Merkle Tree seralised in tree.json
darakht <file> -f=tree.json -prove 4 10
```

## Development

```bash
# Build
make build

# Generate random file
make rng
./rngF <fileName> <IntendedFileSize (MBs)>

# Run
./darakht <fileName> # OR: `make` to run on the sample file

# Run tests
make test
```
