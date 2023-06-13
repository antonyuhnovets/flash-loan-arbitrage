// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0 <=0.8.17;
pragma abicoder v2;

import "@aave/core-v3/contracts/flashloan/base/FlashLoanSimpleReceiverBase.sol";
import "@aave/core-v3/contracts/interfaces/IPoolAddressesProvider.sol";
import "@aave/core-v3/contracts/dependencies/openzeppelin/contracts/IERC20.sol";

import "@uniswap/v3-core/contracts/libraries/LowGasSafeMath.sol";
import '@uniswap/v3-periphery/contracts/interfaces/ISwapRouter.sol';

contract FlashLoan is FlashLoanSimpleReceiverBase {
    using LowGasSafeMath for uint256;
    using LowGasSafeMath for int256;

    address payable owner;

    address public token0 = 0x65aFADD39029741B3b8f0756952C74678c9cEC93;
    address public token1 = 0xCCB14936C2E000ED8393A571D15A2672537838Ad;

    uint24 public fee = 3000;

    ISwapRouter public immutable swapRouter;


    modifier onlyOwner() {
        require(msg.sender == owner, "only owner can call this");
        _;
    }

    constructor(
        address _addressProvider,
        ISwapRouter _swapRouter
    )
        FlashLoanSimpleReceiverBase(IPoolAddressesProvider(_addressProvider))
    {
        owner = payable(msg.sender);
        swapRouter = _swapRouter;
    }

    function requestFlashLoan(address _token, uint256 _amount) public {
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

    /// Calls after flash loan
    function executeOperation(
        address asset,
        uint256 amount,
        uint256 premium,
        address initiator,
        bytes calldata params
    )  external override returns (bool) {
        
        // Swap 

        // Approve the Pool contract allowance to *pull* the owed amount
        uint256 totalAmount = amount + premium;
        IERC20(asset).approve(address(POOL), totalAmount);

        return true;
    }

    receive() external payable {}
}
