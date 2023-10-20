// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {console2} from "forge-std/Test.sol";
import {Script} from "forge-std/Script.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {Referral} from "../src/Referral.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployWalleMon is Script {
    function run() external returns (address) {
        vm.startBroadcast();
        address proxy = deployWalleMon();
        address referral = deployReferral();
        initWalletMon(proxy, referral);
        vm.stopBroadcast();
        return proxy;
    }

    function deployReferral() public returns (address) {
        Referral referral = new Referral();
        return address(referral);
    }

    function deployWalleMon() public returns (address) {
        // vm.startBroadcast();
        WalleMon walleMon = new WalleMon();
        console2.log("walleMon: ", address(walleMon));
        console2.log("this: ", address(this));
        console2.log("msg.sender: ", msg.sender);
        ERC1967Proxy proxy = new ERC1967Proxy(address(walleMon), "");
        // vm.stopBroadcast();
        return address(proxy);
    }

    function initWalletMon(address proxy, address referral) public {
        // vm.startBroadcast();
        WalleMon w = WalleMon(proxy);
        console2.log("initWalletMon msg.sender: ", msg.sender);
        w.initialize(msg.sender, referral);
        Referral r = Referral(referral);
        r.setOwner(msg.sender);
        // vm.stopBroadcast();
    }
}