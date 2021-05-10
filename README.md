# adjust tool
This tool makes http requests and prints the address of the request along with the MD5 hash of the response.

## Run
To run the tool you need to execute `./adjust <url>...` and in case that for any reason the tool gets an error when executing the request the url will be printed with an empty space
For example:
``` sh
l9431@macl4581 adjust % ./adjust google.com adjust.com http://invalidUrl                                 
http://google.com ed65957bdf10cc9049e9a049051eff8b
http://adjust.com 787d1722dc1a1bd33029b8b9e13a0e1e
http://invalidUrl 
```

## Parallel flag
By default the tool does 10 parallel requests but this can be changed using the `-parallel` flag
For example:
``` sh
l9431@macl4581 adjust % ./adjust -parallel google.com adjust.com
```

