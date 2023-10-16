include .env
export $(shell sed 's/=.*//' .env)

.PHONY: test script

# foundry
build:
	forge build
test:
	forge test -vvv

# deploy
script:
	forge script --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} script/WalleMon.s.sol --broadcast --verify -vvvv
scriptGoerli:
	forge script --rpc-url=${GOERLI} --private-key=${GOERLI_PRIVATE_KEY} script/WalleMon.s.sol --broadcast --verify -vvvv
# upgrade
upgrade:
	forge script --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} script/UpgradeWalleMon.s.sol --broadcast --verify -vvvv

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

# uri calls
tokenUri:
	cast call ${CONTRACT} "tokenURI(uint256)(string)" 0
tokenUriGoerli:
	cast call --rpc-url ${GOERLI} ${PROXY} "tokenURI(uint256)(string)" 0
setEggUri:
	cast send --rpc-url ${GOERLI} --private-key ${GOERLI_PRIVATE_KEY} ${PROXY} "setEggURI(string)" "https://ipfs.blocto.app/ipfs/QmUNPQoX9JzA2kmaUvVKFPmsuQY9vmY4ytwPrNZykQSd8Z"
# ERC721 calls
ownerOf:
	cast call ${CONTRACT} "ownerOf(uint256)(address)" 0
balanceOf:
	cast call ${CONTRACT} "balanceOf(address)(uint256)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
total_supply:
	cast call ${CONTRACT} "totalSupply()(uint256)"
safeMint:
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "safeMint(address,string)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 fucnk
safeMintGoerli:
	cast send --rpc-url ${GOERLI} --private-key ${GOERLI_PRIVATE_KEY} ${PROXY} "safeMint(address,string)" ${GOERLI_ADDR} 0
owner:
	cast call ${CONTRACT} "owner()(address)"

# utils
abi:
	forge inspect src/WalleMon.sol:WalleMon abi > ../backend/scripts/w.json

abigen:
	abigen --abi=w.json --pkg=abi --out=../backend/abi/WalleMon.go