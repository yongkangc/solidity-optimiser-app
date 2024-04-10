.PHONY: update-submodules start run-frontend run-backend kill-all

update-submodules:
	@echo "Updating submodules..."
	git submodule update --init --recursive
	@echo "Submodules updated successfully!"

start: run-frontend run-backend

run-frontend:
	tmux new-session -d -s frontend "cd frontend && npm run dev"

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