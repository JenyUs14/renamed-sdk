#!/bin/bash
set -e

# Creates a path-prefixed tag for Go module in subdirectory
# Go modules in subdirectories require tags like sdks/go/v1.0.0

VERSION=${GITHUB_REF#refs/tags/}
GO_TAG="sdks/go/${VERSION}"

echo "Creating Go module tag: $GO_TAG"

if git rev-parse "$GO_TAG" >/dev/null 2>&1; then
  echo "Tag $GO_TAG already exists, skipping"
else
  git tag "$GO_TAG"
  git push origin "$GO_TAG"
  echo "Created and pushed tag: $GO_TAG"
fi

echo "Go module available: go get github.com/${GITHUB_REPOSITORY}/sdks/go@${VERSION}"
