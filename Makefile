.PHONY: go
default : run
run: 
	go run cmd/multas/*.go --type P --number $(plate)
