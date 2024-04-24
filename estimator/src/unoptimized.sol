pragma solidity ^0.8.0;

contract Unoptimized {
  struct UnoptimizedProduct {
    uint256 id;
    bool isAvailable;
    uint256 price;
    uint32 quantity;
    string name;
    uint32 category;
    address seller;
    uint16 ratings;
  }
  uint256 public variable1;
  uint256 public variable2;
  function calculateSumUnoptimized() public view returns (uint256) {
    uint256 sum = variable1 + variable2;
    return sum;
  }
  function calculateSumOptimized() public view returns (uint256) {
    uint256 v1 = variable1;
    uint256 v2 = variable2;
    uint256 sum = v1 + v2;
    return sum;
  }
  function sumOfArrayUnoptimized(uint256[] memory numbers) external pure returns (uint256) {
    uint256 sum = 0;
    for (uint256 i = 0; i < numbers.length; ++i) {
      sum += numbers[i];
    }
    return sum;
  }
  function sumOfArrayOptimized(uint256[] calldata numbers) external pure returns (uint256) {
    uint256 sum = 0;
    for (uint256 i = 0; i < numbers.length; ++i) {
      sum += numbers[i];
    }
    return sum;
  }
}

