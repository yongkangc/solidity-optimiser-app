// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract NotOptimizedStruct {
    struct Employee {
        uint256 id;        // 32 bytes
        uint32 salary;     // 4 bytes
        uint32 age;        // 4 bytes
        bool isActive;     // 1 byte
        address addr;      // 20 bytes
        uint16 department; // 2 bytes
    }
}
