#!/bin/bash

# For build docker container's ENTRYPOINT

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
go get github.com/tools/godep && \
cd $GOPATH/src/github.com/cpg1111/maestro/ && \
rm -rf ./Godeps/_workspace/ && \
godep restore ./... && \
go build --ldflags '-extldflags "-static"' -o maestro github.com/cpg1111/maestro/ && \
ldd ./maestro && \
cp ./maestro /opt/bin/
