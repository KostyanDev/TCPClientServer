package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

type Server struct {
	proto   string
	addr    string
	handler func(c *net.TCPConn) error
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
	port := flag.String("port", "8087", "Port")
	packetsize = flag.Int("size", 1500, "Size of packets to send")
	profile := flag.String("profile", "", "write profile to file with following prefix")

	if *profile != "" {
		go doprofile(*profile)
	}
	if flag.NArg() != 0 {
		log.Println("Usage:")
		flag.PrintDefaults()
		return
	}
	b = make([]byte, *packetsize)
	s := &Server{proto: "tcp", addr: net.JoinHostPort(*host, *port), handler: TCPConnWrite}
	s.ListenAndGo()
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
