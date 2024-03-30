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

contract GasOptimizationExample {
    struct User {
        string name;
        uint256 age;
    }

    User[] private users;

    // Gas-optimized function using calldata
    function addUserOptimized(string calldata _name, uint256 _age) external {
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

**addUserOptimized:**

- This function uses calldata for the \_name parameter, which is an external function input.
- By using calldata, the function reads the data directly from the call data, avoiding the need to copy it to memory.
- This results in gas savings compared to using memory.

**addUserNonOptimized:**

- This function uses memory for the \_name parameter.
- When using memory, the data is first copied to memory before being used, which consumes additional gas.
