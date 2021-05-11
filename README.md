# hash-tool
This tool makes http requests and prints the address of the request along with the MD5 hash of the response.

## Build
To build this tool run `go build .`

## Run
To run the tool you need to execute `./hash-tool <url>...` and in case that for any reason the tool gets an error when executing the request the url will be printed with an empty space
For example:
``` sh
l9431@macl4581 % ./hash-tool google.com http://invalidUrl                                 
http://google.com ed65957bdf10cc9049e9a049051eff8b
http://invalidUrl 
```

## Parallel flag
By default the tool does 10 parallel requests but this can be changed using the `-parallel` flag
For example:
``` sh
l9431@macl4581 % ./hash-tool -parallel google.com
```

