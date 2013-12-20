/*
消息头格式（Message Header）
*/

package cmpp2

import (
	"bufio"
	//"errors"
)

type Cmpp_header struct {
	Total_Length uint32     //消息总长度(含消息头及消息体)
	Command_id   COMMAND_ID //命令或响应类型
	Sequence_Id  uint32     //消息流水号,顺序累加,步长为1,循环使用（一对请求和应答消息的流水号必须相同）
}

func (h *Cmpp_header) write(w *bufio.Writer) (err error) {
	p := make([]byte, 12)
	copy(p[0:4], packUi32(uint32(h.Total_Length)))
	copy(p[4:8], packUi32(uint32(h.Command_id)))
	copy(p[8:12], packUi32(uint32(h.Sequence_Id)))
	_, err = w.Write(p)
	if err != nil {
		return
	}
	err = w.Flush()
	return

}

func (h *Cmpp_header) read(r *bufio.Reader) (err error) {
	p := make([]byte, 12)
	_, err = r.Read(p)
	if err != nil {
		return
	}
	h.Total_Length = unpackUi32(p[0:4])
	h.Command_id = COMMAND_ID(unpackUi32(p[4:8]))
	h.Sequence_Id = unpackUi32(p[8:12])
	return
}
