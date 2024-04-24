pragma solidity ^0.8.13;

import {Test, console} from "forge-std/Test.sol";
import { Unoptimized } from "../src/unoptimized.sol";

contract UnoptimizedTest is Test {
    Unoptimized public myContract;
    function setUp() public {
        myContract = new Unoptimized();
    }
    
// Example test code
function test() public view {
    // use the myContract variable to interact with the contract
    myContract.calculateSumUnoptimized();
}
    
}
