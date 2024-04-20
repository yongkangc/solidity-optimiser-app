pragma solidity ^0.8.20;

contract Optimized {
  uint256 public number;
  int256[] public arr;
  function increment() public view returns (uint256) {
    uint256 cached_number = number;
    require(cached_number < 10);
    uint256 incremented = cached_number + 1;
    return incremented;
  }
  function sum() public view returns (int256) {
    int256[] memory cached_arr = arr;
    int256 sum = 0;
    for (uint256 i = 0; i < cached_arr.length; i++) {
      sum += cached_arr[i];
    }
    return sum;
  }
}

