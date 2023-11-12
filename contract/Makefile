include .env
export $(shell sed 's/=.*//' .env)

.PHONY: test script

# foundry
build:
	forge build
test:
	forge test -vvv

# deploy
deployPreheat:
	if [ "$(network)" = "goerli" ]; then \
		forge script --rpc-url=${GOERLI} --private-key=${GOERLI_PRIVATE_KEY} script/Preheat.s.sol --broadcast --verify -vvvv; \
	elif [ "$(network)" = "arbitrum" ]; then \
		forge script --rpc-url=${ARB} --private-key=${ARB_PRIVATE_KEY} script/Preheat.s.sol --broadcast --verify -vvvv; \
	else \
		forge script --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} script/Preheat.s.sol --broadcast --verify -vvvv; \
	fi
deployWallemon:
	if [ "$(network)" = "goerli" ]; then \
		forge script --rpc-url=${GOERLI} --private-key=${GOERLI_PRIVATE_KEY} script/WalleMon.s.sol --broadcast --verify -vvvv; \
	elif [ "$(network)" = "arbitrum" ]; then \
		forge script --rpc-url=${ARB} --private-key=${ARB_PRIVATE_KEY} script/WalleMon.s.sol --broadcast --verify -vvvv; \
	else \
		forge script --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} script/WalleMon.s.sol --broadcast --verify -vvvv; \
	fi

estimateGasPrice:
	cast gas-price --rpc-url=${ARB}
	
safeMintPreheat:
	cast send --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512 "safeMint(address)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
safeMintWallemon:
	if [ "$(network)" = "goerli" ]; then \
		cast send --rpc-url=${GOERLI} --private-key=${GOERLI_PRIVATE_KEY} 0x78a201027fea5a36ba53bfc122Ff5370bc2625b6 "safeMint(address,string)" 0xf80620C02ef934f760fCBe8e8311BE7027C6201b ""; \
	elif [ "$(network)" = "arbitrum" ]; then \
		cast send --rpc-url=${GOERLI} --private-key=${GOERLI_PRIVATE_KEY} 0xB3545c6bb37E4854d57ABE36db9943dbf795068d "safeMint(address)" 0x0; \
	else \
		cast send --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9 "safeMint(address,string)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 ""; \
	fi

script:
	forge script --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} script/WalleMon.s.sol --broadcast --verify -vvvv
scriptGoerli:
	forge script --rpc-url=${GOERLI} --private-key=${GOERLI_PRIVATE_KEY} script/WalleMon.s.sol --broadcast --verify -vvvv
setReferral:
	if [ "$(network)" = "goerli" ]; then \
		forge script --rpc-url=${GOERLI} --private-key=${GOERLI_PRIVATE_KEY} script/SetReferral.s.sol --broadcast --verify -vvvv; \
	elif [ "$(network)" = "arbitrum" ]; then \
		forge script --rpc-url=${ARB} --private-key=${ARB_PRIVATE_KEY} script/SetReferral.s.sol --broadcast --verify -vvvv; \
	else \
		forge script --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} script/SetReferral.s.sol --broadcast --verify -vvvv; \
	fi

getReferral:
	if [ "$(network)" = "goerli" ]; then \
		cast call  --rpc-url=${GOERLI} 0x6972D5282c530fE0F92797578582fdeb5aC7414D "getReferralAmounts(string)(uint32)" qa; \
	elif [ "$(network)" = "arbitrum" ]; then \
		@echo "not implemented"; \
	else \
		@echo "not implemented"; \
	fi

getReferralOwner:
	cast call  --rpc-url=${LOCALHOST} 0x5FbDB2315678afecb367f032d93F642f64180aa3 "getOwner()(address)"
setReferralOwner:
	cast send --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0 "setOwner(address)()" 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
reveal:
	if [ "$(network)" = "goerli" ]; then \
		cast send --rpc-url=${GOERLI} --private-key=${GOERLI_PRIVATE_KEY} 0xa3C204CF706069E7D4898482D11610903916C575 "setRevealed(bool)()" true; \
	elif [ "$(network)" = "arbitrum" ]; then \
		@echo not yet; \
	else \
		cast send --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9 "setRevealed(bool)()" true; \
	fi
revealed:
	if [ "$(network)" = "goerli" ]; then \
		cast call --rpc-url=${GOERLI} 0xa3C204CF706069E7D4898482D11610903916C575 "revealed()(bool)"; \
	elif [ "$(network)" = "arbitrum" ]; then \
		@echo not yet; \
	else \
		cast call --rpc-url=${LOCALHOST} 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9 "revealed()(bool)"; \
	fi

health:
	if [ "$(network)" = "goerli" ]; then \
		cast call --rpc-url=${LOCALHOST} 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9 "health(uint256)(uint8)" $(tokenID); \
	elif [ "$(network)" = "arbitrum" ]; then \
		@echo not yet; \
	else \
		cast call --rpc-url=${LOCALHOST} 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9 "health(uint256)(uint8)"  $(tokenID); \
	fi

toBeSickList:
	cast call --rpc-url=${LOCALHOST} 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9 "toBeSickList()(uint256[])"

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
	cast call ${CONTRACT} "health(uint256)(uint256)" 0
getLastMealTime:
	cast call ${CONTRACT} "lastMealTime(uint256)(uint256)" 0
healthList:
	cast call 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9 "healthList()(uint8[])"

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
totalSupply:
	if [ "$(network)" = "goerli" ]; then \
		cast call --rpc-url=${GOERLI} 0x0 "totalSupply()(uint256)"; \
	elif [ "$(network)" = "arbitrum" ]; then \
		@echo not yet; \
	else \
		cast call --rpc-url=${LOCALHOST} 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9 "totalSupply()(uint256)"; \
	fi
safeMint:
	cast send --private-key ${PRIVATE_KEY} ${CONTRACT} "safeMint(address,string)" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 ""
safeMintGoerli:
	cast send --rpc-url ${GOERLI} --private-key ${GOERLI_PRIVATE_KEY} 0x78a201027fea5a36ba53bfc122Ff5370bc2625b6 "safeMint(address,string)" ${GOERLI_ADDR} 0
owner:
	cast call ${CONTRACT} "owner()(address)"
# preheat
totalSupplyPreheat:
	cast call --rpc-url=${LOCALHOST} 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512 "totalSupply()(uint256)"
burnPreheat:
	cast send --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512 "burn(uint256)" 0
burnPreheatAll:
	cast send --rpc-url=${LOCALHOST} --private-key=${PRIVATE_KEY} 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512 "burnAll()"

# utils
abi:
	forge inspect src/WalleMon.sol:WalleMon abi > ../backend/config/abi.json
abiPreheat:
	forge inspect src/Preheat.sol:Preheat abi > ../backend/config/abiPreheat.json
abiReferral:
	forge inspect src/Referral.sol:Referral abi > ../backend/config/abiReferral.json

abigen:
	abigen --abi=w.json --pkg=abi --out=../backend/abi/WalleMon.go



#   [1050581] → new Referral@0x5FbDB2315678afecb367f032d93F642f64180aa3
#     └─ ← 4518 bytes of code

#   [1697908] → new Preheat@0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
#     ├─ emit OwnershipTransferred(previousOwner: 0x0000000000000000000000000000000000000000, newOwner: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266)
#     └─ ← 7142 bytes of code
