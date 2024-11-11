# Swagger Implementation
* Add the necessary code under the pkg directory. 
* Install the required Go packages with the following command:
``` bash
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/fiber-swagger
```
* Then, run the following init command, specifying the folder of the relevant main.go file:
```bash
swag init -g cmd/server/main.go -o docs
```
---
# Unit Test Implementation
* Install the `gomock` package.
``` bash
go get github.com/golang/mock/gomock
```
* Then, run the command to generate mocks for the service or repositories you need to mock.
* The following example creates a mock for a service:
``` bash
mockery --name=AddressService --dir=internal/address/service --output=internal/address/service/mocks
```
* If the `mockgen` command is not installed, use one of these commands to install it:
``` bash
go install github.com/golang/mock/mockgen@latest
or
go install github.com/vektra/mockery/v2@latest
```
* The following example creates a mock for a repository.
``` bash
mockery --name=AddressRepository --dir=internal/address/repository --output=internal/address/repository/mocks
```