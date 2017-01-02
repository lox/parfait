#!/bin/bash
set -eu -o pipefail

VERSION="$(git describe --tags --candidates=1)"

download_github_release() {
  wget -N https://github.com/c4milo/github-release/releases/download/v1.0.8/github-release_v1.0.8_linux_amd64.tar.gz
  tar -vxf github-release_v1.0.8_linux_amd64.tar.gz
}

github_release() {
  local version="$1"
  ./github-release cardigann/cardigann "$version" "$TRAVIS_COMMIT" "$(git cat-file -p "$version" | tail -n +6)" ""
}

if [[ "$BUILDKITE_TAG" =~ ^v ]] ; then
  download_github_release

  echo "Releasing version ${VERSION} on github.com"
  github_release "${VERSION}"
fi