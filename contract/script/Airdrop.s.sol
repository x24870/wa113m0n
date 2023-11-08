// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Script} from "forge-std/Script.sol";
import {WalleMon} from "../src/WalleMon.sol";

contract Airdrop is Script {
    WalleMon public wallemon;

    function run() external {
        vm.startBroadcast();

        wallemon = WalleMon(address(0x9fc184Dc6A94B43e56478e817e31C15315Ac0757));
        address to = address(0x0);
        wallemon.safeMint(to, "");

        vm.stopBroadcast();
    }
}