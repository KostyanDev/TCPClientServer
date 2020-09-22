package main

import (
	"flag"
	//	"fmt"
	"io"
	"log"
	"net"
	//"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"syscall"
	"time"
)

type Server struct {
	proto   string
	addr    string
	handler func(c *net.TCPConn) error
}

type Client struct {
	proto       string
	addr        string
	handler     func(c *net.TCPConn) error
	concurrency int
	size        int
	nflight     int
	reqres      bool
	saddr       string
}

func (s *Server) ListenAndGo() error {
	tcpaddr, _ := net.ResolveTCPAddr(s.proto, s.addr)
	ln, err := net.ListenTCP(s.proto, tcpaddr)
	if err != nil {
		log.Println("Failed to listen for tcp connections on address ", s.addr, " Error: ", err)
		return err
	}

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Println("Failed to accept connection ", conn, " due to error ", err)
			continue
		}
		log.Println("Client ", conn.RemoteAddr(), " connected")
		go s.handler(conn)
	}
	return nil
}

func (c *Client) NewConnection() (*net.TCPConn, error) {
	srcTcpAddr, _ := net.ResolveTCPAddr(c.proto, c.saddr)
	dstTcpAddr, _ := net.ResolveTCPAddr(c.proto, c.addr)
	return net.DialTCP(c.proto, srcTcpAddr, dstTcpAddr)
}

func (c *Client) ConnectAndGo() error {
	conns := make([]*net.TCPConn, c.concurrency)
	for i := 0; i < c.concurrency; i++ {
		conn, err := c.NewConnection()
		if err != nil {
			log.Println("Failed to connect to tcp server on address ", c.addr, " from source address: ", c.saddr, " Error: ", err)
			conn.Close()
			return err
		}
		log.Println("Client ", conn.RemoteAddr(), " connected")
		conns[i] = conn
	}
	for i := 0; i < c.concurrency; i++ {
		go c.handler(conns[i])
	}
	return nil
}

func TCPConnRead(c *net.TCPConn) error {
	for {
		_, err := c.Read(b)
		if err != nil {
			if err == io.EOF {
				log.Println("Client ", c.RemoteAddr(), " disconnected")
				c.Close()
				return nil
			} else {
				log.Println("Failed reading bytes from conn: ", c, " with error ", err)
				c.Close()
				return err
			}
		} else {
			//log.Println("Read ", n, " bytes from client ", c.RemoteAddr())
		}
	}
	return nil
}

var b []byte

func TCPConnWrite(c *net.TCPConn) error {
	c.SetWriteBuffer(*packetsize)

	for {
		_, err := c.Write(b)
		if err != nil {
			if err == io.EOF {
				log.Println("Client ", c.RemoteAddr(), " disconnected")
				c.Close()
				return nil
			} else {
				log.Println("Failed writing bytes to conn: ", c, " with error ", err)
				c.Close()
				return err
			}
		} else {
			//log.Println("Wrote ", n, " bytes to connection ", c.LocalAddr())
		}
	}
	log.Println("Came out of write loop! ", c)
	return nil
}

var packetsize *int

func main() {
	host := flag.String("host", "127.0.0.1", "Host IP address")
	port := flag.String("port", "8082", "Port")
	shost := flag.String("shost", "127.0.0.1", "Host IP address")
	sport := flag.String("sport", "8082", "Port")
	listen := flag.Bool("listen", false, "Listen")
	packetsize = flag.Int("size", 1500, "Size of packets to send")
	nconn := flag.Int("nconn", 254, "Number of concurrent connections")
	reqres := flag.Bool("reqres", false, "Request/Response protocol")
	nflight := flag.Int("nflight", 1024, "Number of requests in flight before waiting for response")
	profile := flag.String("profile", "", "write profile to file with following prefix")
	flag.Parse()

	if *profile != "/home/kovcharenko/Downloads/1" {
		go doprofile(*profile)
	}
	if flag.NArg() != 0 {
		log.Println("Usage:")
		flag.PrintDefaults()
		return
	}
	b = make([]byte, *packetsize)
	go GoRuntimeStats()
	if *listen {
		s := &Server{proto: "tcp", addr: net.JoinHostPort(*host, *port), handler: TCPConnWrite}
		s.ListenAndGo()
	} else {
		time.Sleep(20 * time.Second)
		c := &Client{proto: "tcp", addr: net.JoinHostPort(*host, *port), handler: TCPConnRead, size: *packetsize, concurrency: *nconn, nflight: *nflight, reqres: *reqres, saddr: net.JoinHostPort(*shost, *sport)}
		c.ConnectAndGo()
		SigIntHandler()
	}

	log.Println("Finished execution!")
}

func SigIntHandler() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	log.Println("CTRL-C; exiting")
	os.Exit(0)
}

func GoRuntimeStats() {
	m := &runtime.MemStats{}
	for {
		time.Sleep(2 * time.Second)
		log.Println("# goroutines: ", runtime.NumGoroutine())
		runtime.ReadMemStats(m)
		log.Println("# malloc       : ", m.Mallocs)
		log.Println("# free         : ", m.Frees)
		log.Println("GC enabled     : ", m.EnableGC)
		log.Println("# GC           : ", m.NumGC)
		log.Println("Last GC time   : ", m.LastGC)
		//runtime.GC()
	}
}
func doprofile(fn string) {
	for i := 1; i > 0; i++ {
		fc, err := os.Create(fn + "-cpu-" + strconv.Itoa(i) + ".prof")
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(fc)
		time.Sleep(300 * time.Second)
		pprof.StopCPUProfile()
		fc.Close()

		fh, err := os.Create(fn + "-heap-" + strconv.Itoa(i) + ".prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(fh)
		fh.Close()

		ft, err := os.Create(fn + "-threadcreate-" + strconv.Itoa(i) + ".prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.Lookup("threadcreate").WriteTo(ft, 0)
		ft.Close()
		log.Println("Created CPU, heap and threadcreate profile of 300 seconds")
	}
}
