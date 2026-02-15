#!/bin/bash

# QEMU Monitor - Quick Start Script

echo "================================"
echo "   QEMU Instance Monitor"
echo "================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo "✓ Go detected: $(go version)"
echo ""

# Build the application
echo "📦 Building application..."
if go build -o qemu-monitor .; then
    echo "✓ Build successful!"
else
    echo "❌ Build failed"
    exit 1
fi

echo ""
echo "🚀 Starting QEMU Monitor..."
echo ""
echo "   Web UI: http://localhost:5450"
echo "   API:    http://localhost:5450/api/instances"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Run the application
./qemu-monitor
