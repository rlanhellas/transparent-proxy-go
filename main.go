package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"syscall"
	"unsafe"
)

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func main() {
	host := flag.String("host", ":8080", "host proxy server")

	flag.Parse()

	listener, err := net.Listen("tcp", *host)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(realServerAddress(&conn))

		//	go func(c net.Conn) {
		//		br := bufio.NewReader(c)
		//		req, err := http.ReadRequest(br)
		//		if err != nil {
		//			log.Println("buffer: ", err)
		//			return
		//		}
		//
		//		if req.Method == http.MethodConnect {
		//
		//			response := &http.Response{
		//				StatusCode: 200,
		//				ProtoMajor: 1,
		//				ProtoMinor: 1,
		//			}
		//			response.Write(c)
		//
		//			destConn, err := net.DialTimeout("tcp", req.URL.Host, 10*time.Second)
		//			if err != nil {
		//				response := &http.Response{
		//					StatusCode: http.StatusRequestTimeout,
		//					ProtoMajor: 1,
		//					ProtoMinor: 1,
		//				}
		//				response.Write(c)
		//				return
		//			}
		//
		//			go transfer(destConn, c)
		//			go transfer(c, destConn)
		//
		//		} else {
		//			response := &http.Response{
		//				StatusCode: http.StatusRequestTimeout,
		//				ProtoMajor: 1,
		//				ProtoMinor: 1,
		//				Body:       ioutil.NopCloser(strings.NewReader("hello world")),
		//			}
		//			response.Write(c)
		//			c.Close()
		//			return
		//		}
		//	}(conn)
	}
}

type sockaddr struct {
	family uint16
	data   [14]byte
}

const SO_ORIGINAL_DST = 80

// realServerAddress returns an intercepted connection's original destination.
func realServerAddress(conn *net.Conn) (string, error) {
	tcpConn, ok := (*conn).(*net.TCPConn)
	if !ok {
		return "", errors.New("not a TCPConn")
	}

	file, err := tcpConn.File()
	if err != nil {
		return "", err
	}

	// To avoid potential problems from making the socket non-blocking.
	tcpConn.Close()
	*conn, err = net.FileConn(file)
	if err != nil {
		return "", err
	}

	defer file.Close()
	fd := file.Fd()

	var addr sockaddr
	size := uint32(unsafe.Sizeof(addr))
	err = getsockopt(int(fd), syscall.SOL_IP, SO_ORIGINAL_DST, uintptr(unsafe.Pointer(&addr)), &size)
	if err != nil {
		return "", err
	}

	var ip net.IP
	switch addr.family {
	case syscall.AF_INET:
		ip = addr.data[2:6]
	default:
		return "", errors.New("unrecognized address family")
	}

	port := int(addr.data[0])<<8 + int(addr.data[1])

	return net.JoinHostPort(ip.String(), strconv.Itoa(port)), nil
}

func getsockopt(s int, level int, name int, val uintptr, vallen *uint32) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_GETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(unsafe.Pointer(vallen)), 0)
	if e1 != 0 {
		err = e1
	}

	return
}
