// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Base {
    uint256 public x;
}

contract Derived is Base {
    uint256 public y;
}