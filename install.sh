#!/bin/sh

# Iron CLI Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/ironlabsdev/iron/main/install.sh | bash
# Or with custom version: curl -fsSL https://raw.githubusercontent.com/ironlabsdev/iron/main/install.sh | VERSION=v1.0.0 bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="ironlabsdev/iron"
BINARY_NAME="iron"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
detect_os_arch() {
    OS=$(uname -s)
    ARCH=$(uname -m)

    case $OS in
        Linux*)
            OS="linux"
            ;;
        Darwin*)
            OS="mac"
            ;;
        CYGWIN*|MINGW*|MSYS*)
            OS="windows"
            ;;
        *)
            echo "${RED}Unsupported operating system: $OS${NC}"
            exit 1
            ;;
    esac

    case $ARCH in
        x86_64|amd64)
            ARCH="x86_64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="armv7"
            ;;
        *)
            echo "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
}

# Get latest version from GitHub releases
get_latest_version() {
    if [ -n "${VERSION:-}" ]; then
        echo "Using specified version: $VERSION"
        return
    fi

    echo "${BLUE}Fetching latest version...${NC}"
    VERSION=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

    if [ -z "$VERSION" ]; then
        echo "${RED}Failed to get latest version${NC}"
        exit 1
    fi

    echo "${GREEN}Latest version: $VERSION${NC}"
}

# Check if running as root for system-wide installation
check_permissions() {
    if [ "$EUID" -ne 0 ] && [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
        echo "${YELLOW}Note: Installing to /usr/local/bin requires sudo privileges${NC}"
        INSTALL_DIR="$HOME/.local/bin"
        echo "${BLUE}Installing to $INSTALL_DIR instead${NC}"
        mkdir -p "$INSTALL_DIR"

        # Add to PATH if not already there
        SHELL_CONFIG=""
        if [ -n "${ZSH_VERSION:-}" ]; then
            SHELL_CONFIG="$HOME/.zshrc"
        elif [ -n "${BASH_VERSION:-}" ]; then
            SHELL_CONFIG="$HOME/.bashrc"
        fi

        if [ -n "$SHELL_CONFIG" ] && [ -f "$SHELL_CONFIG" ]; then
            if ! grep -q "$INSTALL_DIR" "$SHELL_CONFIG"; then
                echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$SHELL_CONFIG"
                echo "${YELLOW}Added $INSTALL_DIR to PATH in $SHELL_CONFIG${NC}"
                echo "${YELLOW}Please run 'source $SHELL_CONFIG' or restart your shell${NC}"
            fi
        fi
    fi
}

# Download and install binary
install_binary() {
    # Construct download URL
    FILENAME="${BINARY_NAME}-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        FILENAME="${FILENAME}.zip"
        ARCHIVE_TYPE="zip"
    else
        FILENAME="${FILENAME}.tar.gz"
        ARCHIVE_TYPE="tar.gz"
    fi

    DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$VERSION/$FILENAME"

    echo "${BLUE}Downloading $DOWNLOAD_URL...${NC}"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    # Download archive
    if command -v curl >/dev/null 2>&1; then
        curl -sLO "$DOWNLOAD_URL"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$DOWNLOAD_URL"
    else
        echo "${RED}Error: curl or wget is required${NC}"
        exit 1
    fi

    # Extract archive
    echo "${BLUE}Extracting archive...${NC}"
    if [ "$ARCHIVE_TYPE" = "zip" ]; then
        unzip -q "$FILENAME"
    else
        tar -xzf "$FILENAME"
    fi

    # Find the binary (it might be in a subdirectory)
    BINARY_PATH=$(find . -name "$BINARY_NAME" -type f | head -1)
    if [ -z "$BINARY_PATH" ]; then
        BINARY_PATH=$(find . -name "${BINARY_NAME}.exe" -type f | head -1)
    fi

    if [ -z "$BINARY_PATH" ]; then
        echo "${RED}Error: Binary not found in archive${NC}"
        exit 1
    fi

    # Install binary
    echo "${BLUE}Installing to $INSTALL_DIR...${NC}"

    if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
        sudo mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi

    # Cleanup
    cd - >/dev/null
    rm -rf "$TMP_DIR"
}

# Verify installation
verify_installation() {
    echo "${BLUE}Verifying installation...${NC}"

    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        VERSION_OUTPUT=$("$BINARY_NAME" --version 2>/dev/null || "$BINARY_NAME" version 2>/dev/null || echo "unknown")
        echo "${GREEN}✅ Iron CLI installed successfully!${NC}"
        echo "${GREEN}Version: $VERSION_OUTPUT${NC}"
        echo "${BLUE}Run 'iron --help' to get started${NC}"
    else
        echo "${YELLOW}⚠️  Installation completed, but 'iron' is not in your PATH${NC}"
        echo "${BLUE}You can run it directly: $INSTALL_DIR/$BINARY_NAME${NC}"

        if [ "$INSTALL_DIR" != "/usr/local/bin" ]; then
            echo "${YELLOW}Please restart your shell or run: export PATH=\"$INSTALL_DIR:\$PATH\"${NC}"
        fi
    fi
}

# Main installation flow
main() {
    echo "${GREEN}Iron CLI Installer${NC}"
    echo "${BLUE}==================${NC}"

    detect_os_arch
    echo "${BLUE}Detected OS: $OS, Architecture: $ARCH${NC}"

    get_latest_version
    check_permissions
    install_binary
    verify_installation

    echo "${GREEN}🎉 Installation complete! $${NC}"
}

# Check if required tools are available
check_requirements() {
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        echo "${RED}Error: curl or wget is required for installation${NC}"
        exit 1
    fi

    if ! command -v tar >/dev/null 2>&1; then
        echo "${RED}Error: tar is required for installation${NC}"
        exit 1
    fi
}

# Run main function
check_requirements
main