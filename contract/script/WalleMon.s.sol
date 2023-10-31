// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {console2} from "forge-std/Test.sol";
import {Script} from "forge-std/Script.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {ERC6551Registry} from "../src/ERC6551Registry.sol";
import {Referral} from "../src/Referral.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployWalleMon is Script {
    function run() external returns (address) {
        vm.startBroadcast();
        address proxy = deployWalleMon();
        // address referral = deployReferral();
        address referral = address(0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512);
        address registry = deployERC6551Registry();
        initWalletMon(proxy, registry, referral);
        vm.stopBroadcast();
        return proxy;
    }

    // The deployment of Preheat contract also deployed Referral contract
    // So we don't need to deploy Referral contract again
    // function deployReferral() public returns (address) {
    //     Referral referral = new Referral();
    //     return address(referral);
    // }

    function deployWalleMon() public returns (address) {
        // vm.startBroadcast();
        WalleMon walleMon = new WalleMon();
        console2.log("deployWalleMon...");
        console2.log("walleMon: ", address(walleMon));
        console2.log("this: ", address(this));
        console2.log("msg.sender: ", msg.sender);
        ERC1967Proxy proxy = new ERC1967Proxy(address(walleMon), "");
        // vm.stopBroadcast();
        return address(proxy);
    }

    function deployERC6551Registry() public returns (address) {
        ERC6551Registry registry = new ERC6551Registry();
        console2.log("deployRegistry...");
        console2.log("registry: ", address(registry));
        console2.log("this: ", address(this));
        console2.log("msg.sender: ", msg.sender);
        return address(registry);
    }

    function initWalletMon(address proxy, address registry, address referral) public {
        // vm.startBroadcast();
        WalleMon w = WalleMon(proxy);
        console2.log("initWalletMon msg.sender: ", msg.sender);
        w.initialize(msg.sender, registry, referral);
        Referral r = Referral(referral);
        r.setOwner(msg.sender);
        // vm.stopBroadcast();
    }
}