#!/bin/bash
# build.sh - Simple build script for MyRPG

set -e  # Exit on any error

echo "🎮 Building MyRPG..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the game
go build -o ./bin/myrpg ./cmd/myrpg

echo "✅ Build complete! Binary created at: ./bin/myrpg"
echo ""
echo "To run the game:"
echo "  ./bin/myrpg"
echo ""
echo "Or use the Makefile:"
echo "  make run"