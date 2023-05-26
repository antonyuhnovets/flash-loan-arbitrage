// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "@aave/core-v3/contracts/dependencies/openzeppelin/contracts/IERC20.sol";
import "@aave/core-v3/contracts/flashloan/base/FlashLoanSimpleReceiverBase.sol";
import "@aave/core-v3/contracts/interfaces/IPoolAddressesProvider.sol";


contract ArbitrargeFlashLoan is FlashLoanSimpleReceiverBase {
    address payable owner;

    constructor(address _addressProvider)
        FlashLoanSimpleReceiverBase(IPoolAddressesProvider(_addressProvider))
    {
        owner = payable(msg.sender);
    }

    // Take the loan
    function fn_RequestFlashLoan(address _token, uint256 _amount) public {
        address receiverAddress = address(this);
        address asset = _token;
        uint256 amount = _amount;
        bytes memory params = "";
        uint16 referralCode = 0;

        POOL.flashLoanSimple(
            receiverAddress,
            asset,
            amount,
            params,
            referralCode
        );
    }

    // Arbitrage action
    function arbitrargeUSDC(address _tokenAddress, uint256 _amount) private  returns(bool) {
        uint256 arbitraged_amount = (_amount / 10);       
        IERC20 token = IERC20(_tokenAddress);
        return token.transfer(owner, arbitraged_amount);
    }

    //  
    function executeOperation(
        address asset, 
        uint256 amount, 
        uint256 premium, 
        address initiator, 
        bytes calldata params
        ) external override returns (bool) {
        
        // Perform all logic before repaying the loan
        bool status = arbitrargeUSDC(asset, amount);

        // Repay the loan
        uint256 repayAmount = amount + premium;
        IERC20(asset).approve(address(POOL), totalAmount);

        return status;
    }

    // Recieve the Ether
    receive() external payable {}
}
