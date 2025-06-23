#!/bin/bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default increment type
INCREMENT="patch"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --minor)
      INCREMENT="minor"
      shift
      ;;
    --major)
      INCREMENT="major"
      shift
      ;;
    --patch)
      INCREMENT="patch"
      shift
      ;;
    -h|--help)
      echo "Usage: $0 [--patch|--minor|--major]"
      echo "  --patch  Increment patch version (default)"
      echo "  --minor  Increment minor version"
      echo "  --major  Increment major version"
      exit 0
      ;;
    *)
      echo "Unknown option $1"
      exit 1
      ;;
  esac
done

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Check if working directory is clean
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: Working directory is not clean. Please commit or stash your changes.${NC}"
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${RED}Error: You must be on the main branch to create a release${NC}"
    exit 1
fi

# Pull latest changes
echo -e "${BLUE}Pulling latest changes...${NC}"
git pull origin main

# Get the latest tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo -e "${BLUE}Latest tag: ${LATEST_TAG}${NC}"

# Remove 'v' prefix for version calculation
VERSION=${LATEST_TAG#v}

# Split version into components
IFS='.' read -ra VERSION_PARTS <<< "$VERSION"
MAJOR=${VERSION_PARTS[0]:-0}
MINOR=${VERSION_PARTS[1]:-0}
PATCH=${VERSION_PARTS[2]:-0}

# Increment version based on type
case $INCREMENT in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
esac

NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"

echo -e "${YELLOW}Creating new ${INCREMENT} release: ${NEW_VERSION}${NC}"

# Confirm with user
read -p "Do you want to proceed with release ${NEW_VERSION}? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Release cancelled${NC}"
    exit 0
fi

# Run linter
echo -e "${BLUE}Running linter...${NC}"
if command -v golangci-lint &> /dev/null; then
    if ! golangci-lint run; then
        echo -e "${RED}Linting failed. Aborting release.${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}Warning: golangci-lint not found. Skipping linting.${NC}"
fi

# Create and push tag
echo -e "${BLUE}Creating and pushing tag ${NEW_VERSION}...${NC}"
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"
git push origin "$NEW_VERSION"

echo -e "${GREEN}✅ Successfully created release ${NEW_VERSION}${NC}"
echo -e "${BLUE}GitHub Actions will now build and publish the release.${NC}"
echo -e "${BLUE}Check the progress at: https://github.com/ironlabsdev/iron/actions${NC}"