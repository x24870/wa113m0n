// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {console2} from "forge-std/Test.sol";
import {Script} from "forge-std/Script.sol";
import {Preheat} from "../src/Preheat.sol";
import {Referral} from "../src/Referral.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployPreheat is Script {
    function run() external returns (address, address) {
        vm.startBroadcast();
        address preheat = deployPreheat();
        address referral = deployReferral();
        vm.stopBroadcast();
        return (preheat, referral);
    }

    function deployReferral() public returns (address) {
        Referral referral = new Referral();
        console2.log("deployReferral...");
        console2.log("refertal: ", address(referral));
        console2.log("this: ", address(this));
        console2.log("msg.sender: ", msg.sender);
        console2.log("owner: ", referral.getOwner());
        return address(referral);
    }

    function deployPreheat() public returns (address) {
        Preheat preheat = new Preheat();
        console2.log("deployPreheat...");
        console2.log("preheat: ", address(preheat));
        console2.log("this: ", address(this));
        console2.log("msg.sender: ", msg.sender);
        return address(preheat);
    }
}