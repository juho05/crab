#!/bin/bash

if ! command -v go &> /dev/null
then
    echo "Go is not installed! Please install Go version 1.18 or above."
    exit
fi

echo "Installing 'crab' into '/usr/local/bin'..."
sudo GOBIN=/usr/local/bin GOPRIVATE=github.com/Bananenpro/crab go install github.com/Bananenpro/crab@latest
