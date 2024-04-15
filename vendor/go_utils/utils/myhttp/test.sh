GOFLAGS="-count=1" go test -v -run TestMutualTLSServer
GOFLAGS="-count=1" go test -v -run TestGenerateCerts
GOFLAGS="-count=1" go test -v -run TestClearCerts
