include .env
export $(shell sed 's/=.*//' .env)

echo:
	@echo LOCALHOST is $(LOCALHOST)
	@echo PRIVATE_KEY is $(PRIVATE_KEY)

build:
	forge build

test:
	forge test

run_local_node:
	anvil

deploy_local:
	forge create --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} WalleMon

CONTRACT=0x0165878A594ca255338adfa4d48449f69242Eb8F
add:
	cast call ${CONTRACT} "add(uint256,uint256)" 2 6 | tr -d '"' | cast to-dec

token_uri:
	cast call ${CONTRACT} "tokenURI(uint256)(string)" 0

token_owner:
	cast call ${CONTRACT} "ownerOf(uint256)(address)" 3

balanceOf:
	cast call ${CONTRACT} "balanceOf(address)(uint256)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8

total_supply:
	cast call ${CONTRACT} "totalSupply()(uint256)"

safe_mint:
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "safeMint(address)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8

abi:
	forge inspect src/WalleMon.sol:WalleMon abi > w.json

abigen:
	abigen --abi=w.json --pkg=abi --out=../backend/abi/WalleMon.go

