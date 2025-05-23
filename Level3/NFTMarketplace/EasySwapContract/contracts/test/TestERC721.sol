// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract TestERC721 is ERC721Enumerable, Ownable {
    string private _tokenURI;

    constructor() ERC721("Troll", "Troll") Ownable(msg.sender) {}

    function mint(address to, uint256 tokenId) external onlyOwner {
        _mint(to, tokenId);
    }

    function tokenURI(
        uint256 /* tokenId */
    ) public view override returns (string memory) {
        return _tokenURI;
    }

    function setTokenURI(string memory newTokenURI) external onlyOwner {
        _tokenURI = newTokenURI;
    }
}
