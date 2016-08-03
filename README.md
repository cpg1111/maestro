# maestro
Deploy only what's changed for your multiple services in mono-repos

[![Go Report Card](https://goreportcard.com/badge/github.com/cpg1111/maestro)](https://goreportcard.com/report/github.com/cpg1111/maestro) [![Code Climate](https://codeclimate.com/github/cpg1111/maestro/badges/gpa.svg)](https://codeclimate.com/github/cpg1111/maestro)

## How it Works

Maestro pulls a given repository then builds a dependency tree based on a given config file.
Once the dependency tree is created, Maestro diffs against a previous commit, with the pathspec being the root directory for each service.
Maestro the flags only the changed services for the pipeline, which is then ran concurrently per teir of dependencies, therefore siblings will build, test and deploy concurrently, but parents and children dependencies will always be built in the correct order.

For more details see this talk: https://www.youtube.com/watch?v=dGM8mYj8nz4&feature=youtu.be

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
    # This one handles everything but running maestro
    docker pull cpg1111/maestro
```

then

```
    docker build -t maestro_c -f Dockerfile_c .
    docker build -t maestro_build -f Dockerfile_build .
```

then either:

```
    docker build -t maestro .
```

or

```
    docker build -t maestro_bin_deps -f Dockerfile_bin .
    docker build -t maestro -f Dockerfile_fully_loaded . # gives you certain common executables for building, testing and deploying
```

## Test

```
    go test ./...
```

## Run

```
    Usage of maestro:
      --branch string
            Git branch to checkout for project (default "master")
      --clone-path string
            Local path to clone repo to defaults to PWD (default "./")
      --config string
            Path to the config for maestro to use (default "./conf.toml")
      --deploy
            Whether or not to deploy this build                             # defaults to false
      --prev-commit string
            Previous commit to compare to                                   # required
```

```
    maestro --branch <git branch to build> --conf <project config> --prev-commit <commit to compare to> --deploy <whether to deploy build or not> --clone-path <tmp path to clone repo into>
```

or

```
    docker run -v <path of conf>:<target> -v <path to ssh credentials if using ssh for git>:<target> maestro --branch <git branch to build> --conf <project config> --prev-commit <commit to compare to> --deploy <whether to deploy build or not> --clone-path <tmp path to clone repo into>
```

## Example Config

```
    [Environment] # Environment will run before anything else, ExecSync will execute commands in the array synchronously, while exec will execute them concurrently
    Env=["node_env:test", "docker_tls_verify:1"] # set environment variables all lowercase keys and keys are separated from values with ':'
    ExecSync=["apt-get install -y docker node go"]
    Exec=["docker pull someOrg/logger", "docker pull someOrg/models", "docker pull someOrg/auth", "docker pull someOrg/client"]

    [Project]
    RepoURL="git@github.com:someOrg/someRepo.git"
    CloneCMD="git clone"
    AuthType="SSH"
    SSHPrivKeyPath="~/.ssh/id_rsa"
    SSHPubKeyPath="~/.ssh/id_rsa.pub"
    Username="git" # github's ssh user is git, but this can vary
    Password=""
    PromptForPWD=false # when requiring a password, you prompt for a password

    [[Services]] # Services are either actual services or libraries / packages / separately compiled objects
    Name="logger"
    Tag="0.1.0"
    TagType="git"
    Path="./src/logger"
    BuildCMD=["docker build -t logger ."] # '.' is relative to the given path field of the service
    TestCMD=["go test ./..."]
    CheckCMD=["bash -c 'docker images -a | grep logger'"]
    CreateCMD=[
        "docker tag logger <org>/logger:{{.Curr}}",
        "docker push <org>/logger:{{.Curr}}"
    ]
    UpdateCMD=[
        "docker tag logger <org>/logger:{{.Curr}}",  # {{.Curr}} will template the current commit hash into the command
        "docker push <org>/logger:{{.Curr}}"
    ]
    DependsOn=[]

    [[Services]]
    Name="models"
    Tag="0.1.0"
    TagType="git"
    Path="./src/models"
    BuildCMD=["docker build -t models ."]
    TestCMD=["go test ./..."]
    CheckCMD=["bash -c 'docker images -a | grep models'"]
    CreateCMD=[
        "docker tag models <org>/models:{{.Curr}}",
        "docker push <org>/models:{{.Curr}}"
    ]
    UpdateCMD=[
        "docker tag models <org>/models:{{.Curr}}",
        "docker push <org>/models:{{.Curr}}"
    ]
    DependsOn=["logger"] # Assume Dockerfile contains FROM logger

    [[Services]]
    Name="auth"
    Tag="0.1.0"
    TagType="git"
    Path="./src/auth"
    BuildCMD=["docker build -t auth ."]
    TestCMD=["go test ./..."]
    CheckCMD=["bash -c 'docker ps -a | grep auth'"]
    CreateCMD=[
        "docker tag auth <org>/auth:{{.Curr}}",
        "docker tag auth <org>/auth:{{.Curr}}",
        "docker run --rm -d <org>/auth:{{.Curr}}"
    ]
    UpdateCMD=[
        "docker tag auth <org>/auth:{{.Curr}}",
        "docker tag auth <org>/auth:{{.Curr}}",
        "docker run --rm -d <org>/auth:{{.Curr}}"
    ]
    DependsOn=["database"] # Assume Dockerfile contains FROM models

    [[Services]]
    Name="client"
    Tag="0.1.0"
    TagType="git"
    Path="./src/client"
    BuildCMD=["docker build -t client ."]
    TestCMD=["npm test"]
    CheckCMD=["bash -c 'docker ps -a | grep client'"]
    CreateCMD=[
        "docker tag client <org>/client:{{.Curr}}",
        "docker push <org>/client:{{.Curr}}",
        "docker run --rm -d <org>/client:{{.Curr}}"
    ]
    UpdateCMD=[
        "docker tag client <org>/client:{{.Curr}}",
        "docker push <org>/client:{{.Curr}}",
        "docker run --rm -d <org>/client:{{.Curr}}"
    ]
    DependsOn=[]

    [CleanUp]
    AdditionalCMDs=["docker inspect auth", "docker export -o ./dist/auth.tgz auth"] # Will execute synchronously
    InDaemon=false # COMING SOON for maestrod
        [[CleanUp.Artifacts]] # Artifacts are saved concurrently
        RuntimeFilePath="./dist/auth.tgz"
        SaveFilePath="/opt/data/auth.tgz"
```

## Roadmap

- Pipeline tests
- Allow larger log buffers
- More possible dependency structures
- Maestrod integration on clean up
- Encrypted Environment variable values
- Log Versbosity control
- Debug with bash session

### Daemon
See [this](https://github.com/cpg1111/maestrod) (https://github.com/cpg1111/maestrod) for a manager daemon for handling git push hooks and multiple concurrent builds and repos.

### Warning

Some dependency structures are not supported yet.  Best to refer to immediate dependencies only and avoid circlular references.
