// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract GasOptimizationExample {
    struct User {
        string name;
        uint256 age;
    }

    User[] private users;

    // // Gas-optimized function using calldata
    // function addUserOptimized(string calldata _name, uint256 _age) external {
    //     users.push(User(_name, _age));
    // }

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