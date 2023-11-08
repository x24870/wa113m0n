// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Script} from "forge-std/Script.sol";
import {Preheat} from "../src/Preheat.sol";

contract Airdrop is Script {
    Preheat public preheat;

    function run() external {
        vm.startBroadcast();

        preheat = Preheat(address(0x149196B0C40a0A12d2b201BEd925beD2813Db744));
        preheat.burnAll();

        vm.stopBroadcast();
    }
}