.PHONY: build clean

# Binary name
BINARY_NAME=cmdeagle

# Build directory
BUILD_DIR=bin

# Go build flags
GOBUILD=go build -o $(BUILD_DIR)/$(BINARY_NAME)

# Ensure build directory exists
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)
	chmod 755 $(BUILD_DIR)

# Build the binary
build: $(BUILD_DIR)
	$(GOBUILD)
	chmod 755 $(BUILD_DIR)/$(BINARY_NAME)
	
# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR) 


release:
	git tag -a $(VERSION) -m "$(VERSION)"
	git push origin $(VERSION)
	git push --tags
	goreleaser release --clean

test:
	go test ./...

