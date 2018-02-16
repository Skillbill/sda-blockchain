fs = require('fs')
solc = require('solc')
Web3 = require('web3')

let web3
let compiled
let contract
let candidates
let contractAddress

function send(tx, account) {
    return tx.send({
	from: account,
	gas: 900000,
        gasPrice: '20000000000'	
    }).on('receipt', console.log)
}

const init = (path, providerAddr) => {
    provider = new Web3.providers.HttpProvider(providerAddr || 'http://localhost:7545')
    web3 = new Web3(provider)
    const src = fs.readFileSync(path || '../sol/voting.sol').toString()    
    compiled = solc.compile(src)
    const abi = JSON.parse(compiled.contracts[':Voting'].interface)    
    contract = new web3.eth.Contract(abi)
}

const deployContract = (candidates, account) => {
    const tx = contract.deploy({
	from: account,
        data: compiled.contracts[':Voting'].bytecode,
        arguments: [ candidates.map(Web3.utils.asciiToHex) ]
    })
    return send(tx, account).then((receipt) => {
	contract._address = receipt._address
    })
}

const totalVotesFor = (candidate, account) => {
    contract.methods.totalVotesFor(Web3.utils.fromAscii(candidate)).call()
	.then(count => console.log(`total votes for ${candidate}: ${count}`))
}

const voteFor = (candidate, account) => {
    const tx = contract.methods.voteFor(Web3.utils.fromAscii(candidate))
    return send(tx, account)
}

const _getContract = () => { return contract }

module.exports = {
    init,
    deployContract,
    voteFor,
    totalVotesFor,
    _getContract
}
