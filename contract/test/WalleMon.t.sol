// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {DeployWalleMon} from "../script/WalleMon.s.sol";

contract WalleMonTest is Test {
    DeployWalleMon public deployWalleMon;
    WalleMon public walletMon;
    address public proxy;
    address public owner;
    address public mintTo;

    function setUp() public {
        owner = msg.sender;
        vm.startPrank(owner);

        // deploy WalleMon
        deployWalleMon = new DeployWalleMon();
        proxy = deployWalleMon.deployWalleMon();
        deployWalleMon.initWalletMon(proxy);

        // setup
        walletMon = WalleMon(proxy);
        mintTo = address(0x1);

        vm.stopPrank();
    }

    function testWalleMon() public {
        vm.startPrank(owner);

        // test version
        uint256 expectedValue = 1;
        assertEq(expectedValue, WalleMon(proxy).version());

        // test balanceOf
        assertEq(0, WalleMon(proxy).balanceOf(mintTo));

        // test owner
        console2.log("************* proxyAddress: ", proxy);
        console2.log("msg.sender: ", msg.sender);
        console2.log("walleMon", walletMon.owner());
        assertEq(msg.sender, walletMon.owner());

        vm.stopPrank();
    }

    function testWalleMonGame() public {
        vm.prank(owner);
        walletMon.safeMint(mintTo, "tokenURI");
        assertEq(1, walletMon.balanceOf(mintTo));
        assertEq(mintTo, walletMon.ownerOf(0));
        uint256 bornMealTime = block.timestamp;
        assertEq(block.timestamp, bornMealTime);
        vm.warp(block.timestamp+ 1 minutes);
        vm.roll(block.number+1);

        // feed
        vm.prank(mintTo);
        walletMon.feed(0);
        assertEq(block.timestamp, walletMon.lastMealTime(0));
        assertEq(bornMealTime + 1 minutes, walletMon.lastMealTime(0));
     
        // sick
        vm.prank(owner);
        walletMon.sick(0);
        assertEq(uint8(WalleMon.Health.SICK), walletMon.health(0));

        // heal
        vm.prank(mintTo);
        walletMon.heal(0);
        assertEq(uint8(WalleMon.Health.HEALTHY), walletMon.health(0));

        // kill
        vm.prank(owner);
        walletMon.sick(0);
        vm.prank(owner);
        walletMon.kill(0);
        assertEq(uint8(WalleMon.Health.DEAD), walletMon.health(0));

    }
}
