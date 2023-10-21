// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721EnumerableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {Referral} from "./Referral.sol";

contract WalleMon is Initializable, ERC721Upgradeable, ERC721EnumerableUpgradeable, ERC721URIStorageUpgradeable, OwnableUpgradeable, UUPSUpgradeable {
    enum Health { HEALTHY, SICK, DEAD }
    struct State {
        Health health;
        uint32 lastMealTime;
        uint32 lastSickTime;
        uint32 lastHealTime;
    }

    // misc params
    bool private _revealed;
    string private _eggURI;
    uint32 private _hungryDuration; // the time interval from last feed time, after which a WalleMon is sick
    uint32 private _sickDuration; // the time interval from last get sick time, after which a WalleMon is dead
    uint32 private _invincibleDuration; // the time interval after a WalleMon is healed, during which it cannot be sick again

    // roles
    Referral private _referral;
    address private _refOwner;
    // token state
    uint256 private _nextTokenId;
    mapping(uint256 => State) private _states;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address initialOwner, address referral) initializer public {
        __ERC721_init("WalleMon", "WLM");
        __ERC721Enumerable_init();
        __ERC721URIStorage_init();
        __Ownable_init(initialOwner);
        __UUPSUpgradeable_init();
        _hungryDuration = 5 seconds;
        _sickDuration = 3 seconds;
        _invincibleDuration = 3 seconds;
        _referral = Referral(referral);
        _refOwner = initialOwner;
        _revealed = false;
        _eggURI = "https://ipfs.blocto.app/ipfs/Qmbnvcgrwjo9amrTGy9kc5VDejYN1ZyrLckXQFoMkDRS9B";
    }

    function _baseURI() internal pure override returns (string memory) {
        return "https://ipfs.blocto.app/ipfs/Qmbnvcgrwjo9amrTGy9kc5VDejYN1ZyrLckXQFoMkDRS9B";
    }

    function setEggURI(string calldata eggURI) public onlyOwner {
        _eggURI = eggURI;
    }

    function setRevealed(bool revealed) public onlyOwner {
        _revealed = revealed;
    }

    function safeMint(address to, string memory uri) public onlyOwner {
        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);
        feed(tokenId);
    }

    function userMint(string calldata refCode, bytes calldata sig) public {
        // here we want to keep the original msg.sender, but also access Referral contract storage
        // so send original msg.sender as first param, then call the Referral contract
        _referral.claim(msg.sender, refCode, sig);
        uint256 tokenId = _nextTokenId++;
        _safeMint(msg.sender, tokenId);
        _setTokenURI(tokenId, "");
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
        _states[tokenId].lastMealTime = uint32(block.timestamp);
    }

    function sick(uint256 tokenId) public onlyOwner() {
        require(
            _states[tokenId].health == Health.HEALTHY,
            "WalleMon: dead or sick"
        );
        _states[tokenId].health = Health.SICK;
    }

    function heal(uint256 tokenId) public onlyOwnerOrTokenOwner(tokenId) {
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

    function batchSick(uint256[] calldata tokenIds) public onlyOwner() {
        for (uint256 i = 0; i < tokenIds.length; i++) {
            sick(tokenIds[i]);
        }
    }

    function batachKill(uint256[] calldata tokenIds) public onlyOwner() {
        for (uint256 i = 0; i < tokenIds.length; i++) {
            kill(tokenIds[i]);
        }
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

    // TODO: maybe set to onlyOwner
    function toBeSickList() public view returns (uint256[] memory) {
        uint256 counter = 0;

        // First pass to determine the size
        for (uint256 i = 0; i < _nextTokenId; i++) {
            if (_states[i].health == Health.HEALTHY && 
                block.timestamp - _states[i].lastMealTime > _hungryDuration &&
                block.timestamp - _states[i].lastHealTime > _invincibleDuration
            ) {
                counter++;
            }
        }

        // Allocate the memory array
        uint256[] memory result = new uint256[](counter);

        // Second pass to populate the results
        uint256 resultIndex = 0;
        for (uint256 i = 0; i < _nextTokenId; i++) {
            if (_states[i].health == Health.HEALTHY && 
                block.timestamp - _states[i].lastMealTime > _hungryDuration &&
                block.timestamp - _states[i].lastHealTime > _invincibleDuration
            ) {
                result[resultIndex] = i;
                resultIndex++;
            }
        }

        return result;
    }


    // TODO: maybe set to onlyOwner
    function toBeDeadList() public view returns (uint256[] memory) {
        uint256 counter = 0;

        // First pass to determine the size
        for (uint256 i = 0; i < _nextTokenId; i++) {
            if (_states[i].health == Health.SICK && 
                block.timestamp - _states[i].lastSickTime > _sickDuration
            ) {
                counter++;
            }
        }

        // Allocate the memory array
        uint256[] memory result = new uint256[](counter);

        // Second pass to populate the results
        uint256 resultIndex = 0;
        for (uint256 i = 0; i < _nextTokenId; i++) {
            if (_states[i].health == Health.SICK && 
                block.timestamp - _states[i].lastSickTime > _sickDuration
            ) {
                result[resultIndex] = i;
                resultIndex++;
            }
        }

        return result;
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
        if (!_revealed) {
            return _eggURI;
        }
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

    function blocktime() public view returns (uint256) {
        return block.timestamp;
    }
}