FROM golang:1.5
ADD . /go/src/github.com/cpg1111/maestro/
WORKDIR /go/src/github.com/cpg1111/maestro/
RUN echo $PATH && apt-get update && apt-get install -y cmake build-essential pkgconf && make get-deps && make && make install
ENTRYPOINT ["maestro"]
