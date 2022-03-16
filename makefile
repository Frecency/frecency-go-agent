# Go parameters
GOCMD=go
GOGET=$(GOCMD) get

# fetch dependencies
deps:
	$(GOGET) github.com/denisbrodbeck/machineid
	$(GOGET) golan