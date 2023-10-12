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
upgrade:
	forge script --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} script/UpgradeWalleMon.s.sol --broadcast --verify -vvvv

#   [2265977] → new WalleMon@0x1291Be112d480055DaFd8a610b7d1e203891C274
#     ├─ emit Initialized(version: 18446744073709551615 [1.844e19])
#     └─ ← 10128 bytes of code

#   [130457] → new ERC1967Proxy@0x5f3f1dBD7B74C6B46e8c44f98792A1dAf8d69154
#     ├─ emit Upgraded(implementation: WalleMon: [0x1291Be112d480055DaFd8a610b7d1e203891C274])
#     └─ ← 170 bytes of code

CONTRACT=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512

# wallemon calls
feed:
	cast send --private-key ${USER_PRIVATE_KEY} ${CONTRACT} "feed(uint256)" 0
sick:
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "sick(uint256)" 0
heal:
	cast send --private-key ${USER_PRIVATE_KEY} ${CONTRACT} "heal(uint256)" 0
kill:
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "kill(uint256)" 0
# view calls
getHealth:
	cast call ${CONTRACT} "getHealth(uint256)(uint256)" 0
getLastMealTime:
	cast call ${CONTRACT} "getLastMealTime(uint256)(uint256)" 0

# standard calls
token_uri:
	cast call ${CONTRACT} "tokenURI(uint256)(string)" 0

ownerOf:
	cast call ${CONTRACT} "ownerOf(uint256)(address)" 0

balanceOf:
	cast call ${CONTRACT} "balanceOf(address)(uint256)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8

total_supply:
	cast call ${CONTRACT} "totalSupply()(uint256)"

safeMint:
	@echo "Sending transaction to contract ${CONTRACT}"
	@echo "Using private key: ${PRIVATE_KEY}"
	# cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "safeMint(address)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "safeMint(address,string)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 fucnk

owner:
	cast call ${CONTRACT} "owner()(address)"


# playground
getNum:
	cast call ${CONTRACT} "getNum()(uint256)"
setNum:
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "setNum(uint256)" 100
incr:
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "incrementNum()"
add:
	cast call ${CONTRACT} "add(uint256,uint256)" 2 6 | tr -d '"' | cast to-dec
blocktime:
	cast call ${CONTRACT} "blocktime()(uint256)"

# utils
abi:
	forge inspect src/WalleMon.sol:WalleMon abi > ../backend/scripts/w.json

abigen:
	abigen --abi=w.json --pkg=abi --out=../backend/abi/WalleMon.go