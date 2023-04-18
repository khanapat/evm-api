// SPDX-License-Identifier: MIT

pragma solidity 0.8.13;

contract GetSet {
    uint256 public a;

    event SetA(uint256 a);

    function setA(uint256 _a) public {
        a = _a;

        emit SetA(_a);
    }
}
