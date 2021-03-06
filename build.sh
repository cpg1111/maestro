#!/bin/bash

# For build docker container's ENTRYPOINT

CMAKE_INCLUDE_PATH=$CMAKE_INCLUDE_PATH:/usr/include/:/usr/local/include/
CMAKE_LIBRARY_PATH=$CMAKE_LIBRARY_PATH:/usr/lib/:/usr/local/lib/
OPENSSL_INCLUDE_DIR=/usr/include/
OPENSSL_OPENSSL_LIBRARIES=/usr/lib/

curl -L -o v0.22.0.tar.gz -z v0.22.0.tar.gz https://github.com/libgit2/libgit2/archive/v0.22.0.tar.gz && \
tar xzvf v0.22.0.tar.gz && \
cd libgit2-0.22.0 && \
pwd && \
mkdir build && \
cd build && \
pwd && \
cmake -DTHREADSAFE=ON \
      -DBUILD_CLAR=OFF \
      -DBUILD_SHARED_LIBS=OFF \
      -DCMAKE_C_FLAGS=-fPIC \
      -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
      .. && \
cmake --build . && \
make install && \
mkdir -p /usr/lib/lib/git2 && \
mv /usr/local/include/git2/ /usr/include/git2/ && \
mv /usr/local/include/git2.h /usr/include/git2.h && \
cd $GOPATH/src/github.com/cpg1111/maestro/ && \
rm -rf libgit2-0.22.0 && \
rm v0.22.0.tar.gz && \
cd $GOPATH && \
go get github.com/kardianos/govendor && \
cd $GOPATH/src/github.com/cpg1111/maestro/ && \
govendor sync && \
go build --ldflags '-extldflags "-static"' -o maestro github.com/cpg1111/maestro/ && \
cp ./maestro /opt/bin/
