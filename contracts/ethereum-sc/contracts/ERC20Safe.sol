//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

// TODO: Rename to ERC20Bridge
contract ERC20Safe {
    struct Deposit {
        address tokenAddress;
        uint256 amount;
        address depositor;
        bytes recipient;
    }

    event ERC20Deposited(    
        uint64 depositIndex
    );

    uint64 public depositsCount;

    mapping(uint64 => Deposit) deposits;
  
  /**

      @notice It assumes that tokenAddress is a corect address for an ERC20 token. No checks whatsoever for this (yet)
      @param tokenAddress Address of the contract for the ERC20 token that will be deposited
      @param amount number of tokens that need to be deposited
      @param recipientAddress address of the r eceiver of tokens on Elrond Network
      @notice emits {ERC20Deposited} event
   */
    function deposit(address tokenAddress, uint256 amount, bytes calldata recipientAddress) public {        
        uint64 depositIndex = ++depositsCount;
        deposits[depositIndex] = Deposit(tokenAddress, amount, msg.sender, recipientAddress);
        
        lockTokens(tokenAddress, msg.sender, amount);
        emit ERC20Deposited(depositIndex);
    }

    function lockTokens(address tokenAddress, uint256 amount, address owner) internal {
        IERC20 erc20 = IERC20(tokenAddress);
        _safeTransferFrom(erc20, owner, address(this), amount);
    }
`
 
    /**
        @notice Gets information about a deposit into the bridge
        @param depositIndex Index of the deposit. Also represents the n-th deposit that was made
        @return Deposit which consists of:
        - tokenAddress Address used when {deposit} was executed.
        - amount Number of tokens that were deposited
        - depositor Address of the account that deposited the tokens
        - recipient Address where tokens will be minted on Elrond Network
    */
    function getDeposit(uint64 depositIndex) external view returns (Deposit memory) {
        return deposits[depositIndex];
    }

    // function _safeTransfer(IERC20 token, address to, uint256 value) private {
    //     _safeCall(token, abi.encodeWithSelector(token.transfer.selector, to, value));
    // }
    
    // function _safeTransferFrom(IERC20 token, address from, address to, uint256 value) private {
    //     _safeCall(token, abi.encodeWithSelector(token.transferFrom.selector, from, to, value));
    // }

    // function _safeCall(IERC20 token, bytes memory data) private {        
    //     (bool success, bytes memory returndata) = address(token).delegatecall(data);
    //     require(success, "ERC20: call failed");

    //     if (returndata.length > 0) {
    //         require(abi.decode(returndata, (bool)), "ERC20: operation did not succeed");
    //     }
    // }
}