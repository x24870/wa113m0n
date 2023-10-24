// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {console2} from "forge-std/Test.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

contract Referral {
    address private owner;
    // claimed address => referral code
    mapping(address => string) private claimed;
    // referral code => claimed count
    // TODO: use byte32 may save gas
    mapping(string => uint32) private referralCount;
    // refferal code => allowed amount
    mapping(string => uint32) private referralAmount;

    constructor() {
        owner = msg.sender;
    }

    // getter and setter
    function getOwner() public view returns (address) {
        return owner;
    }

    function setOwner(address _owner) public {
        require(msg.sender == owner, "Referral: not owner");
        owner = _owner;
    }

    function getClaimed(address _claimed) public view returns (string memory) {
        return claimed[_claimed];
    }

    function getReferralCount(string calldata _referralCode) public view returns (uint32) {
        return referralCount[_referralCode];
    }

    function getReferralAmounts(string calldata _referralCode) public view returns (uint32) {
        return referralAmount[_referralCode];
    }

    function setReferralAmounts(string calldata _referralCode, uint32 _amount) public {
        require(msg.sender == owner, "Referral: not owner");
        referralAmount[_referralCode] = _amount;
    }

    // 
    function claim(address minter, string calldata _referralCode, bytes calldata signature) public {
        require(Strings.equal(claimed[minter], ""), "Referral: already claimed");
        require(verify(owner, minter, _referralCode, signature), 
            "Referral: invalid signature");
        require(referralAmount[_referralCode] > 0, "Referral: no amount left");
        claimed[msg.sender] = _referralCode;
        referralCount[_referralCode] += 1;
        referralAmount[_referralCode] -= 1;
    }

    // Owner signed message format: keecak256(claimedAddress,referralCode)
    function getMessageHash(
        address _claimTo,
        string memory _referralCode
    ) public pure returns (bytes32) {
        return keccak256(abi.encodePacked(_claimTo, _referralCode));
    }

    function getEthSignedMessageHash(
        bytes32 _messageHash
    ) public pure returns (bytes32) {
        /*
        Signature is produced by signing a keccak256 hash with the following format:
        "\x19Ethereum Signed Message\n" + len(msg) + msg
        */
        return
            keccak256(
                abi.encodePacked("\x19Ethereum Signed Message:\n32", _messageHash)
            );
    }

    function verify(
        address _signer,
        address _claimTo,
        string memory _referralCode,
        bytes memory signature
    ) public pure returns (bool) {
        bytes32 messageHash = getMessageHash(_claimTo, _referralCode);
        bytes32 ethSignedMessageHash = getEthSignedMessageHash(messageHash);

        return recoverSigner(ethSignedMessageHash, signature) == _signer;
    }

    function recoverSigner(
        bytes32 _ethSignedMessageHash,
        bytes memory _signature
    ) public pure returns (address) {
        (bytes32 r, bytes32 s, uint8 v) = splitSignature(_signature);

        return ecrecover(_ethSignedMessageHash, v, r, s);
    }

    function splitSignature(
        bytes memory sig
    ) public pure returns (bytes32 r, bytes32 s, uint8 v) {
        require(sig.length == 65, "invalid signature length");

        assembly {
            /*
            First 32 bytes stores the length of the signature

            add(sig, 32) = pointer of sig + 32
            effectively, skips first 32 bytes of signature

            mload(p) loads next 32 bytes starting at the memory address p into memory
            */

            // first 32 bytes, after the length prefix
            r := mload(add(sig, 32))
            // second 32 bytes
            s := mload(add(sig, 64))
            // final byte (first byte of the next 32 bytes)
            v := byte(0, mload(add(sig, 96)))
        }

        // implicitly return (r, s, v)
    }
}