pragma solidity ^0.8.13;

import {Test, console} from "forge-std/Test.sol";
import { Optimized } from "../src/optimized.sol";

contract OptimizedTest is Test {
    Optimized public myContract;
    function setUp() public {
        myContract = new Optimized();
    }
    
// Example test code
function test() public view {
    // use the myContract variable to interact with the contract
    myContract.calculateSumUnoptimized();
}
    
}
