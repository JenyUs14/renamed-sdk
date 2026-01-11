#!/bin/bash
set -e

VERSION=$(node -p "require('./package.json').version")
TAG="latest"

if [[ "$VERSION" == *"-"* ]]; then
  # Prerelease version (contains hyphen: beta, alpha, rc, etc.)
  TAG=$(echo "$VERSION" | sed 's/.*-\([a-zA-Z]*\).*/\1/')
fi

echo "Publishing version $VERSION with tag $TAG"

OUTPUT=$(npm publish --provenance --access public --tag "$TAG" 2>&1) || {
  EXIT_CODE=$?
  echo "$OUTPUT"
  if echo "$OUTPUT" | grep -q "cannot publish over the previously published"; then
    echo "Package already exists, skipping..."
    exit 0
  fi
  exit $EXIT_CODE
}

echo "$OUTPUT"
