#!/bin/bash
set -eu -o pipefail

version=$(buildkite-agent meta-data get release-version)

download_github_release() {
  wget -N https://github.com/buildkite/github-release/releases/download/v1.0/github-release-linux-amd64
  mv github-release-linux-amd64 github-release
  chmod +x ./github-release
}

github_release() {
  local version="$1"
  ./github-release "$version" build/* --commit "${BUILDKITE_COMMIT}" \
                                      --tag "v${version}" \
                                      --github-repository "lox/parfait"
}

download_github_release

echo "--- Downloading build artifacts"
buildkite-agent artifact download 'build/*' .

echo "+++ Releasing version v${version} on github.com"
github_release "v${version}"
