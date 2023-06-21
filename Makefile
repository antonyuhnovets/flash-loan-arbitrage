deploy:
	npx hardhat run --network $(network) scripts/deploy.js

verify:
	npx hardhat verify $(address) $(input) --contract 'contracts/$(contract).sol:$(contract)' --network $(network)

build: 
	solc --bin --abi --include-path node_modules/ --base-path . ./contracts/$(contract).sol -o build
	abigen --abi=./build/$(contract).abi --bin=./build/$(contract).bin --pkg=api --out=./internal/api/$(contract).go

delete:
	rm -rf ./build
	rm -rf ./internal/api/$(contract).go