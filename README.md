# maestro
Deploy only what's changed for your multiple services in mono-repos

## How it Works

Maestro pulls a given repository then builds a dependency tree based on a given cofig file.
Once the dependency tree is create, Maestro then diffs against a previous commit, with the pathspec being the root directory for each service.
Maestro the flags only the changed services for the pipeline, which is then ran concurrently per teir of dependencies, therefore siblings will build, test and deploy concurrently, but parents and children dependencies will always be built in the correct order.

## Install

`curl -L -o <target path> -z <target path> https://github.com/cpg1111/maestro/releases/download/v0.1.0/maestro-0.1.0-<arch>.tar.gz`

or

`wget https://github.com/cpg1111/maestro/releases/download/v0.1.0/maestro-0.1.0-<arch>.tar.gz`

then

```
    tar -C <target directory> xzvf maestro-0.1.0-<arch>.tar.gz
    cd <target directory>/maestro-0.1.0-<arch>
    sudo make get-deps
    make
    sudo make install
```
