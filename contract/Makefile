include .env
export $(shell sed 's/=.*//' .env)

.PHONY: test script

echo:
	@echo LOCALHOST is $(LOCALHOST)
	@echo PRIVATE_KEY is $(PRIVATE_KEY)

build:
	forge build

test:
	forge test -vv

run_local_node:
	anvil

deploy_local:
	forge create --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} WalleMon --constructor-args ${ADDR}
	# forge create --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} WalleMon

script:
	forge script --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} script/WalleMon.s.sol --broadcast --verify -vvvv

#   [1679136] → new WalleMon@0x322813Fd9A801c5507c9de605d63CEA4f2CE6c44
#     ├─ emit Initialized(version: 18446744073709551615 [1.844e19])
#     └─ ← 7410 bytes of code

#   [130469] → new ERC1967Proxy@0xa85233C63b9Ee964Add6F2cffe00Fd84eb32338f
#     ├─ emit Upgraded(implementation: WalleMon: [0x322813Fd9A801c5507c9de605d63CEA4f2CE6c44])
#     └─ ← 170 bytes of code

CONTRACT=0xa85233C63b9Ee964Add6F2cffe00Fd84eb32338f
getNum:
	cast call ${CONTRACT} "getNum()(uint256)"
setNum:
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "setNum(uint256)" 100
incr:
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "incrementNum()"
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

safeMint:
	@echo "Sending transaction to contract ${CONTRACT}"
	@echo "Using private key: ${PRIVATE_KEY}"
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "safeMint(address)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
	# cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "safeMint(address,string)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 fucnk

owner:
	cast call ${CONTRACT} "owner()(address)"

abi:
	forge inspect src/WalleMon.sol:WalleMon abi > ../backend/scripts/w.json

abigen:
	abigen --abi=w.json --pkg=abi --out=../backend/abi/WalleMon.go