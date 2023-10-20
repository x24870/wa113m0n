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
setReferral:
	forge script --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} script/SetReferral.s.sol --broadcast --verify -vvvv
getReferral:
	cast call  --rpc-url=${LOCALHOST} 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0 "getReferralAmounts(string)(uint32)" wallemon
getReferralOwner:
	cast call  --rpc-url=${LOCALHOST} 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0 "getOwner()(address)"
setReferralOwner:
	cast send --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0 "setOwner(address)()" 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

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



#   [2568795] → new WalleMon@0x5FbDB2315678afecb367f032d93F642f64180aa3
#     ├─ emit Initialized(version: 18446744073709551615 [1.844e19])
#     └─ ← 11531 bytes of code

#   [135226] → new ERC1967Proxy@0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
#     ├─ emit Upgraded(implementation: WalleMon: [0x5FbDB2315678afecb367f032d93F642f64180aa3])
#     └─ ← 177 bytes of code

#   [957022] → new Referral@0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0
#     └─ ← 4084 bytes of code

