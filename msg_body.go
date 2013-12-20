/*消息定义*/
package cmpp2

import (
	"bufio"
	"encoding/binary"
	//"errors"
	//"os"
	//"reflect"
)

func unpackUi32(b []byte) (n uint32) {
	n = binary.BigEndian.Uint32(b)
	return
}

func packUi32(n uint32) (b []byte) {
	b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return
}

func unpackUi16(b []byte) (n uint16) {
	n = binary.BigEndian.Uint16(b)
	return
}

func packUi16(n uint16) (b []byte) {
	b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return
}

func packUi8(n uint8) (b []byte) {
	b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(n))
	return b[1:]
}

func ParseMessage(b []byte) (*CMPP_struct, error) {
	var err error
	cmpp := &CMPP_struct{}
	if len(b) < 12 {
		return cmpp, err
	}

	header := &Cmpp_header{}
	header.Total_Length = unpackUi32(b[0:4])
	header.Command_id = COMMAND_ID(unpackUi32(b[4:8]))
	header.Sequence_Id = unpackUi32(b[8:12])
	cmpp.Header = header
	switch header.Command_id {
	case CMPP_CONNECT_RESP:
		conresp, errs := NewCmpp_connect_resp(b[12:])
		if err != nil {
			return cmpp, errs
		}
		cmpp.Body = MSG(conresp)
	}
	return cmpp, nil

}

type MSG interface {
	read(r *bufio.Reader) (err error)
	write(w *bufio.Writer) (err error)
	SetSeqNum(uint32)
	setHeader(h *Cmpp_header)
	GetHeader() *Cmpp_header
	GetStruct() interface{}
}

type CMPP_struct struct {
	Header *Cmpp_header
	Body   MSG
}

func (m *CMPP_struct) setHeader(h *Cmpp_header) {
	m.Header = h
}

func (m *CMPP_struct) SetSeqNum(i uint32) {
	m.Header.Sequence_Id = i
}

func (m *CMPP_struct) GetHeader() *Cmpp_header {
	return m.Header
}

func (m *CMPP_struct) GetStruct() interface{} {
	return *m
}

//CMPP_CONNECT消息定义
type Cmpp_connect struct {
	Source_Add          string //此处为SP_Id，即SP的企业代码
	AuthenticatorSource string //用于鉴别源地址 MD5（Source_Addr+9 字节的0 +shared secret+timestamp）
	Version             uint32 //双方协商的版本号
	Timestamp           uint32 //MMDDHHMMSS
}

func (cc *Cmpp_connect) write(w *bufio.Writer) (err error) {

	return
}

func (cc *Cmpp_connect) read(r *bufio.Reader) (err error) {
	return
}

//CMPP_CONNECT响应消息定义
/*状态
0：正确
1：消息结构错
 2：非法源地址
 3：认证错
 4：版本太高
 5~ ：其他错误*/
type Cmpp_connect_resp struct {
	Status            uint32
	AuthenticatorISMG string //ISMG认证码 认证出错时，此项为空 MD5（Status+AuthenticatorSource+shared secret）
	Version           uint32 //服务器支持的最高版本号
}

func NewCmpp_connect_resp(b []byte) (*Cmpp_connect_resp, error) {
	ccr := &Cmpp_connect_resp{}
	ccr.Status = unpackUi32(b[0:1])
	ccr.AuthenticatorISMG = string(b[1:16])
	ccr.Version = unpackUi32(b[16:17])
	return ccr, nil
}

func (cc *Cmpp_connect_resp) write(w *bufio.Writer) (err error) {

	return
}

func (cc *Cmpp_connect_resp) read(r *bufio.Reader) (err error) {

	return
}

func (cc *Cmpp_connect_resp) GetHeader() *Cmpp_header {
	return &Cmpp_header{}
}

func (cc *Cmpp_connect_resp) GetStruct() interface{} {
	return *cc
}

func (cc *Cmpp_connect_resp) setHeader(h *Cmpp_header) {
}

func (cc *Cmpp_connect_resp) SetSeqNum(i uint32) {
	cc.GetHeader().Sequence_Id = i
}

/*SP或ISMG请求拆除连接（CMPP­_TERMINATE）操作*/
//CMPP­_TERMINATE消息定义 无消息体
type Cmpp_terminate interface{}

//CMPP­_TERMINATE响应消息定义 无消息体
type Cmpp_terminate_resp interface{}

/*ISMG或SP以CMPP_TERMINATE_RESP消息响应请求*/

//CMPP­_SUBMIT消息定义
type Cmpp_submit struct {
	Msg_Id              uint32 //信息标识，由SP侧短信网关本身产生，本处填空。
	Pk_total            uint32 //相同Msg_Id的信息总条数，从1开始
	Pk_number           uint32 //相同Msg_Id的信息序号，从1开始
	Registered_Delivery uint32 //是否要求返回状态确认报告：0：不需要1：需要2：产生SMC话单
	Msg_level           uint32 //信息级别
	Service_Id          string //业务类型，是数字、字母和符号的组合。
	Fee_UserType        uint32 //计费用户类型字段 0：对目的终端MSISDN计费；1：对源终端MSISDN计费；2：对SP计费;3：表示本字段无效，对谁计费参见Fee_terminal_Id字段。
	Fee_terminal_Id     uint32 //被计费用户的号码（如本字节填空，则表示本字段无效，对谁计费参见Fee_UserType字段，本字段与Fee_UserType字段互斥）
	TP_pId              uint32 //GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9
	TP_udhi             uint32 //GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.23,仅使用1位，右对齐
	Msg_Fmt             uint32 //信息格式  0：ASCII串  3：短信写卡操作  4：二进制信息  8：UCS2编码15：含GB汉字  。。。。。。
	Msg_src             string //信息内容来源(SP_Id)
	FeeType             string //资费类别 01：对“计费用户号码”免费 02：对“计费用户号码”按条计信息费03：对“计费用户号码”按包月收取信息费04：对“计费用户号码”的信息费封顶05：对“计费用户号码”的收费是由SP实现
	FeeCode             string //资费代码（以分为单位）
	ValId_Time          string //存活有效期，格式遵循SMPP3.3协议
	At_Time             string //定时发送时间，格式遵循SMPP3.3协议
	Src_Id              string //源号码SP的服务代码或前缀为服务代码的长号码, 网关将该号码完整的填到SMPP协议Submit_SM消息相应的source_addr字段，该号码最终在用户手机上显示为短消息的主叫号码
	DestUsr_tl          uint32 //接收信息的用户数量(小于100个用户)
	Dest_terminal_Id    string //接收短信的MSISDN号码
	Msg_Length          uint32 //信息长度(Msg_Fmt值为0时：<160个字节；其它<=140个字节)
	Msg_Content         string //信息内容
	Reserve             string //保留

}

//CMPP­_SUBMIT_RESP响应消息定义
/*
信息标识，生成算法如下：
采用64位（8字节）的整数：
（1） 时间（格式为MMDDHHMMSS，即月日时分秒）：bit64~bit39，其中
bit64~bit61：月份的二进制表示；
bit60~bit56：日的二进制表示；
bit55~bit51：小时的二进制表示；
bit50~bit45：分的二进制表示；
bit44~bit39：秒的二进制表示；
（2） 短信网关代码：bit38~bit17，把短信网关的代码转换为整数填写到该字段中。
（3） 序列号：bit16~bit1，顺序增加，步长为1，循环使用。
各部分如不能填满，左补零，右对齐。

（SP根据请求和应答消息的Sequence_Id一致性就可得到CMPP_Submit消息的Msg_Id）
*/
type Cmpp_submit_resp struct {
	Msg_Id uint32
	Result uint32 //结果0：正确1：消息结构错
}

//CMPP_QUERY消息的定义
type Cmpp_query struct {
}

//CMPP_QUERY_RESP消息的定义

//CMPP_DELIVER消息定义
//CMPP_DELIVER_RESP消息定义

//CMPP_CANCEL消息定义

//CMPP_CANCEL_RESP消息定义

//CMPP_ACTIVE_TEST定义 无消息体

//CMPP_ACTIVE_TEST_RESP定义
