// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ArraySumCalculator {

    // This pure function calculates the sum of an array of integers.
    // The array is passed as calldata to optimize gas usage.
    function sumOfArray(uint256[] memory numbers) external pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }

        // This pure function calculates the sum of an array of integers.
    // The array is passed as calldata to optimize gas usage.
    function sumOfArrayOptimized(uint256[] calldata numbers) external pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }

    function shouldntOptimiseThis(uint256[] memory numbers) returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
            numbers[i] = 0;
        }
        return sum;
    }
}
