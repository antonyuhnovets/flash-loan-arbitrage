deploy:
	npx hardhat run --network $(network) scripts/deploy$(contract).js

verify:
	npx hardhat verify $(address) $(token) --contract 'contracts/$(contract).sol:$(contract)' --network $(network)
