# Import dumps from Personagraph

## Build

1. Install Go language (https://golang.org/doc/install)

    `$ brew install go`

2. Install build subsystem (https://getgb.io/)

    `$ go get github.com/constabulary/gb/...`

3. Fetch dependencies

    `$ gb vendor restore`

4. Build project

    `$ gb build`


## Usage

1. Import dumps in Gzip format
    
    `$ ./bin/import dumps/*.gz`
