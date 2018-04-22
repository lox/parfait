#!/bin/bash
set -euo pipefail

export VERSION=$(awk -F\" '/const Version/ {print $2}' version/version.go)

echo "Checking if $VERSION is a tag..."

# If there is already a release (which means a tag), we want to avoid trying to create
# another one, as this will fail and cause partial broken releases
# If there is already a release (which means a tag), we want to avoid trying to create
# another one, as this will fail and cause partial broken releases
if git ls-remote --tags origin | grep "refs/tags/v${VERSION}" ; then
  echo "Tag refs/tags/v${VERSION} already exists"
  exit 0
fi

buildkite-agent meta-data set release-version "$VERSION"
buildkite-agent pipeline upload .buildkite/pipeline.release.yml
