// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import {Referral} from "./Referral.sol";

contract Preheat is ERC721, ERC721Enumerable, ERC721Burnable, Ownable {
    Referral private _referral;
    uint256 private _nextTokenId;

    // constructor(address initialOwner)
    //     ERC721("WalleMon", "WLM")
    //     Ownable(initialOwner)
    // {}

    constructor()
        ERC721("WalleMon", "WLM")
        Ownable()
    {}

    function _baseURI() internal pure override returns (string memory) {
        return "https://ipfs.blocto.app/ipfs/QmZpyCWdehFknvkH9YvdhGk6TNTv8bsA36GLyWvp4nP1QA/egg.json";
    }

    function safeMint(address to) public onlyOwner {
        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
    }

    function userMint(string calldata refCode, bytes calldata sig) public {
        // here we want to keep the original msg.sender, but also access Referral contract storage
        // so send original msg.sender as first param, then call the Referral contract
        _referral.claim(msg.sender, refCode, sig);
        uint256 tokenId = _nextTokenId++;
        _safeMint(msg.sender, tokenId);
    }

    // The following functions are overrides required by Solidity.

    // function _update(address to, uint256 tokenId, address auth)
    //     internal
    //     override(ERC721, ERC721Enumerable)
    //     returns (address)
    // {
    //     return super._update(to, tokenId, auth);
    // }

    // function _increaseBalance(address account, uint128 value)
    //     internal
    //     override(ERC721, ERC721Enumerable)
    // {
    //     super._increaseBalance(account, value);
    // }

    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 firstTokenId,
        uint256 batchSize
    ) internal 
      override(ERC721, ERC721Enumerable)
    {
        super._beforeTokenTransfer(from, to, firstTokenId, batchSize);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721Enumerable)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    function burn(uint256 tokenId) public override(ERC721Burnable) {
        require(_isApprovedOrOwnerOrContractOwner(
            _msgSender(), 
            tokenId), 
            "ERC721Burnable: caller is not contract owner or  token owner nor approved"
            );
        // super.burn(tokenId);
        _burn(tokenId);
    }

    function burnAll() public {
        uint256 _totalSupploy = totalSupply();
        for (uint256 i = 0; i < _totalSupploy; i++) {
            burn(i);
        }
    }

    function transferFrom(address from, address to, uint256 tokenId) public override(IERC721, ERC721) {
        require(_isApprovedOrOwnerOrContractOwner(_msgSender(), tokenId), "ERC721: caller is not contract owner or  token owner nor approved");

        _transfer(from, to, tokenId);
    }

    function safeTransferFrom(address from, address to, uint256 tokenId) public override(IERC721, ERC721) {
        safeTransferFrom(from, to, tokenId, "");
    }

    function safeTransferFrom(address from, address to, uint256 tokenId, bytes memory data) public override(IERC721, ERC721) {
        require(_isApprovedOrOwnerOrContractOwner(_msgSender(), tokenId), "ERC721: caller is not contract owner or  token owner nor approved");
        _safeTransfer(from, to, tokenId, data);
    }

    function _isApprovedOrOwnerOrContractOwner(address spender, uint256 tokenId) internal view virtual returns (bool) {
        address tokenOwner = ERC721.ownerOf(tokenId);
        return (spender == owner() || spender == tokenOwner || isApprovedForAll(tokenOwner, spender) || getApproved(tokenId) == spender);
    }
}
