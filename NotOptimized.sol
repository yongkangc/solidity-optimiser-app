// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

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