package cmpp2

import (
	"bufio"
	//"fmt"
	"errors"
	"net"
	//"os"
	"strconv"
)

type cmpp2 struct {
	conn      net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	connected bool
	bound     bool
	async     bool
	sequence  uint32
}

func (cmpp2 *cmpp2) connect(host string, port int) (err error) {
	// Create TCP connection
	cmpp2.conn, err = net.Dial("tcp", "", host+":"+strconv.Itoa(port))
	if err != nil {
		return
	}
	cmpp2.connected = true
	// Setup buffered reader/writer
	cmpp2.reader = bufio.NewReader(cmpp2.conn)
	cmpp2.writer = bufio.NewWriter(cmpp2.conn)
	return
}

func (cmpp2 *cmpp2) close() (err error) {
	err = cmpp2.conn.Close()
	cmpp2.connected = false
	return
}

func (cmpp2 *cmpp2) Async(async bool) {
	cmpp2.async = async
}

func (cmpp2 *cmpp2) bind() (err error) {
	cmpp2.sequence++
	h := new(Cmpp_header)
	h.Command_id = CMPP_CONNECT
	return
}

func (cmpp2 *cmpp2) GetResp() {

}
