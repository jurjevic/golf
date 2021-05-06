#!/usr/bin/env bash

if [ -z "$1" ]; then
  echo "Please provide a semantic version to create a release. e.g. 1.0.0"
  exit 1
fi

if git diff-index --quiet HEAD --; then
    echo "File changes cleaned up..."
else
    echo "Changes found! Please commit dem first"
    exit 1
fi

new_version="$1"
tag_version="v$new_version"
latest="latest"

$(go env GOPATH)/bin/golf -v main.go main.go -- '
  var NewVersion string = "'$new_version'"
'
git add main.go
git commit -m "$tag_version release build with version increment."

git push origin :refs/tags/$latest
git tag -fa $latest -m "$tag_version release build with version increment."
git tag -fa $tag_version -m "$tag_version release build with version increment."
git push origin main --tags --force

cd $(brew --repository jurjevic/homebrew-tap)
download="https://github.com/jurjevic/golf/archive/$tag_version.tar.gz"
wget $download
hash=$(sha256sum $tag_version.tar.gz)
rm "$tag_version.tar.gz"

golf Formula/golf.rb Formula/golf.rb -- '
  var HashOutput string = "'$hash'"
  var Hash string = Split(HashOutput, " ")[0]
  var Download string = "'$download'"
'

git add Formula/golf.rb
git commit -m "golf $tag_version added."
git push origin --force
