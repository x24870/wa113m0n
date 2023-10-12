// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {console2} from "forge-std/Test.sol";
import {Script} from "forge-std/Script.sol";
import {WalleMon} from "../src/WalleMon.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract UpgradeWalleMon is Script {
    function run() external returns (address) {
        vm.startBroadcast();
        WalleMon newWalleMon = new WalleMon();
        address proxy = upgradeWalleMon(
            address(0x998abeb3E57409262aE5b751f60747921B33613E), // legacy implementation address
            address(newWalleMon)
            );
        vm.stopBroadcast();
        return proxy;
    }

    function upgradeWalleMon(
        address proxyAddr,
        address newImplementation
    ) public returns (address) {
        WalleMon proxy = WalleMon(payable(proxyAddr));
        proxy.upgradeToAndCall(address(newImplementation), "");
        return address(proxy);
    }

}