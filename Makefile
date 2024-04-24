.PHONY: build build-optimizer build-backend update-submodules start run-frontend run-backend kill-all install-dependencies optimiser-test

build:
	mkdir -p build
	make build-optimizer
	make build-backend

build-optimizer:
	@echo "Building Optimizer..."
	go build -o ./build/optimizer ./main.go
	@echo "Optimizer built successfully!"

build-backend:
	@echo "Building Backend..."
	go build -o ./build/backend ./backend/main.go
	@echo "Backend built successfully!"

optimiser-test:
	@echo "Running Optimiser tests..."
	go test -v ./...
	@echo "Optimiser tests passed successfully!"


update-submodules:
	@echo "Updating submodules..."
	git submodule update --init --recursive
	@echo "Submodules updated successfully!"

start: run-frontend run-backend

run-frontend:
	tmux new-session -d -s frontend "cd frontend && npm install && npm run dev"

run-backend:
	tmux new-session -d -s backend "cd backend && go run main.go"

view-backend:
	tmux attach-session -t backend

view-frontend:
	tmux attach-session -t frontend

kill-all:
	@echo "Killing all Tmux sessions for frontend and backend..."
	@tmux list-sessions | grep -E 'frontend|backend' | awk -F':' '{print $$1}' | xargs -I{} tmux kill-session -t {}
	@echo "All Tmux sessions have been killed."

# For setting up on a fresh computer
install-dependencies:
	@echo "Installing dependencies..."

	@echo "Installing Git..."
	sudo apt-get update
	sudo apt-get install -y git

	@echo "Installing nvm..."
	curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
	export NVM_DIR="$$HOME/.nvm"
	[ -s "$$NVM_DIR/nvm.sh" ] && \. "$$NVM_DIR/nvm.sh"

	@echo "Installing Node.js v20..."
	nvm install 20
	nvm use 20

	@echo "Installing Golang..."
	wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
	sudo tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
	export PATH=$$PATH:/usr/local/go/bin
	rm go1.22.2.linux-amd64.tar.gz

	@echo "Installing frontend dependencies..."
	cd frontend && npm install

	@echo "Dependencies installed successfully!"

install-solc:
	mkdir -p ~/.solc/releases
	wget -O solc-static-linux https://github.com/ethereum/solidity/releases/download/v0.8.4/solc-static-linux
	chmod +x solc-static-linux
	mv solc-static-linux ~/.solc/releases/solc-v0.8.4
	ln -s ~/.solc/releases/solc-v0.8.4 ~/.solc/solc
	~/.solc/solc --version
