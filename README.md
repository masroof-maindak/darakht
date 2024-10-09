# Merkle Tree

The relevant code files are in `cmd/test/main.go` & `pkg/merkletree/merkletree.go`.

### Usage

```bash
# Build
make build

# Generate a random file
make rng
./rngF <fileName> <IntendedFileSize (MBs)>

# Run
./merkle <fileName> # OR: `make` to run on the sample file
```
