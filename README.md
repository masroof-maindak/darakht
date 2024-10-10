# Darakht

This library provides a somewhat robust & fairly performant* implementation to manage Merkle Trees in Go. The user-facing API provides functionality to construct a Merkle Tree from a provided file of arbitrary size, and to verify a content block's existence in said tree.

Currently, it does not provide support for sorting, which would allow us to return `false` from the Verify API sooner. This has been deliberately overlooked as sorting an arbitrarily-sized file is not viable in all scenarios. On known, or upper-bound file sizes (e.g the number of transactions in a bitcoin block), this is accomplish-able, but would be antithetical to the purpose of this project.

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
