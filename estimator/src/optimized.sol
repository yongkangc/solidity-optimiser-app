pragma solidity ^0.8.0;

contract Optimized {
  struct UnoptimizedProduct {
    uint256 id;
    uint256 price;
    string name;
    address seller;
    uint32 quantity;
    uint32 category;
    uint16 ratings;
    bool isAvailable;
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
  function sumOfArrayUnoptimized(uint256[] calldata numbers) external pure returns (uint256) {
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

