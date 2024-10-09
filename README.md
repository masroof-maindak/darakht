# Merkle Tree

The relevant code files are in `cmd/test/main.go` & `pkg/merkletree/merkletree.go`.

### Usage

```bash
# Build
make build

# Generate a random file
make rng
./generator <file_name> <intended_file_size (in MBs)>

# Run
./merkle <file_name> # OR: `make` to run on the sample file
```
