# Makefile

.PHONY: update-submodules

update-submodules:
	@echo "Updating submodules..."
	git submodule update --init --recursive
	@echo "Submodules updated successfully!"


run:
	go run main.go