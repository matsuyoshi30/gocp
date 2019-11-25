# [WIP] gocp

**SUPER ROUGH** Command line tool for the competitive programming written in Go

Currently, only for [AtCoder](https://atcoder.jp/)

### Usage

```sh 
# Login
$ gocp login

# Check session status
$ gocp session

# Make directory and template files.
$ gocp prepare [contest No]

$ cd ./[contest No]/A

# Solve task ( If you use language need to compile, you do it )

# Run test cases (e.g. task A).
$ gocp test a.out

# Submit (e.g. task A).
$ gocp submit main.cpp

# Logout
$ gocp logout
```

### License

MIT
