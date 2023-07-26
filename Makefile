run:
	build/app/webapi

# Build
build-up:
	docker compose build
	docker compose up -d

build-down:
	docker compose down

build-rm:
	sudo rm -rf pg_data pg_data_test
	docker rm flash-postgres
	docker rm flash-webapi

build-restart:
	docker compose down 
	docker compose up -d

build-reload:
	docker compose down
	sudo rm -rf pg_data pg_data_test
	docker compose build
	docker compose up -d

# App
app-start:
	docker start flash-webapi
	docker exec flash-webapi make run

app-stop:
	docker stop flash-webapi

app-restart:
	make app-stop
	make app-start

app-reload:
	make app-stop
	make go-build
	make app-start

app-build:
	make go-build

app-logs:
	docker logs flash-webapi

# Database
db-logs:
	docker logs flash-pg

# Go
go-build:
	swag init -g internal/delivery/rest/v1/router.go
	go build -a -o ./build/app/webapi ./cmd/daemon/main.go
	chmod +x ./build/app/webapi
	
# Contract
contract-build: 
	solc --bin --abi --include-path node_modules/ --base-path . ./contracts/$(contract).sol -o build/contract
	abigen --abi=./build/contract/$(contract).abi --bin=./build/contract/$(contract).bin --pkg=api --out=./internal/api/$(contract).go

contract-deploy:
	npx hardhat run --network $(network) scripts/deploy.js

contract-verify:
	npx hardhat verify $(address) $(input) --contract 'contracts/$(contract).sol:$(contract)' --network $(network)

# Install required pkg/dep
install-dependencies:
	go mod download
	npm install

install-node:
	sudo apt install -U nodejs
	sudo apt install -U npm

install-solc:
	sudo apt install -y libboost-all-dev cmake
	git clone --recursive --branch v0.8.10 https://github.com/ethereum/solidity.git ~/github.com/ethereum/solidity \
	&& sh ~/github.com/ethereum/solidity/scripts/build.sh

install-geth:
	git clone https://github.com/ethereum/go-ethereum.git ~/github.com/ethereum/go-ethereum
	cd ~/github.com/ethereum/go-ethereum/ && make devtools

install-protoc:
	sudo apt install -U unzip curl
	PROTOC_ZIP=protoc-3.14.0-linux-x86_64.zip
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$PROTOC_ZIP
	sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
	sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
	rm -f $PROTOC_ZIP
