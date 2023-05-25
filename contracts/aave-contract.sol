// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "@aave/core-v3/contracts/dependencies/openzeppelin/contracts/IERC20.sol";
import "@aave/core-v3/contracts/flashloan/base/FlashLoanSimpleReceiverBase.sol";


contract ArbitrargeFlashLoan is FlashLoanSimpleReceiverBase {
    address payable owner;

    constructor(address _addressProvider)
        FlashLoanSimpleReceiverBase(IPoolAddressesProvider(_addressProvider))
    {
        owner = payable(msg.sender);
    }
}