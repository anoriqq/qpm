#!/bin/sh

set -au

get_qpm_download_url() {
  local os_arch="$1"

  curl -sSL 'https://api.github.com/repos/anoriqq/qpm/releases/latest' | jq -r '.assets[].browser_download_url' | grep "$os_arch" | head -n 1
}

install() {
  local os_arch="$1"

  # Get qpm download URL
  NEXT_WAIT_TIME=0
  until QPM_DOWNLOAD_URL=$(get_qpm_download_url "$os_arch") || [ $NEXT_WAIT_TIME -eq 8 ]; do
    echo $NEXT_WAIT_TIME
    sleep $(( NEXT_WAIT_TIME++ ))
  done
  echo "QPM_DOWNLOAD_URL: $QPM_DOWNLOAD_URL"

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

case "$(uname)" in
  Linux*) install "linux_$(uname -m)" ;;
  Darwin*) install "darwin_$(uname -m)" ;;
  *) echo "Unsupported OS" ;;
esac

