# README for go-nc


## What is this?
This is a simple clone of `netcat` in golang. I wrote it many moons
ago as my first go program. It is slightly different from the
traditional `netcat` in the following ways:

* adds md5 checksum on network I/O
* traffic stats
* hexdump of the traffic
* Bi directional I/O is optional; server always writes to `stdout`
  and clinet always reads from `stdin`. This makes it easy to
  remember.

## How do I build it?
You will need Go 1.5 or later:
 
    git clone https://github.com/opencoff/go-nc
    cd go-nc
    make

## How do I run it?
`gonc` comes with helpful commandline options. Try:

    ./bin/gonc -help



