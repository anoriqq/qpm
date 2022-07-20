#!/bin/sh

set -au

get_qpm_download_url() {
  curl -sSL 'https://api.github.com/repos/anoriqq/qpm/releases/latest' | jq -r '.assets[].browser_download_url' | grep linux_amd64
}

linux() {
  # Get qpm download URL
  NEXT_WAIT_TIME=0
  until QPM_DOWNLOAD_URL=$(get_qpm_download_url) || [ $NEXT_WAIT_TIME -eq 8 ]; do
    echo $NEXT_WAIT_TIME
    sleep $(( NEXT_WAIT_TIME++ ))
  done

  # Create tmp dir
  mkdir ./tmp.qpm

  # Download qpm to unarchive
  (cd ./tmp.qpm && curl -sSL ${QPM_DOWNLOAD_URL} | tar -zx)

  # Move Binary to PATH
  sudo mv ./tmp.qpm/qpm /usr/local/bin

  # Remove tmp dir
  rm -rf ./tmp.qpm

  # Show qpm version
  /usr/local/bin/qpm version
}

if [ "$(uname)" = "Linux" ]; then
  linux
fi

