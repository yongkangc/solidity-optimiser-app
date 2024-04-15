// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract OptimizationShowcase {
    // Unoptimized struct
    struct UnoptimizedProduct {
        uint256 id;         // 32 bytes
        bool isAvailable;   // 1 byte
        uint256 price;      // 32 bytes
        uint32 quantity;    // 4 bytes
        string name;        // dynamic size
        uint32 category;    // 4 bytes
        address seller;     // 20 bytes
        uint16 ratings;     // 2 bytes
    }

    // // Optimized struct
    // struct OptimizedProduct {
    //     uint256 id;         // 32 bytes
    //     uint256 price;      // 32 bytes
    //     address seller;     // 20 bytes
    //     uint32 quantity;    // 4 bytes
    //     uint32 category;    // 4 bytes
    //     uint16 ratings;     // 2 bytes
    //     bool isAvailable;   // 1 byte
    //     string name;        // dynamic size
    // }

    uint256 public variable1;
    uint256 public variable2;

    // Unoptimized function
    function calculateSumUnoptimized() public view returns (uint256) {
        uint256 sum = variable1 + variable2;
        return sum;
    }

    // should not be optimized as variable is only read once
    function calculateSumOptimized() public view returns (uint256) {
        uint256 v1 = variable1;
        uint256 v2 = variable2;
        uint256 sum = v1 + v2;
        return sum;
    }

    // Unoptimized function with memory array
    function sumOfArrayUnoptimized(uint256[] memory numbers) external pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }

    // // Optimized function with calldata array
    function sumOfArrayOptimized(uint256[] calldata numbers) external pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }
}
