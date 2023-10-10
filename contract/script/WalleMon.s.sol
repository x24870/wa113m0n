// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Script} from "forge-std/Script.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployWalleMon is Script {
    function run() external returns (address) {
        address proxy = deployWalleMon();
        return proxy;
    }

    function deployWalleMon() public returns (address) {
        vm.startBroadcast();
        WalleMon walleMon = new WalleMon();
        ERC1967Proxy proxy = new ERC1967Proxy(address(walleMon), "");
        vm.stopBroadcast();
        return address(proxy);
    }
}