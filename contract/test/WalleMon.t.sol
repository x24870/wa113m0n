// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {DeployWalleMon} from "../script/WalleMon.s.sol";

contract WalleMonTest is Test {
    DeployWalleMon public deployWalleMon;
    address public OWNER = address(1);

    function setUp() public {
        deployWalleMon = new DeployWalleMon();
    }

    function testWalleMonWorks() public {
        address proxyAddress = deployWalleMon.deployWalleMon();
        uint256 expectedValue = 1;
        assertEq(expectedValue, WalleMon(proxyAddress).version());
    }

    // function testAdd() public {
    //     console2.log("testAdd-----------------------------");
    //     assertEq(walleMon.add(1, 2), 3);
    // }

}
