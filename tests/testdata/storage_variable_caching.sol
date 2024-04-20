// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Counter1 {
    uint256 public number;
    int256[] public arr;

    function increment() public view returns (uint256) {
        require(number < 10);
        uint256 incremented = number + 1;
        return incremented;
    }

    function sum() public view returns (int256) {
        int256 sum = 0;
        for (uint256 i = 0; i < arr.length; i++) {
            sum += arr[i];
        }
        return sum;
    }

}
