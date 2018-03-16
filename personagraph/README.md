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

- Import profiles from dumps in Gzip format
    
    `$ ./bin/import dumps/*.gz`

- Import profiles filtered by categories from dumps in Gzip format

    `$./bin/import -f with-gender:55,56 -f with-age:40,41,42,43,44,45,46 ./dumps/*.gz`

- Other options can be found with -h flag

   `$ ./bin/import -h`
