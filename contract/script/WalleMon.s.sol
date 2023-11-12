// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {console2} from "forge-std/Test.sol";
import {Script} from "forge-std/Script.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {Referral} from "../src/Referral.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {ERC6551Registry} from "../src/ERC6551Registry.sol";
import {ERC6551AccountProxy} from "../src/ERC6551Upgradeable/ERC6551AccountProxy.sol";
import {ERC6551AccountUpgradeable} from "../src/ERC6551Upgradeable/ERC6551AccountUpgradeable.sol";

contract DeployWalleMon is Script {
    function run() external returns (address) {
        vm.startBroadcast();
        address proxy = deployWalleMon();
        // address referral = deployReferral();
        address referral = address(0x6972D5282c530fE0F92797578582fdeb5aC7414D);// TODO: replace to referral address
        // address referral = address(0x5FbDB2315678afecb367f032d93F642f64180aa3); // local
        // revert("Hey, did you replace the referral address and remove console2.log?");
    
        // ERC6551 contracts
        address registry = deployERC6551Registry();
        address implementation = deployERC6551AccountUpgradeable();
        address payable accountProxy = deployERC6551AccountProxy(implementation);

        // init WalleMon
        initWalletMon(proxy, registry, accountProxy, referral);

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
        WalleMon walleMon = new WalleMon();
        console2.log("deployWalleMon...");
        console2.log("walleMon: ", address(walleMon));
        console2.log("this: ", address(this));
        console2.log("msg.sender: ", msg.sender);
        ERC1967Proxy proxy = new ERC1967Proxy(address(walleMon), "");
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

    function initWalletMon(
        address proxy, 
        address registry, 
        address payable accountProxy, 
        address referral
    ) public {
        WalleMon w = WalleMon(proxy);
        console2.log("initWalletMon msg.sender: ", msg.sender);
        w.initialize(msg.sender, registry, accountProxy, referral);
        Referral r = Referral(referral);
        r.setOwner(msg.sender);
    }

    function deployERC6551AccountUpgradeable() public returns (address) {
        ERC6551AccountUpgradeable account = new ERC6551AccountUpgradeable();
        console2.log("deployAccount...");
        console2.log("accountUpgradeable: ", address(account));
        console2.log("this: ", address(this));
        console2.log("msg.sender: ", msg.sender);
        return address(account);
    }

    function deployERC6551AccountProxy(address implementation) public returns (address payable) {
        ERC6551AccountProxy proxy = new ERC6551AccountProxy(implementation);
        console2.log("deployAccountProxy...");
        console2.log("accountProxy: ", address(proxy));
        console2.log("this: ", address(this));
        console2.log("msg.sender: ", msg.sender);
        return payable(address(proxy));
    }
}