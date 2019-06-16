# README for go-nc

## What is this?
This is a simple clone of `netcat` in golang. I wrote it many moons
ago as my first go program. It is slightly different from the
traditional `netcat` in the following ways:

* adds sha256 checksum on network I/O
* traffic stats (number of bytes sent/received)
* hexdump of the traffic
* Bi directional I/O is optional; server always writes to `stdout`
  and client always reads from `stdin`. This makes it easy to
  remember.

## How do I build it?
You will need Go 1.5 or later:
 
    git clone https://github.com/opencoff/go-nc
    cd go-nc
    make

## How do I run it?
`gonc` comes with helpful commandline options. Try:

    ./bin/gonc -h

Or the long-form:

    ./bin/gonc --help

### Simple Example
Lets say that we wish to send a directory's contents to a server whose IP address is a.b.c.d.
On that server create a listening instance of 'gonc':

    ./bin/gonc -c -l :9090 > foo.tar.gz

And on a client try:

    tar cf - . | gzip -9 | ./bin/gonc -c -v a.b.c.d:9090


In the example above, all I/O is checksummed on both sides.

