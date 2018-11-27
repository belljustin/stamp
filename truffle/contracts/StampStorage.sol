pragma solidity ^0.4.21;

contract StampStorage {
    address public stamper;
    mapping (bytes32 => uint) public stamps;

    constructor() public {
        stamper = msg.sender;
    }

    event Stamped(bytes32 hash, uint timestamp);

    function addStamp(bytes32 hash) public {
        if (msg.sender != stamper) return;
        if (stamps[hash] > 0) return;
        uint timestamp = now;
        emit Stamped(hash, timestamp);
    }
}
