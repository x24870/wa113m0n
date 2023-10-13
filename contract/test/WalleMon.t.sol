// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {DeployWalleMon} from "../script/WalleMon.s.sol";

contract WalleMonTest is Test {
    DeployWalleMon public deployWalleMon;
    // address public onwer = address(this);

    function setUp() public {
        deployWalleMon = new DeployWalleMon();
    }

    // function testWalleMonWorks() public {
    //     address proxyAddress = deployWalleMon.deployWalleMon();
    //     deployWalleMon.initWalletMon(proxyAddress);
    //     uint256 expectedValue = 1;
    //     assertEq(expectedValue, WalleMon(proxyAddress).version());

    //     address to = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
    //     assertEq(0, WalleMon(proxyAddress).balanceOf(to));

    //     WalleMon w = WalleMon(proxyAddress);
    //     uint256 b = w.balanceOf(to);
    //     assertEq(0, b);

    //     console2.log("************* proxyAddress: ", proxyAddress);
    //     console2.log("msg.sender: ", msg.sender);
    //     console2.log("walleMon", w.owner());

    //     w.safeMint(to, "tokenURI");
    //     assertEq(1, w.balanceOf(to));

    //     assertEq(69, w.getNum());
    // }

}
