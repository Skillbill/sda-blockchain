const Web3 = require('web3');
const cc = require('./compile');
const provider = new Web3.providers.HttpProvider('http://localhost:7545')
const contractOwner = '0x627306090abaB3A6e1400e9345bC60c78a8BEf57';
const candidates = ['Zaphod Beeblebrox', 'Foobar', 'Qux'];

cc.compile('../sol/voting.sol', 'Voting').then(data => {
    const {abi, bytecode} = data;
    const web3 = new Web3(provider);
    const contract = new web3.eth.Contract(JSON.parse(abi));
    return contract.deploy({
        data: bytecode,
        arguments: [ candidates.map(Web3.utils.fromAscii) ]
    }).send({
        from: contractOwner,
        gas: 900000,
        gasPrice: '20000000000' // wei
    });
}).then(receipt => {
    /* ... */
}).catch(err => {
    /* ... */
});
