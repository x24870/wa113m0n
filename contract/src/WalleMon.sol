// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721EnumerableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

contract WalleMon is Initializable, ERC721Upgradeable, ERC721EnumerableUpgradeable, ERC721URIStorageUpgradeable, OwnableUpgradeable, UUPSUpgradeable {
    enum Health { HEALTHY, SICK, DEAD }
    struct State {
        Health health;
        uint256 lastMealTime;
        uint256 lastHealTime;
    }

    uint256 private _nextTokenId;
    uint256 private _num;
    mapping(uint256 => State) private _states;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address initialOwner) initializer public {
        __ERC721_init("WalleMon", "WLM");
        __ERC721Enumerable_init();
        __ERC721URIStorage_init();
        __Ownable_init(initialOwner);
        __UUPSUpgradeable_init();
    }

    function _baseURI() internal pure override returns (string memory) {
        return "https://wallemon.xyz";
    }

    function safeMint(address to, string memory uri) public onlyOwner {
        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);
        feed(tokenId);
    }

    function _authorizeUpgrade(address newImplementation)
        internal
        onlyOwner
        override
    {}

    // WalletMon logic functions
    function feed(uint256 tokenId) public onlyOwnerOrTokenOwner(tokenId) {
        require(
            _states[tokenId].health == Health.HEALTHY,
            "WalleMon: dead or sick"
        );
        _states[tokenId].lastMealTime = block.timestamp;
    }

    function sick(uint256 tokenId) public onlyOwner() {
        require(
            _states[tokenId].health == Health.HEALTHY,
            "WalleMon: dead or sick"
        );
        _states[tokenId].health = Health.SICK;
    }

    function heal(uint256 tokenId) public onlyTokenOwner(tokenId) {
        require(
            _states[tokenId].health == Health.SICK,
            "WalleMon: not sick"
        );
        _states[tokenId].health = Health.HEALTHY;
    }

    function kill(uint256 tokenId) public onlyOwner() {
        require(
            _states[tokenId].health == Health.SICK,
            "WalleMon: not sick"
        );
        _states[tokenId].health = Health.DEAD;
    }

    // View functions
    function health(uint256 tokenId) public view returns (uint8) {
        return uint8(_states[tokenId].health);
    }

    function lastMealTime(uint256 tokenId) public view returns (uint256) {
        return _states[tokenId].lastMealTime;
    }

    function lastHealTime(uint256 tokenId) public view returns (uint256) {
        return _states[tokenId].lastHealTime;
    }


    // Modifiers
    modifier onlyTokenOwner(uint256 tokenId) {
        require(msg.sender == ownerOf(tokenId), "WalleMon: not token owner");
        _;
    }

    modifier onlyOwnerOrTokenOwner(uint256 tokenId) {
        require(
            msg.sender == ownerOf(tokenId) || msg.sender == owner(),
            "WalleMon: not token owner or owner"
        );
        _;
    }

    // The following functions are overrides required by Solidity.
    function _update(address to, uint256 tokenId, address auth)
        internal
        override(ERC721Upgradeable, ERC721EnumerableUpgradeable)
        returns (address)
    {
        return super._update(to, tokenId, auth);
    }

    function _increaseBalance(address account, uint128 value)
        internal
        override(ERC721Upgradeable, ERC721EnumerableUpgradeable)
    {
        super._increaseBalance(account, value);
    }

    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721Upgradeable, ERC721URIStorageUpgradeable)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721Upgradeable, ERC721EnumerableUpgradeable, ERC721URIStorageUpgradeable)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

        function add (uint256 a, uint256 b) public pure returns (uint256) {
        return a + b;
    }

    // playground functions
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

    function blocktime() public view returns (uint256) {
        return block.timestamp;
    }
}