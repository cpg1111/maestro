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
	mkdir -p tmp && \
	cd tmp && \
	curl -o zlib-1.2.8.tar.gz -z zlib-1.2.8.tar.gz http://zlib.net/zlib-1.2.8.tar.gz && \
    tar xzvf zlib-1.2.8.tar.gz && \
    cd zlib-1.2.8 && \
    ./configure && \
    make && make install && \
    curl -o openssl-1.0.2h.tar.gz -z openssl-1.0.2h.tar.gz https://openssl.org/source/openssl-1.0.2h.tar.gz && \
    tar xzvf openssl-1.0.2h.tar.gz && \
    cd openssl-1.0.2h && \
    ./config --prefix=/usr \
             --openssldir=/etc/ssl \
             --libdir=lib \
             shared \
             zlib-dynamic && \
    make depend && \
    make && make install && \
	curl -L -o http-parser.tar.gz -z http-parser.tar.gz https://github.com/nodejs/http-parser/archive/v2.7.0.tar.gz && \
	tar xzvf http-parser.tar.gz && \
	cd http-parser-2.7.0 && \
	PREFIX=/usr make package && PREFIX=/usr/ make install && ls /usr/include/ && ls /usr/lib/
	curl -L -o libssh.tar.gz -z libssh.tar.gz https://www.libssh2.org/download/libssh2-1.4.2.tar.gz
	tar xzvf libssh.tar.gz
	cd libssh2-1.4.2 && \
	./configure && \
	make && make install
	curl -L -o v0.22.0.tar.gz -z v0.22.0.tar.gz https://github.com/libgit2/libgit2/archive/v0.22.0.tar.gz
	tar xzvf v0.22.0.tar.gz
	cd libgit2-0.22.0 && \
	pwd && \
	mkdir build && \
	cd build && \
	pwd && \
	cmake .. \
		-DCMAKE_INSTALL_PREFIX=/usr/ \
		-DTHREADSAFE=ON \
	    -DBUILD_CLAR=OFF \
		&& \
	cmake --build . --target install && \
	cd -
	rm -rf libgit2-0.22.0
	rm v0.22.0.tar.gz
	cd ${GOPATH} && \
	go get -u github.com/kardianos/govendor && \
	cd - && \
	govendor sync
	rm -rf tmp
build:
	govendor sync
	go build -linkshared -o maestro github.com/cpg1111/maestro/
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
docker:
	docker build -t maestro_c -f Dockerfile_c .
	docker build -t maestro_build -f Dockerfile_build .
	docker build -t maestro_bin_deps -f Dockerfile_bin .
	docker build -t maestro -f Dockerfile_fully_loaded .
e2e-test:
	docker run --rm -it \
	-v ${HOME}/.ssh/:/root/.ssh/ \
	-v `pwd`:/etc/maestro/ \
	-v ${DOCKER_CERT_PATH}:${DOCKER_CERT_PATH} \
	-e DOCKER_HOST=${DOCKER_HOST} \
	-e DOCKER_MACHINE_NAME=${DOCKER_MACHINE_NAME} \
	-e DOCKER_TLS_VERIFY=1 \
	-e DOCKER_CERT_PATH=${DOCKER_CERT_PATH} \
	maestro \
	--clone-path=/tmp/build \
	--branch=master \
	--prev-commit=ca30ac184cd46fc1c7d59d7973f87350050e39ee \
	--curr-commit=eaeca0254dc1bd04413f7823fe03f81583ed6b9c \
	--config=/etc/maestro/test_conf.toml \
	--deploy=true
