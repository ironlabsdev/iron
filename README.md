# Iron CLI Tool

Iron is a CLI tool that helps students of Ironlabs scaffold the project and focus on the core learning while abiding the
bootstrap process.

**Homepage**: [ironlabs.dev](https://ironlabs.dev)

# Installation

Iron CLI is available for Linux, macOS, and Windows. Choose your preferred installation method below.

## Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/ironlabsdev/iron/main/install.sh | bash
```

This script will:
- Detect your operating system and architecture
- Download the latest version
- Install to `/usr/local/bin` (with sudo) or `~/.local/bin` (without sudo)
- Automatically add to your PATH if needed

### Homebrew (macOS/Linux)

```bash
brew install ironlabsdev/tap/iron
```

This is the recommended method for macOS and Linux users who have Homebrew installed.

## Manual Installation

### Download Binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/ironlabsdev/iron/releases/latest):

| Platform | Architecture | Download |
|----------|-------------|----------|
| **Linux** | x86_64 | [iron-linux-x86_64.tar.gz](https://github.com/ironlabsdev/iron/releases/latest/download/iron-linux-x86_64.tar.gz) |
| **Linux** | ARM64 | [iron-linux-arm64.tar.gz](https://github.com/ironlabsdev/iron/releases/latest/download/iron-linux-arm64.tar.gz) |
| **macOS** | x86_64 (Intel) | [iron-mac-x86_64.tar.gz](https://github.com/ironlabsdev/iron/releases/latest/download/iron-mac-x86_64.tar.gz) |
| **macOS** | ARM64 (Apple Silicon) | [iron-mac-arm64.tar.gz](https://github.com/ironlabsdev/iron/releases/latest/download/iron-mac-arm64.tar.gz) |
| **Windows** | x86_64 | [iron-windows-x86_64.zip](https://github.com/ironlabsdev/iron/releases/latest/download/iron-windows-x86_64.zip) |
| **Windows** | ARM64 | [iron-windows-arm64.zip](https://github.com/ironlabsdev/iron/releases/latest/download/iron-windows-arm64.zip) |

### Extract and Install

#### Linux/macOS
```bash
# Download and extract (replace URL with your platform)
curl -LO https://github.com/ironlabsdev/iron/releases/latest/download/iron-linux-x86_64.tar.gz
tar -xzf iron-linux-x86_64.tar.gz

# Move to a directory in your PATH
sudo mv iron /usr/local/bin/
# or for user-only installation:
# mv iron ~/.local/bin/

# Make executable
chmod +x /usr/local/bin/iron
```

#### Windows
1. Download the `.zip` file for your architecture
2. Extract the `iron.exe` file
3. Place it in a directory that's in your PATH, or add the directory to your PATH

## Package Managers

### Linux Package Managers

#### Debian/Ubuntu (.deb)
```bash
# Download the .deb package
curl -LO https://github.com/ironlabsdev/iron/releases/latest/download/iron_linux_amd64.deb

# Install with dpkg
sudo dpkg -i iron_linux_amd64.deb

# Or install with apt (resolves dependencies)
sudo apt install ./iron_linux_amd64.deb
```

#### RedHat/CentOS/Fedora (.rpm)
```bash
# Download the .rpm package  
curl -LO https://github.com/ironlabsdev/iron/releases/latest/download/iron_linux_amd64.rpm

# Install with rpm
sudo rpm -i iron_linux_amd64.rpm

# Or install with dnf/yum (resolves dependencies)
sudo dnf install ./iron_linux_amd64.rpm
```

#### Alpine Linux (.apk)
```bash
# Download the .apk package
curl -LO https://github.com/ironlabsdev/iron/releases/latest/download/iron_linux_amd64.apk

# Install with apk
sudo apk add --allow-untrusted ./iron_linux_amd64.apk
```

## Verify Installation

After installation, verify that Iron CLI is working:

```bash
iron --version
```

You should see output similar to:
```
Iron CLI v1.0.0
Build Date: 2025-01-XX
Git Commit: abc1234
Go Version: go1.24.2
Platform: linux/amd64
```

## Getting Started

Once installed, you can start using Iron CLI:

```bash
# See available commands
iron --help

# Generate a new OAuth project
iron generate oauth my-project

# Get help for a specific command
iron generate --help
```

## Troubleshooting

### Command not found

If you get a "command not found" error after installation:

1. **Check if the binary is in your PATH:**
   ```bash
   which iron
   ```

2. **Add to PATH manually** (if using `~/.local/bin`):
   ```bash
   # For bash users
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   
   # For zsh users  
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
   source ~/.zshrc
   ```

3. **Restart your terminal** or run `source ~/.bashrc` (or `~/.zshrc`)

### Permission denied

If you get permission errors:

```bash
# Make the binary executable
chmod +x /path/to/iron

# Or reinstall with proper permissions
sudo chown $(whoami) /usr/local/bin/iron
sudo chmod +x /usr/local/bin/iron
```

### macOS Security Warning

On macOS, you might see a security warning about an unidentified developer:

1. Go to **System Preferences** → **Security & Privacy**
2. Click **"Allow Anyway"** next to the Iron CLI warning
3. Or run: `sudo xattr -rd com.apple.quarantine /usr/local/bin/iron`

## Uninstall

To remove Iron CLI:

```bash
# Remove the binary
sudo rm /usr/local/bin/iron
# or
rm ~/.local/bin/iron

# For package manager installations:
# Debian/Ubuntu: sudo apt remove iron
# RedHat/CentOS/Fedora: sudo dnf remove iron
# Alpine: sudo apk del iron
```

## Need Help?

- 📖 **Documentation**: Check our [GitHub repository](https://github.com/ironlabsdev/iron)
- 🐛 **Issues**: Report bugs on [GitHub Issues](https://github.com/ironlabsdev/iron/issues)
- 💬 **Discussions**: Join the conversation on [GitHub Discussions](https://github.com/ironlabsdev/iron/discussions)
