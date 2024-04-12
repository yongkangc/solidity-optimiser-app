# All Examples

## Example of struct packing

input

```solidity
    struct RequestMeta {
        uint64 completedRequests;
        Custom.Datatype data;
        address requestingContract;
        uint72 adminFee; // in wei
        address subscriptionOwner;
        bytes32 flags; // 32 bytes of flags
        uint96 availableBalance; // in wei. 0 if not specified.
        uint64 subscriptionId;
        uint64 initiatedRequests;// number of requests initiated by this contract
        uint32 callbackGasLimit;
        uint16 dataVersion;
    }
```

Expected output

```solidity
    struct RequestMeta {
        Custom.Datatype data; //
        bytes32 flags; //                  32 bytes of flags
        address requestingContract; // ──╮
        uint96 availableBalance; // ─────╯ in wei. 0 if not specified.
        address subscriptionOwner; // ───╮
        uint64 completedRequests; //     │
        uint32 callbackGasLimit; // ─────╯
        uint72 adminFee; // ─────────────╮ in wei
        uint64 subscriptionId; //        │
        uint64 initiatedRequests; //     │ number of requests initiated by this contract
        uint16 dataVersion; // ──────────╯
    }
```

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
