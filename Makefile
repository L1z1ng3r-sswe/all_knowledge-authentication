SRC = app/cmd/main.go
DESTINATION = app/cmd/build
BINARY_APP = binary_app
	
build:
	go build -o $(DESTINATION)/$(BINARY_APP) $(SRC)

run: build
	./$(DESTINATION)/$(BINARY_APP)

clean: 
	rm -rf cmd/build/*