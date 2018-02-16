pragma solidity ^0.4.18;

contract Voting {
  mapping (bytes32 => uint8) public votes;
  mapping (bytes32 => bool) public validCandidates;
  
  function Voting(bytes32[] candidates) public {
    for (uint i=0; i<candidates.length; i++) {
       validCandidates[candidates[i]] = true;
    }
  }

  function totalVotesFor(bytes32 candidate) view public returns (uint8) {
    require(isValid(candidate));
    return votes[candidate];
  }

  function voteFor(bytes32 candidate) public {
    require(isValid(candidate));
    votes[candidate] += 1;
  }

  function isValid(bytes32 candidate) view public returns (bool) {
     return validCandidates[candidate];
  }
}
