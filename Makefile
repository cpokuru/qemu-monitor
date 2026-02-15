.PHONY: build run clean install

# Binary name
BINARY=qemu-monitor

# Build the application
build:
	@echo "Building $(BINARY)..."
	@go build -o $(BINARY) .
	@echo "Build complete!"

# Run the application
run: build
	@echo "Starting QEMU Monitor on http://0.0.0.0:5450"
	@./$(BINARY)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY)
	@echo "Clean complete!"

# Install to /usr/local/bin (requires sudo on macOS)
install: build
	@echo "Installing to /usr/local/bin..."
	@sudo cp $(BINARY) /usr/local/bin/
	@echo "Installed! Run with: qemu-monitor"

# Uninstall from /usr/local/bin
uninstall:
	@echo "Uninstalling..."
	@sudo rm -f /usr/local/bin/$(BINARY)
	@echo "Uninstalled!"

# Help
help:
	@echo "QEMU Monitor - Makefile commands:"
	@echo ""
	@echo "  make build      - Build the application"
	@echo "  make run        - Build and run the application"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make install    - Install to /usr/local/bin (requires sudo)"
	@echo "  make uninstall  - Remove from /usr/local/bin"
	@echo "  make help       - Show this help message"
