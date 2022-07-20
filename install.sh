#!/bin/sh

set -au

get_qpm_download_url() {
  local OS_ARCH="$1"

  curl -sSL 'https://api.github.com/repos/anoriqq/qpm/releases/latest' | jq -r '.assets[].browser_download_url' | grep "$OS_ARCH" | head -n 1
}

install() {
  local OS_ARCH="$1"
  echo "OS_ARCH: $OS_ARCH"

  # Get qpm download URL
  NEXT_WAIT_TIME=1
  until QPM_DOWNLOAD_URL=$(get_qpm_download_url "$OS_ARCH") || [ $NEXT_WAIT_TIME -eq 8 ]; do
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

main() {
  echo "uname -a: $(uname -a)"

  ARCH=""
  case "$(uname -m)" in
    arm64) ARCH="arm64" ;;
    *) ARCH="amd64" ;;
  esac

  case "$(uname)" in
    Linux*) install "linux_${ARCH}" ;;
    Darwin*) install "darwin_${ARCH}" ;;
    *) echo "Unsupported OS" ;;
  esac
}

main

