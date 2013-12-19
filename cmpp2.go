package cmpp2

import (
	"bufio"
	//"fmt"
	"errors"
	"net"
	//"os"
	"strconv"
	"sync"
)

type cmpp2 struct {
	conn      net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	connected bool
	bound     bool
	mu        sync.Mutex
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

func (s *cmpp2) NewSeqNum() uint32 {
	defer s.mu.Unlock()
	s.mu.Lock()
	s.sequence++
	return s.sequence
}

func (cmpp2 *cmpp2) bind(ms MSG_struct) (err error) {

	h := new(Cmpp_header)
	h.Command_id = CMPP_CONNECT
	h.Sequence_Id = cmpp2.NewSeqNum()
	ms.setHeader(h)
	cmpp2.GetResp(ms)
	return
}

func (cmpp2 *cmpp2) GetResp(ms MSG_struct) (err error) {
	h := new(Cmpp_header)
	err = h.read(cmpp2.reader)
	if err != nil {
		return
	}

}
