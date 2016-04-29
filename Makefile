ifeq ($(OS),Windows_NT)
	LDD_CMD = echo "Pff. Seriously? Windows?"; exit 1;
    CCFLAGS += -D WIN32
    ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        CCFLAGS += -D AMD64
    endif
    ifeq ($(PROCESSOR_ARCHITECTURE),x86)
        CCFLAGS += -D IA32
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        CCFLAGS += -D LINUX
		LDD_CMD = ldd
    endif
    ifeq ($(UNAME_S),Darwin)
        CCFLAGS += -D OSX
		LDD_CMD = otool -l
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif
all: build
get-deps:
	curl -L https://github.com/libgit2/libgit2/archive/v0.22.0.tar.gz > v0.22.0.tar.gz
	tar xzvf v0.22.0.tar.gz
	cd libgit2-0.22.0 && \
	pwd && \
	mkdir build && \
	cd build && \
	pwd && \
	cmake .. && \
	cmake --build . && \
	make install && \
	cd -
	rm -rf libgit2-0.22.0
	rm v0.22.0.tar.gz
	cd ${GOPATH} && \
	go get github.com/tools/godep && \
	cd -
	rm -rf ./Godeps/_workspace/
	godep restore ./...
build:
	rm -rf ./Godeps/_workspace/
	godep restore ./...
	go build --ldflags '-w' -o maestro github.com/cpg1111/maestro/
	$(LDD_CMD) ./maestro
install:
	cp ./maestro /usr/bin/maestro
	mkdir /etc/maestro/
	cp ./test_conf.toml /etc/maestro/conf.toml
clean:
	rm -rf $GOPATH/bin/github.com/cpg1111/maestro $GOPATH/pkg/github.com/cpg1111/maestro $GOPATH/src/github.com/cpg1111/kubongo/maestro
test:
	go test -v ./...
uninstall:
	rm -rf /etc/maestro
	rm -rf /usr/bin/maestro
