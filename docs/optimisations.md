# Example Solidity Optimisation

```solidity
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

    // Optimized struct
    struct OptimizedProduct {
        uint256 id;         // 32 bytes
        uint256 price;      // 32 bytes
        address seller;     // 20 bytes
        uint32 quantity;    // 4 bytes
        uint32 category;    // 4 bytes
        uint16 ratings;     // 2 bytes
        bool isAvailable;   // 1 byte
        string name;        // dynamic size
    }

    uint256 public variable1;
    uint256 public variable2;

    // Unoptimized function
    function calculateSumUnoptimized() public view returns (uint256) {
        uint256 sum = variable1 + variable2;
        return sum;
    }

    // Optimized function with cached storage variables
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

    // Optimized function with calldata array
    function sumOfArrayOptimized(uint256[] calldata numbers) external pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
        }
        return sum;
    }

    // Function that shouldn't be optimized with calldata
    function shouldntOptimizeThis(uint256[] memory numbers) public pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < numbers.length; ++i) {
            sum += numbers[i];
            numbers[i] = 0;
        }
        return sum;
    }
}
```

In this Solidity file:

1. Struct Data Packing:

   - The `UnoptimizedProduct` struct represents the unoptimized version with suboptimal storage layout.
   - The `OptimizedProduct` struct demonstrates the optimized version where variables are arranged to minimize storage gaps.

2. Caching Storage Variables:

   - The `calculateSumUnoptimized` function shows the unoptimized version that directly accesses the storage variables `variable1` and `variable2`.
   - The `calculateSumOptimized` function demonstrates the optimized version where the storage variables are cached in local variables `v1` and `v2` before performing the calculation.

3. Calldata Efficiency:
   - The `sumOfArrayUnoptimized` function takes the array as `memory`, which is less gas-efficient.
   - The `sumOfArrayOptimized` function takes the array as `calldata`, which is more gas-efficient when the array is not modified.
   - The `shouldntOptimizeThis` function intentionally uses `memory` because it modifies the array, making it unsuitable for `calldata` optimization.

This file provides a clear comparison between the unoptimized and optimized versions of each optimization technique, allowing you to see the differences and understand how the optimizations can be applied in practice.

## Struct Packing Optimization in Solidity

When working with structs in Solidity, it's important to consider the storage layout and optimize the struct fields to minimize wasted space and reduce gas costs. In this writeup, we'll explore an example that demonstrates struct packing optimization.

## Unoptimized Struct

Let's start with an unoptimized struct:

```solidity
contract NotOptimizedStruct {
    struct Product {
        uint256 id;         // 32 bytes
        bool isAvailable;   // 1 byte
        uint256 price;      // 32 bytes
        uint32 quantity;    // 4 bytes
        string name;        // dynamic size
        uint32 category;    // 4 bytes
        address seller;     // 20 bytes
        uint16 ratings;     // 2 bytes
    }
}
```

In the `NotOptimizedStruct` contract, the `Product` struct is not optimized for storage. The variables are not arranged in a way that minimizes gaps caused by Solidity's storage layout. Let's analyze the storage layout of the unoptimized struct:

- `uint256 id` occupies 32 bytes (1 slot)
- `bool isAvailable` occupies 1 byte, wasting 31 bytes in the slot
- `uint256 price` occupies 32 bytes (1 slot)
- `uint32 quantity` occupies 4 bytes, wasting 28 bytes in the slot
- `string name` is a dynamic size variable and is stored separately
- `uint32 category` occupies 4 bytes, wasting 28 bytes in the slot
- `address seller` occupies 20 bytes, wasting 12 bytes in the slot
- `uint16 ratings` occupies 2 bytes, wasting 30 bytes in the slot

In total, the unoptimized struct wastes a significant amount of storage space.

## Optimized Struct

Now, let's optimize the struct for better storage efficiency:

```solidity
contract OptimizedStruct {
    struct Product {
        uint256 id;         // 32 bytes
        uint256 price;      // 32 bytes
        address seller;     // 20 bytes
        string name;        // dynamic size
        uint32 quantity;    // 4 bytes
        uint32 category;    // 4 bytes
        uint16 ratings;     // 2 bytes
        bool isAvailable;   // 1 byte
    }
}
```

In the `OptimizedStruct` contract, the `Product` struct is arranged in a way that minimizes wasted space:

- `uint256 id` occupies 32 bytes (1 slot)
- `uint256 price` occupies 32 bytes (1 slot)
- `address seller` occupies 20 bytes, leaving 12 bytes for other variables
- `string name` is a dynamic size variable and is stored separately
- `uint32 quantity` occupies 4 bytes, packed with `category` and `ratings`
- `uint32 category` occupies 4 bytes, packed with `quantity` and `ratings`
- `uint16 ratings` occupies 2 bytes, packed with `quantity` and `category`
- `bool isAvailable` occupies 1 byte, utilizing the remaining byte in the slot

The optimized struct arranges the variables in a way that minimizes wasted space. The `uint256` variables are placed first since they occupy full slots. The `address` variable is placed next, leaving space for smaller variables to be packed together. The `uint32`, `uint16`, and `bool` variables are packed together in the remaining slots, utilizing the available space efficiently.

## Example of storage variable caching

## Example of calldata optimization

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract GasOptimizationExample {
    struct User {
        string name;
        uint256 age;
    }

    User[] private users;

    // Non-optimized function using memory
    function addUserNonOptimized(string memory _name, uint256 _age) external {
        users.push(User(_name, _age));
    }

    // Function to get user details
    function getUser(uint256 _index) external view returns (string memory, uint256) {
        require(_index < users.length, "Invalid user index");
        User memory user = users[_index];
        return (user.name, user.age);
    }

    // Function to get the number of users
    function getUserCount() external view returns (uint256) {
        return users.length;
    }
}
```

```solidity
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
}
```

**Original Function: sumOfArray**

- The original sumOfArray function takes an array of uint256 as a memory parameter. When this function is called, the array passed to the function is copied from calldata (where the input data of a transaction is stored) to memory. This copy operation consumes additional gas, and since memory is a more expensive place to store data than calldata, it increases the cost of executing the function.

**Optimized Function: sumOfArrayOptimized**

- The optimized sumOfArrayOptimized function takes the same array but declares it as calldata instead of memory. This means that when the function is called, the EVM doesn't need to copy the array from calldata to memory; it can read the values directly from calldata.
