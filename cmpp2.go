package cmpp2

import (
	"bufio"
	//"fmt"
	//"errors"
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

func NewCmpp2Connect(host string, port int) (*cmpp2, error) {
	c := &cmpp2{}
	err := c.connect(host, port)
	return c, err
}

func (cmpp2 *cmpp2) connect(host string, port int) (err error) {
	// Create TCP connection
	cmpp2.conn, err = net.Dial("tcp", host+":"+strconv.Itoa(port))
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

func (s *cmpp2) Write(b []byte) error {
	_, err := s.conn.Write(b)

	return err
}

func (s *cmpp2) Read() (*CMPP_struct, error) {
	l := make([]byte, 4)
	cs := &CMPP_struct{}
	_, err := s.conn.Read(l)
	if err != nil {
		return cs, err
	}

	pduLength := unpackUi32(l) - 4
	if pduLength > 4096 {
		return cs, err
	}

	data := make([]byte, pduLength)

	i, err := s.conn.Read(data)
	if err != nil {
		return cs, err
	}

	if i != int(pduLength) {
		return cs, err
	}

	pkt := append(l, data...)
	cmpp, err := ParseMessage(pkt)
	if err != nil {
		return cs, err
	}
	return cmpp, err
	//pdu, err := ParsePdu(pkt)
	//if err != nil {
	//	return nil, err
	//}

	//return pdu, nil
}

func (s *cmpp2) NewSeqNum() uint32 {
	defer s.mu.Unlock()
	s.mu.Lock()
	s.sequence++
	return s.sequence
}

func (cmpp2 *cmpp2) bind(ms CMPP_struct) (err error) {

	h := new(Cmpp_header)
	h.Command_id = CMPP_CONNECT
	h.Sequence_Id = cmpp2.NewSeqNum()
	ms.setHeader(h)
	cmpp2.GetResp(ms)
	return
}

func (cmpp2 *cmpp2) GetResp(ms CMPP_struct) (err error) {
	h := new(Cmpp_header)
	err = h.read(cmpp2.reader)
	if err != nil {
		return
	}
	return
}
