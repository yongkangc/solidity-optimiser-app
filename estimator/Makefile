.PHONY: install build estimate

install:
	bash ./install.sh

build:
	forge build --silent

estimate:
	forge test --gas-report

run:
	make build && make estimate
