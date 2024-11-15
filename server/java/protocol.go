package java

import (
	"bytes"
	"encoding/binary"
)

// 封装通信协议
type Protocol struct {
	Version int
	Load    bool
	NameLen int32
	Name    string
	Data    []byte
}

func NewProtocol(name string, data []byte, load bool) *Protocol {
	return &Protocol{
		Version: 1,
		Load:    load,
		Name:    name,
		Data:    data,
		NameLen: int32(len(name)),
	}
}

// Encode 将 Protocol 对象编码为二进制
func (p *Protocol) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)

	// 写入 Version
	if err := binary.Write(buffer, binary.BigEndian, p.Version); err != nil {
		return nil, err
	}

	// 写入 Load
	load := byte(0)
	if p.Load {
		load = 1
	}
	if err := binary.Write(buffer, binary.BigEndian, load); err != nil {
		return nil, err
	}

	// 写入 NameLen
	p.NameLen = int32(len(p.Name))
	if err := binary.Write(buffer, binary.BigEndian, p.NameLen); err != nil {
		return nil, err
	}

	// 写入 Name
	if err := binary.Write(buffer, binary.BigEndian, []byte(p.Name)); err != nil {
		return nil, err
	}

	// 写入 Data
	if err := binary.Write(buffer, binary.BigEndian, p.Data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Decode 从二进制数据解码到 Protocol 对象
func (p *Protocol) Decode(data []byte) error {
	buffer := bytes.NewReader(data)

	// 读取 Version
	if err := binary.Read(buffer, binary.BigEndian, &p.Version); err != nil {
		return err
	}

	// 读取 Load
	var load byte
	if err := binary.Read(buffer, binary.BigEndian, &load); err != nil {
		return err
	}
	p.Load = load != 0

	// 读取 NameLen
	if err := binary.Read(buffer, binary.BigEndian, &p.NameLen); err != nil {
		return err
	}

	// 读取 Name
	nameBytes := make([]byte, p.NameLen)
	if err := binary.Read(buffer, binary.BigEndian, &nameBytes); err != nil {
		return err
	}
	p.Name = string(nameBytes)

	// 读取 Data
	p.Data = make([]byte, buffer.Len())
	if err := binary.Read(buffer, binary.BigEndian, &p.Data); err != nil {
		return err
	}

	return nil
}
