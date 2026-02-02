// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract TestToken is ERC20 {
    uint256 public constant RATE = 100000000; // 1 ETH = 100M tokens

    constructor() ERC20("Test Token", "TST") {}

    function mint() public payable {
        require(msg.value >= 0.001 ether, "Min 0.001 ETH");
        uint256 tokensToMint = (msg.value * RATE);
        _mint(msg.sender, tokensToMint);
    }

    // 可以发送 0.001 ETH 到合约获取测试代币
}