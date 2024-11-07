# Swagger Implementasyonu
* Pkg altına gerekli kodlar eklenir. 
* Aşağıdaki komut ile gerekli go paketleri yüklenir.
``` bash
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/fiber-swagger
```
* Daha sonra ilgili main.go'nun folder'ını göstererek aşağıdaki init komutu çalıştırılmalı.
```bash
swag init -g cmd/server/main.go -o docs
```
---
# Unit Test Implementasyonu
* gomock paketini kurmalısın.
``` bash
go get github.com/golang/mock/gomock
```
* Daha sonra mocklanacak servis, repositoryler için gomock için komut çalıştırman gerekli.
* Aşağıdaki kod, örnek bir servis için mock işlemi yapar.
``` bash
mockery --name=AddressService --dir=internal/address/service --output=internal/address/service/mocks
```
* mockgen komutu yüklü değilse aşağıdaki kodla yükle.
``` bash
go install github.com/golang/mock/mockgen@latest
veya
go install github.com/vektra/mockery/v2@latest
```
* Aşağıdaki kod, örnek bir repository için mock işlemi yapar.
``` bash
mockery --name=AddressRepository --dir=internal/address/repository --output=internal/address/repository/mocks
```