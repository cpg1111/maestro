# maestro
Deploy only what's changed for your multiple services in mono-repos

## How it Works

Maestro pulls a given repository then builds a dependency tree based on a given config file.
Once the dependency tree is created, Maestro diffs against a previous commit, with the pathspec being the root directory for each service.
Maestro the flags only the changed services for the pipeline, which is then ran concurrently per teir of dependencies, therefore siblings will build, test and deploy concurrently, but parents and children dependencies will always be built in the correct order.

## Install

```
    go get -d github.com/cpg1111/maestro
    cd maestro
```

or

```
    git clone git@github.com/cpg1111/maestro.git
    cd maestro
```

or

```
    curl -L -o <target path> -z <target path> https://github.com/cpg1111/maestro/releases/download/v0.1.0/maestro-0.1.0-<arch>.tar.gz
    tar -C <target directory> xzvf maestro-0.1.0-<arch>.tar.gz
    cd <target directory>/maestro-0.1.0-<arch>
```

or

```
    wget https://github.com/cpg1111/maestro/releases/download/v0.1.0/maestro-0.1.0-<arch>.tar.gz
    tar -C <target directory> xzvf maestro-0.1.0-<arch>.tar.gz
    cd <target directory>/maestro-0.1.0-<arch>
```

or

```
    docker pull cpg1111/maestro
```

then

dynamically linked

```
    sudo make get-deps # requires libgit2, so get-deps downloads, builds and installs libgit2
    make
    sudo make install
```

or

statically linked

```
    ./build.sh
```

then

```
    maestro --branch <git branch to build> --conf <project config> --prev-commit <commit to compare to> --clone-path <tmp path to clone repo into>
```

or

```
    docker run -v <path of conf>:<target> -v <path to ssh credentials if using ssh for git>:<target> maestro --branch <git branch to build> --conf <project config> --prev-commit <commit to compare to> --clone-path <tmp path to clone repo into>
```
