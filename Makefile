CompileDate = $(shell date +'%Y%m%d%H%M%S')
VersionNumber = v1.0.0
.PHONY: all
all: oplian oplian-gateway oplian-op oplian-op-c2
oplian:
        rm -f oplian
        go build -ldflags "-X main.GitTag=${VersionNumber} -X main.BuildTime=${CompileDate}"  ./cmd/oplian/
oplian-gateway:
        rm -f oplian-gateway
        go build -ldflags "-X main.GitTag=${VersionNumber} -X main.BuildTime=${CompileDate}"  ./cmd/oplian-gateway/
oplian-op:
        rm -f oplian-op
        go build -ldflags "-X main.GitTag=${VersionNumber} -X main.BuildTime=${CompileDate}"  ./cmd/oplian-op/
oplian-op-c2:
        rm -f oplian-op-c2
        go build -ldflags "-X main.GitTag=${VersionNumber} -X main.BuildTime=${CompileDate}"  ./cmd/oplian-op-c2/

.PHONY: clean
clean:
        rm -f oplian oplian-gateway oplian-op oplian-op-c2