// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {console2} from "forge-std/Test.sol";
import {Script} from "forge-std/Script.sol";
import {Referral} from "../src/Referral.sol";

contract DeployWalleMon is Script {
    Referral public referral;

    function run() external {
        vm.startBroadcast();

        // set single referral amount
        referral = Referral(address(0x6972D5282c530fE0F92797578582fdeb5aC7414D));
        // referral.setReferralAmounts("wallemon", 3);

        // set multiple referral amounts
        string[] memory referralCodes = new string[](5);
        uint32[] memory amounts = new uint32[](5);
        referralCodes[0] = "qa";
        // referralCodes[1] = "wallemon2";
        // referralCodes[2] = "wallemon3";
        // referralCodes[3] = "wallemon4";
        // referralCodes[4] = "wallemon5";
        amounts[0] = 50;
        // amounts[1] = 2;
        // amounts[2] = 3;
        // amounts[3] = 4;
        // amounts[4] = 5;
        referral.batchSetReferralAmounts(referralCodes, amounts);

        vm.stopBroadcast();
    }


}