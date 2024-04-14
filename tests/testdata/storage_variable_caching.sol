// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

contract Counter1 {
    uint256 public number;
    int256[] public arr;

    function increment() public view {
        require(number < 10);
        number = number + 1;
    }

    // probably shouldnt optimise this
    function sum() public returns (int256) {
        int256 sum = 0;
        for (uint256 i = 0; i < arr.length; i++) {
            sum += arr[i];
        }
        return sum;
    }

}
