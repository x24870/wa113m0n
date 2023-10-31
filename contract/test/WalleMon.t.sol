// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {Preheat} from "../src/Preheat.sol";
import {Referral} from "../src/Referral.sol";
import {ERC6551Registry} from "../src/ERC6551Registry.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {DeployWalleMon} from "../script/WalleMon.s.sol";
import {DeployPreheat} from "../script/PreHeat.s.sol";

contract WalleMonTest is Test {
    DeployWalleMon public deployWalleMon;
    DeployPreheat public deployPreheat;
    Preheat public preheat;
    ERC6551Registry public registry;
    Referral public referral;
    WalleMon public walleMon;
    address public refAddr;
    address public preheatAddr;
    address public registryAddr;
    address public proxy;
    address public owner;
    address public mintTo;

    function setUp() public {
        owner = msg.sender;
        vm.startPrank(owner);
        console2.log("owner: ", owner);
        console2.log("msg.sender: ", msg.sender);

        // deploy contracts
        deployPreheat = new DeployPreheat();
        deployWalleMon = new DeployWalleMon();
        
        // preheatAddr = deployPreheat.deployPreheat();
        // ref = deployPreheat.deployReferral();
        (bool success, bytes memory result) = address(deployPreheat).delegatecall(abi.encodeWithSignature("deployPreheat()"));
        assertEq(success, true);
        preheatAddr = abi.decode(result, (address));
        (success, result) = address(deployPreheat).delegatecall(abi.encodeWithSignature("deployReferral()"));
        assertEq(success, true);
        refAddr = abi.decode(result, (address));
        
        registryAddr = deployWalleMon.deployERC6551Registry();
        proxy = deployWalleMon.deployWalleMon();
        // deployWalleMon.initWalletMon(proxy, registry, ref);
        (success, result) = address(deployWalleMon).delegatecall(abi.encodeWithSignature("initWalletMon(address,address,address)", proxy, registryAddr, refAddr));
        assertEq(success, true);

        // setup
        walleMon = WalleMon(proxy);
        referral = Referral(refAddr);
        registry = ERC6551Registry(registryAddr);
        preheat = Preheat(preheatAddr);
        mintTo = address(0x1);

        console2.log("!!!preheat owner: ", preheat.owner());

        vm.stopPrank();
    }

    // function testPreheat() public {
    //     vm.startPrank(owner);

    //     // test balanceOf
    //     assertEq(0, preheat.balanceOf(mintTo));

    //     // test burn
    //     console2.log("************* preheatAddress: ", preheatAddr);
    //     console2.log("************* msg.sender: ", msg.sender);
    //     console2.log("************* owner: ", preheat.owner());
    //     preheat.safeMint(mintTo);
    //     assertEq(1, preheat.balanceOf(mintTo));
    //     preheat.burn(0);
    //     assertEq(0, preheat.balanceOf(mintTo));

        
    //     vm.stopPrank();
    // }

    function testWalleMon() public {
        vm.startPrank(owner);

        // test balanceOf
        assertEq(0, WalleMon(proxy).balanceOf(mintTo));

        // test owner
        console2.log("************* proxyAddress: ", proxy);
        console2.log("msg.sender: ", msg.sender);
        console2.log("walleMon owner", walleMon.owner());
        assertEq(msg.sender, walleMon.owner());

        // test set ref code
        string memory refCode = "wallemon";
        console2.log("referral owner: ", referral.getOwner());
        referral.setReferralAmounts(refCode, 1);
        assertEq(referral.getReferralAmounts(refCode), 1);

        vm.stopPrank();
    }

    function testWalleMonGame() public {
        vm.prank(owner);
        walleMon.setRevealed(true);
        vm.prank(owner);
        walleMon.safeMint(mintTo, "tokenURI");
        assertEq(1, walleMon.balanceOf(mintTo));
        assertEq(mintTo, walleMon.ownerOf(0));
        
        uint256 bornMealTime = block.timestamp;
        assertEq(block.timestamp, bornMealTime);
        vm.warp(block.timestamp+ 1 minutes);
        vm.roll(block.number+1);

        // feed
        vm.prank(mintTo);
        walleMon.feed(0);
        assertEq(block.timestamp, walleMon.lastMealTime(0));
        assertEq(bornMealTime + 1 minutes, walleMon.lastMealTime(0));
     
        // sick
        vm.prank(owner);
        walleMon.sick(0);
        assertEq(uint8(WalleMon.Health.SICK), walleMon.health(0));

        // heal
        vm.prank(mintTo);
        walleMon.heal(0);
        assertEq(uint8(WalleMon.Health.HEALTHY), walleMon.health(0));

        // kill
        vm.prank(owner);
        walleMon.sick(0);
        vm.prank(owner);
        walleMon.kill(0);
        assertEq(uint8(WalleMon.Health.DEAD), walleMon.health(0));
    }

    function testWalleMonOwnership() public {
        vm.startPrank(owner);
        console2.log("wallemon owner: ", walleMon.owner());
        assertEq(owner, walleMon.owner());
        address newOwner = address(0x1);
        walleMon.transferOwnership(newOwner);
        assertEq(newOwner, walleMon.owner());
        vm.stopPrank();

        vm.startPrank(newOwner);
        walleMon.safeMint(address(0x2), "");
        assertEq(1, walleMon.balanceOf(address(0x2)));
        // return ownership
        walleMon.transferOwnership(owner);
        assertEq(owner, walleMon.owner());
        vm.stopPrank();
    }
}
