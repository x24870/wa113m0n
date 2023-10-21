// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {Referral} from "../src/Referral.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {DeployWalleMon} from "../script/WalleMon.s.sol";

contract WalleMonTest is Test {
    DeployWalleMon public deployWalleMon;
    Referral public referral;
    WalleMon public walletMon;
    address public ref;
    address public proxy;
    address public owner;
    address public mintTo;

    function setUp() public {
        owner = msg.sender;
        vm.startPrank(owner);

        // deploy WalleMon
        deployWalleMon = new DeployWalleMon();
        proxy = deployWalleMon.deployWalleMon();
        ref = deployWalleMon.deployReferral();
        deployWalleMon.initWalletMon(proxy, ref);

        // setup
        walletMon = WalleMon(proxy);
        referral = Referral(ref);
        mintTo = address(0x1);

        vm.stopPrank();
    }

    function testWalleMon() public {
        vm.startPrank(owner);

        // test balanceOf
        assertEq(0, WalleMon(proxy).balanceOf(mintTo));

        // test owner
        console2.log("************* proxyAddress: ", proxy);
        console2.log("msg.sender: ", msg.sender);
        console2.log("walleMon owner", walletMon.owner());
        assertEq(msg.sender, walletMon.owner());

        // test set ref code
        string memory refCode = "wallemon";
        console2.log("referral owner: ", referral.getOwner());
        referral.setReferralAmounts(refCode, 1);
        assertEq(referral.getReferralAmounts(refCode), 1);

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

    function getOwnerSignedMsg(address _claimTo, string memory _refCode, uint256 ownerPk) internal view returns (bytes memory) {
        bytes32 msgHash = keccak256(abi.encodePacked(_claimTo, _refCode));
        bytes32 signedMsgHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", msgHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPk, signedMsgHash);
        bytes memory signature = abi.encodePacked(r, s, v);
        return signature;
    }
}
