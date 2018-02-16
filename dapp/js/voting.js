const abi = JSON.parse('[{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"votes","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"candidate","type":"bytes32"}],"name":"totalVotesFor","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"validCandidates","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"candidate","type":"bytes32"}],"name":"voteFor","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"candidate","type":"bytes32"}],"name":"isValid","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"candidates","type":"bytes32[]"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]');
const bytecode = '6060604052341561000f57600080fd5b60405161037b38038061037b8339810160405280805182019190505060008090505b8151811015610091576001806000848481518110151561004d57fe5b906020019060200201516000191660001916815260200190815260200160002060006101000a81548160ff0219169083151502179055508080600101915050610031565b50506102d9806100a26000396000f30060606040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632b38cd96146100725780632f265cf7146100b35780633390b9c2146100f457806335154986146101335780636a9385671461015a575b600080fd5b341561007d57600080fd5b610097600480803560001916906020019091905050610199565b604051808260ff1660ff16815260200191505060405180910390f35b34156100be57600080fd5b6100d86004808035600019169060200190919050506101b9565b604051808260ff1660ff16815260200191505060405180910390f35b34156100ff57600080fd5b6101196004808035600019169060200190919050506101fe565b604051808215151515815260200191505060405180910390f35b341561013e57600080fd5b61015860048080356000191690602001909190505061021e565b005b341561016557600080fd5b61017f60048080356000191690602001909190505061027b565b604051808215151515815260200191505060405180910390f35b60006020528060005260406000206000915054906101000a900460ff1681565b60006101c48261027b565b15156101cf57600080fd5b600080836000191660001916815260200190815260200160002060009054906101000a900460ff169050919050565b60016020528060005260406000206000915054906101000a900460ff1681565b6102278161027b565b151561023257600080fd5b6001600080836000191660001916815260200190815260200160002060008282829054906101000a900460ff160192506101000a81548160ff021916908360ff16021790555050565b600060016000836000191660001916815260200190815260200160002060009054906101000a900460ff1690509190505600a165627a7a7230582038f6cef23bf86ff8cb46b64d9741b334909152934b75525fca7de68a70ca2abb0029';

let web3
let contract
let candidates

function send(tx, account) {
    return tx.send({
	from: account,
	gas: 900000,
        gasPrice: '20000000000'	
    });
}

function init(providerAddress) {
    const provider = new Web3.providers.HttpProvider(providerAddress || 'http://127.0.0.1:7545')
    web3 = new Web3(provider);
    contract = new web3.eth.Contract(abi);
}

function deploy(account, candidates) {
    const tx = contract.deploy({
	from: account,
        data: bytecode,
        arguments: [ candidates.map(Web3.utils.asciiToHex) ]
    });
    return send(tx, account).then(receipt => {
        // FIXME: I don't think it's the right way to do that
	contract._address = receipt._address;
        return receipt._address;
    });
}

function totalVotesFor(candidate) {
    return contract.methods.totalVotesFor(Web3.utils.fromAscii(candidate)).call();
}

function voteFor(candidate, account) {
    const tx = contract.methods.voteFor(Web3.utils.fromAscii(candidate))
    return send(tx, account)
}
