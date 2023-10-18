// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {Test, console2} from "forge-std/Test.sol";
import {Referral} from "../src/Referral.sol";

contract RefferalTest is Test {
    Referral public referral;
    address internal _owner;
    uint256 internal _privateKey;
    address internal _user;
    uint256 internal _userPrivateKey;

    function setUp() public {
        _owner = address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266);
        _privateKey = 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80;
        _user = address(0x70997970C51812dc3A010C7d01b50e0d17dc79C8);
        _userPrivateKey = 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d;

        vm.prank(_owner);
        referral = new Referral();
        assertEq(referral.getOwner(), _owner);
    }

    function testClaim() public {
        string memory refCode = "CryptoKOL";
        vm.prank(_owner);
        referral.setReferralAmounts(refCode, 1);
        assertEq(referral.getReferralAmounts(refCode), 1);

        // user get claim message from BE
        bytes memory sig = getOwnerSignedMsg(_user, refCode);
        vm.prank(_user);
        referral.claim(refCode, sig);

        assertEq(referral.getClaimed(_user), refCode);
        assertEq(referral.getReferralCount(refCode), 1);
        assertEq(referral.getReferralAmounts(refCode), 0);

        // user claim again
        vm.expectRevert("Referral: already claimed");
        vm.prank(_user);
        referral.claim(refCode, sig);

        // referral code allowed amount is 0
        vm.expectRevert("Referral: no amount left");
        address user2 = address(0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC);
        sig = getOwnerSignedMsg(user2, refCode);
        vm.prank(user2);
        referral.claim(refCode, sig);

        // fake signer
        vm.expectRevert("Referral: invalid signature");
        uint256 user2PrivateKey = 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a;
        bytes32 msgHash = keccak256(abi.encodePacked(user2, refCode));
        bytes32 signedMsgHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", msgHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(user2PrivateKey, signedMsgHash);
        sig = abi.encodePacked(r, s, v);
        vm.prank(user2);
        referral.claim(refCode, sig);
    }

    function testVerify() public {
        string memory refCode = "CryptoKOL";
        bytes memory sig = getOwnerSignedMsg(_user, refCode);
        assertEq(true, referral.verify(_owner, _user, refCode, sig));
    }

    function getOwnerSignedMsg(address _claimTo, string memory _refCode) internal view returns (bytes memory) {
        bytes32 msgHash = keccak256(abi.encodePacked(_claimTo, _refCode));
        bytes32 signedMsgHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", msgHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(_privateKey, signedMsgHash);
        bytes memory signature = abi.encodePacked(r, s, v);
        return signature;
    }
}