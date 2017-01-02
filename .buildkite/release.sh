#!/bin/bash
set -eu -o pipefail

VERSION="$(git describe --tags --candidates=1 2>/dev/null || echo dev)"

download_github_release() {
  wget -N https://github.com/c4milo/github-release/releases/download/v1.0.8/github-release_v1.0.8_linux_amd64.tar.gz
  tar -vxf github-release_v1.0.8_linux_amd64.tar.gz
}

github_release() {
  local version="$1"
  ./github-release lox/parfait "$version" "$BUILDKITE_COMMIT" "$(git cat-file -p "$version" | tail -n +6)" 'build/*'
}

echo "--- Downloading build artifacts"
buildkite-agent artifact download 'build/*' .

if [[ "$BUILDKITE_TAG" =~ ^v ]] ; then
  download_github_release

  echo "+++ Releasing version ${VERSION} on github.com"
  github_release "${VERSION}"
fi