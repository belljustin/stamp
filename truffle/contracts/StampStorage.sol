pragma solidity ^0.4.21;

contract StampStorage {
    struct Stamp {
        uint timestamp;
        bytes32 hash;
    }

    uint stampCounter = 0;

    address public stamper;
    mapping (uint => Stamp) public stamps;

    constructor() public {
        stamper = msg.sender;
    }

    event Stamped(uint counter, uint timestamp, bytes32 hash);

    function addStamp(bytes32 hash) public returns (uint) {
        if (msg.sender != stamper) return;
        uint timestamp = now;
        Stamp memory stamp = Stamp(timestamp, hash);
        stamps[stampCounter] = stamp;
        emit Stamped(stampCounter++, timestamp, hash);
        return stampCounter;
    }
}
