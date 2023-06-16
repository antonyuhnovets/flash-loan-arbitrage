// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0 <=0.8.17;

interface IWETH {
    function deposit() external payable;

    function withdraw(uint256) external;

    function transfer(address to, uint value) external returns (bool);
}
