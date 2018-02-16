const fs = require('fs');
const solc = require('solc'); // npm install solc

module.exports = {
	compile: function(path, klass) {
		return new Promise((resolve, reject) => {
			fs.readFile(path, (err, data) => {
				if (err) {
					return reject(err);
				}
				const output = solc.compile(data.toString());
				const compiledContract = output.contracts[`:${klass}`];
				if (compiledContract === undefined) {
					return reject(`no contract class "${klass}" in "${path}"`);
				}
				resolve({
					abi: compiledContract.interface,
					bytecode: compiledContract.bytecode
				});
			});
		});
	}
};
