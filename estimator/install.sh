#!/bin/bash
curl -L https://foundry.paradigm.xyz | bash
# if terminal is zsh, source .zshenv
if [ -n "$ZSH_VERSION" ]; then
    echo "zsh detected"
    source ~/.zshenv
fi
if [ -n "$BASH_VERSION" ]; then
    echo "bash detected"
    source ~/.bashrc
fi

# check if the foundryup is installed
# using which
if [ -x "$(command -v foundryup)" ]; then
    echo "foundryup is installed"
else
    echo "foundryup is not installed"
    exit 1
fi

# install forge
foundryup

# check if forge is installed
if [ -x "$(command -v forge)" ]; then
    echo "forge is installed"
else
    echo "forge is not installed"
    exit 1
fi

# run forge init
forge init .
