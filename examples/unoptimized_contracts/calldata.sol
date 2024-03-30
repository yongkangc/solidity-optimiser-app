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
}