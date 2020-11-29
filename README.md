#  exec-env

Execute a command using an environment from .env file. `exec-env` also proxies all os signals to child process.

#  Installation
###   Github Release
Visit the [releases page](https://github.com/valenok-husky/exec-env/releases) to download one of the pre-built binaries.

### Go
As secondary option you can use `go install`:
```
go install github.com/valenok-husky/exec-env
```

# Basic Usage
`exec-env -f <path to your env file> <your command with possible launch flags>`

# Example
### Create the environment file `dev.env`
```
PORT=8080
LOG_LVL=debug
PREFIX="app prefix"
```
### Start `printenv` via `exec-env` with our dev.env file
`exec-env -f dev.env printenv`

Printenv is not a part of `exec-env`. Printenv is a unix command. Check the [man page](https://man7.org/linux/man-pages/man1/printenv.1.html).

The output is
```
<your own environment variables>
...
PORT=8080
LOG_LVL=debug
PREFIX=app prefix
```
It's working 🎉
