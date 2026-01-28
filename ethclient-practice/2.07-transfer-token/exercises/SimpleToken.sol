// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract SimpleToken is ERC20 {
    // 构造函数：铸造初始代币
    constructor() ERC20("Simple Token", "SIMP") {
        // 铸造 1,000,000 代币给部署者
        _mint(msg.sender, 1000000 * 10**18);
    }

    // 公开的 mint 函数：任何人都可以铸造代币
    function mint(uint256 amount) public {
        _mint(msg.sender, amount);
    }

    // 公开的 mintTo 函数：给指定地址铸造代币
    function mintTo(address to, uint256 amount) public {
        _mint(to, amount);
    }
}
