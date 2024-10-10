pragma solidity 0.8.19;

contract Toolchain {

    event ToolchainEvent(uint256 indexed a, bytes32 indexed b, bytes32 indexed c);

    function setBytes32(uint8 a, bytes32 b, bytes32 c) public {
        emit ToolchainEvent(a, b, c);
    }
}