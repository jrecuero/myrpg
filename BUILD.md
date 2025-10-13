# Build Instructions for MyRPG

## Quick Start

### Using Makefile (Recommended)
```bash
# Build the game
make build

# Build and run
make run

# Clean build artifacts
make clean

# Show all available commands
make help
```

### Using Build Script
```bash
# Build the game
./build.sh

# Run the game
./bin/myrpg
```

### Manual Build Commands
```bash
# ✅ CORRECT - Build to bin directory
go build -o ./bin/myrpg ./cmd/myrpg

# ❌ WRONG - Creates binary in root directory
go build ./cmd/myrpg
```

## Build Options

### Development Build
```bash
make dev
# Includes race detection for debugging
```

### Release Build
```bash
make release
# Optimized binary with smaller size
```

### Clean Build
```bash
make clean
# Removes all build artifacts
```

## Binary Location

**Correct Location**: `./bin/myrpg`
- This keeps the project root clean
- Binary is properly organized
- Easy to find and run

**Incorrect Location**: `./myrpg` (root directory)
- Clutters the project root
- Can be accidentally committed to git
- Not following Go project conventions

## Running the Game

After building, run the game with:
```bash
# From project root
./bin/myrpg

# Or using make
make run
```

## Troubleshooting

### Binary in Wrong Location
If you accidentally created `myrpg` in the project root:
```bash
# Remove incorrect binary
rm myrpg

# Build correctly
make build
```

### Build Issues
```bash
# Clean and rebuild
make clean
make build

# Check Go version (requires Go 1.18+)
go version
```

### Dependencies
```bash
# Update dependencies
go mod tidy
go mod download
```

## Project Structure
```
myrpg/
├── bin/              # Build outputs (gitignored)
│   └── myrpg        # ✅ Correct binary location
├── cmd/
│   └── myrpg/
│       └── main.go  # Main entry point
├── internal/        # Game code
├── assets/          # Game assets
├── Makefile         # Build automation
├── build.sh         # Build script
└── .gitignore       # Ignores /bin/ and misplaced binaries
```

## Git Ignore

The `.gitignore` file prevents committing:
- `/bin/` directory (build outputs)
- `/myrpg` (misplaced binary)
- `*.log` files
- `.DS_Store` (macOS)

## IDE Integration

### VS Code
Add to `.vscode/tasks.json`:
```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "build",
            "type": "shell",
            "command": "make",
            "args": ["build"],
            "group": {
                "kind": "build",
                "isDefault": true
            }
        }
    ]
}
```

### GoLand/IntelliJ
Configure build configuration:
- **Program**: `go`
- **Arguments**: `build -o ./bin/myrpg ./cmd/myrpg`
- **Working directory**: `$ProjectFileDir$`

## Common Commands Summary

| Command | Purpose |
|---------|---------|
| `make build` | Build game to ./bin/myrpg |
| `make run` | Build and run game |
| `make clean` | Remove build artifacts |
| `make dev` | Development build with race detection |
| `make release` | Optimized release build |
| `./build.sh` | Alternative build script |
| `./bin/myrpg` | Run the game |

## Best Practices

1. **Always use `-o ./bin/myrpg`** when building manually
2. **Use `make build`** for consistency  
3. **Run `make clean`** before committing
4. **Never commit binaries** to git
5. **Check `.gitignore`** includes build artifacts

This ensures clean, organized builds and prevents binary files from cluttering your repository! 🎯✨