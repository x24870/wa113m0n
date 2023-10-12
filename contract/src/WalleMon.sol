// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

contract WalleMon is Initializable, ERC721Upgradeable, OwnableUpgradeable, UUPSUpgradeable {
    uint256 private _nextTokenId;
    uint256 internal _num;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address initialOwner) initializer public {
        __ERC721_init("WalleMon", "WLM");
        __Ownable_init(initialOwner);
        __UUPSUpgradeable_init();
        _num = 69;
    }

    function _baseURI() internal pure override returns (string memory) {
        return "wallemon.xyz";
    }

    function safeMint(address to) public onlyOwner {
        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
    }

    function _authorizeUpgrade(address newImplementation)
        internal
        onlyOwner
        override
    {}

    function add (uint256 a, uint256 b) public pure returns (uint256) {
        return a + b;
    }

    function version() public pure returns (uint256) {
        return 1;
    }

    function getNum() public view returns (uint256) {
        return _num;
    }

    function setNum(uint256 num) public {
        _num = num;
    }

    function incrementNum() public {
        _num++;
    }
}