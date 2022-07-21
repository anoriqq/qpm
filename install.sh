#!/bin/sh

set -au

get_qpm_download_url() {
  local OS_ARCH="$1"

  if [ "${CI:-}" = "true" ]; then
    curl -sSL https://api.github.com/repos/anoriqq/qpm/releases/latest \
      | jq -r '.assets[].browser_download_url' \
      | grep "$OS_ARCH" \
      | head -n 1
  else
    curl -sSL -H 'authorization: Bearer ${{ secrets.GITHUB_TOKEN }}' https://api.github.com/repos/anoriqq/qpm/releases/latest \
      | jq -r '.assets[].browser_download_url' \
      | grep "$OS_ARCH" \
      | head -n 1
  fi
}

install() {
  local OS_ARCH="$1"
  echo "OS_ARCH: $OS_ARCH"

  # Get qpm download URL
  local NEXT_WAIT_TIME=1
  until QPM_DOWNLOAD_URL=$(get_qpm_download_url "$OS_ARCH") || [ "$NEXT_WAIT_TIME" -gt 10 ]; do
    echo $NEXT_WAIT_TIME
    sleep $(( NEXT_WAIT_TIME ))
    NEXT_WAIT_TIME=$(( NEXT_WAIT_TIME * 2 ))
  done
  echo "QPM_DOWNLOAD_URL: $QPM_DOWNLOAD_URL"
  if [ "$QPM_DOWNLOAD_URL" = "" ]; then
    exit 1
  fi

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

