tcptee
======

tcptee is a simple tcp traffic duplicator.

Usage
-----

    ./tcptee -bind :8000 -backends :2015,:2016,:2017

Example
-------

    go run echo_server.go -bind :8001
    go run echo_server.go -bind :8002
    go run tcptee.go -bind :8000 -backends :8001,:8002
    echo 'Hello world' | nc 127.0.0.1 8000

Or run this script: [example.sh](./example.sh)

License
-------

BSD.
