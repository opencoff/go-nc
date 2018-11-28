package main

// Go implementation of NC
//
// (c) 2013 Sudhi Herle <sw-at-herle.net>
// License GPLv2
//
// nc on steroids:
//   - checksum on network I/O
//   - traffic stats

import (
    "net"
    "os"
    "fmt"
    "io"
    "flag"
    "encoding/hex"
    "crypto/md5"
)

// Go command line option parsing truly sucks - especially in
// comparison to getopt_long(3) or python optparse.


// These are command line flags we process
var f_cksum  = flag.Bool("checksum", false, "Compute checksum on all network I/O")
var f_stats  = flag.Bool("statistics", false, "Show traffic statistics")
var f_listen = flag.Bool("listen", false, "Work in listen mode")
var f_bidir  = flag.Bool("bidirectional", false, "Do I/O in both directions")
var f_verbos = flag.Bool("verbose", false, "Show verbose progress messages")
var f_hex    = flag.Bool("hexdump", false, "Show hexdump of traffic")



// I/O hunk size
var bufsiz  int = 4 * 1048576

// object on which we will do I/O
type io_obj struct {
    addr string
}

// If we have multiple goroutines, they will return the result of
// their computation as this struct.
// If err is not nil, then str will have some meaningful status to
// report
type retval struct {
    err error
    str string
}


// Print error message and die
func die(format string, v ...interface{}) {
    s := fmt.Sprintf(format, v...)
    n := len(s)
    if n > 0 && s[n-1] != '\n' {
        s += "\n"
    }

    os.Stderr.WriteString(s)
    os.Exit(1)
}

func verbose(format string, v ...interface{}) {
    if *f_verbos {
        fmt.Fprintf(os.Stderr, format, v...)
    }
}


// Scaffolding for hexdump
// To do a dummy hexdump, we have to implement the io.WriteCloser
// interface.

type dummy_hexdump struct {
    // Nothing here
}


func (d *dummy_hexdump) Write(b []byte) (n int, err error) {
    return 0, nil
}

func (d* dummy_hexdump) Close() (error) {
    return nil
}


// Read from 'rd' and write to 'wr'. For informational purposes,
// 'dir' is a text string indicating "from" or "to". The result of
// the computation is in the chan of retval
func (o *io_obj) counting_io(rd io.Reader, wr io.Writer, dir string, ch chan retval) {
    var err error
    var n int64 = 0
    var m int = 0
    var d io.WriteCloser

    h  := md5.New()
    b  := make([]byte, bufsiz)

    if *f_hex {
        d = hex.Dumper(os.Stderr)
    } else {
        d = &dummy_hexdump{}
    }


    for {
        m, err = rd.Read(b)

        if m == 0 || err == io.EOF {
          break
        } else if err != nil {
            verbose("%s %s: ERROR %s\n", dir, o.addr, err)
            r := retval{err:err}
            ch <- r
            return
        }


        n += int64(m)
        h.Write(b[0:m])

        d.Write(b[0:m])

        m, err = wr.Write(b[0:m])
        if m != len(b) && err != nil {

            r := retval{err:err}
            ch <- r
            return
        }

    }

    // This closes out any pending bytes in the hexdump buffer
    d.Close()

    var msg string

    if *f_stats {
        msg += fmt.Sprintf("%d bytes %s %s", n, dir, o.addr)
    }
    if *f_cksum {
        if len(msg) > 0 {
            msg += " "
        }
        msg += fmt.Sprintf("(%x)", h.Sum(nil))
    }

    //verbose("I/O [%s] %s done: %d\n", dir, o.addr, n)
    ch <- retval{err:nil, str:msg}
}


func main() {

    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        die("Usage: %s [options] server:port\n", os.Args[0])
    }

    addr := args[0]

    n_ch := 1
    if *f_bidir {
        n_ch += 1
    }

    ch := make(chan retval, n_ch)

    var conn net.Conn
    var err error

    if *f_listen {

        verbose("Listening on %s...\n", addr)
        ln, err := net.Listen("tcp", addr)
        if err != nil {
            die("Can't listen on %s: %s\n", addr, err)
        }

        conn, err = ln.Accept()
        if err != nil {
            die("Can't accept on %s: %s\n", addr, err)
        }

        peer := conn.RemoteAddr().String()

        verbose("Accepted from %s\n", peer)

        o := io_obj{peer}
        go o.counting_io(conn, os.Stdout, "from", ch)
        if *f_bidir {
            go o.counting_io(os.Stdin, conn, "to", ch)
        }
    } else {
        verbose("Connecting to %s...\n", addr)
        conn, err = net.Dial("tcp", addr)
        if err != nil {
            die("Can't connect to %s: %s\n", addr, err)
        }

        o := io_obj{conn.RemoteAddr().String()}

        go o.counting_io(os.Stdin, conn, "to", ch)
        if *f_bidir {
            go o.counting_io(conn, os.Stdout, "from", ch)
        }
    }

    verbose("Waiting for %d goroutines to end ..\n", n_ch)
    errcode := 0
    for i:=0; i < n_ch; i++ {
        r := <-ch
        if r.err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", r.err)
            errcode++
        } else {
            fmt.Fprintf(os.Stderr, "%s\n", r.str)
        }
    }

    close(ch)

    conn.Close()
    os.Exit(errcode)
}

// EOF
