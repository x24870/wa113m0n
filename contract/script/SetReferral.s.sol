// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {console2} from "forge-std/Test.sol";
import {Script} from "forge-std/Script.sol";
import {Referral} from "../src/Referral.sol";

contract DeployWalleMon is Script {
    Referral public referral;

    function run() external {
        vm.startBroadcast();
        referral = Referral(address(0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0));
        referral.setReferralAmounts("wallemon", 3);
        vm.stopBroadcast();
    }


}