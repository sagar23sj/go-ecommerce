run: ## Run e-commerce poject on host machine
	go run cmd/main.go

clean: ## Clean database file for a fresh start
	rm test.db