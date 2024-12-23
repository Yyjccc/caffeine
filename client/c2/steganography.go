package c2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/png"
)

// Steganography 图片隐写处理器
type Steganography struct {
	MaxPayloadSize int // 最大载荷大小
}

// NewSteganography 创建新的隐写处理器
func NewSteganography() *Steganography {
	return &Steganography{
		MaxPayloadSize: 1024 * 1024, // 默认最大1MB
	}
}

// EmbedData 将数据隐写到图片中
func (s *Steganography) EmbedData(imgData []byte, payload []byte) ([]byte, error) {
	// 检查payload大小
	if len(payload) > s.MaxPayloadSize {
		return nil, errors.New("payload size exceeds maximum allowed size")
	}

	// 解码原始图片
	img, err := png.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	// 获取图片尺寸
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 检查图片容量是否足够
	maxBytes := (width * height * 3) / 8 // 每个像素3个通道，每个通道可以存储1位
	if len(payload) > maxBytes-4 {       // 减去4字节用于存储长度
		return nil, errors.New("image capacity is not enough for payload")
	}

	// 创建新的RGBA图片
	newImg := image.NewRGBA(bounds)

	// 首先写入payload长度（4字节）
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(payload)))
	payloadWithLength := append(lengthBytes, payload...)

	// 开始嵌入数据
	bitIndex := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			r, g, b, a := pixel.RGBA()

			// 对每个颜色通道的最低位进行修改
			if bitIndex/8 < len(payloadWithLength) {
				byteIndex := bitIndex / 8
				bitPos := 7 - (bitIndex % 8)
				bit := (payloadWithLength[byteIndex] >> bitPos) & 1

				// 修改红色通道
				if bitIndex < len(payloadWithLength)*8 {
					r = (r & 0xFFFE) | uint32(bit)
					bitIndex++
				}

				// 修改绿色通道
				if bitIndex < len(payloadWithLength)*8 {
					g = (g & 0xFFFE) | uint32(bit)
					bitIndex++
				}

				// 修改蓝色通道
				if bitIndex < len(payloadWithLength)*8 {
					b = (b & 0xFFFE) | uint32(bit)
					bitIndex++
				}
			}

			newImg.Set(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}

	// 将新图片编码为PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, newImg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExtractData 从图片中提取隐写数据
func (s *Steganography) ExtractData(imgData []byte) ([]byte, error) {
	// 解码图片
	img, err := png.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 首先读取长度（前4字节）
	lengthBits := make([]byte, 4)
	bitIndex := 0

	// 读取长度信息
	for y := 0; y < height && bitIndex < 32; y++ {
		for x := 0; x < width && bitIndex < 32; x++ {
			pixel := img.At(x, y)
			r, g, b, _ := pixel.RGBA()

			// 从每个通道的最低位提取数据
			if bitIndex < 32 {
				byteIndex := bitIndex / 8
				bitPos := 7 - (bitIndex % 8)
				lengthBits[byteIndex] |= ((byte(r) & 1) << bitPos)
				bitIndex++
			}
			if bitIndex < 32 {
				byteIndex := bitIndex / 8
				bitPos := 7 - (bitIndex % 8)
				lengthBits[byteIndex] |= ((byte(g) & 1) << bitPos)
				bitIndex++
			}
			if bitIndex < 32 {
				byteIndex := bitIndex / 8
				bitPos := 7 - (bitIndex % 8)
				lengthBits[byteIndex] |= ((byte(b) & 1) << bitPos)
				bitIndex++
			}
		}
	}

	// 解析payload长度
	payloadLength := binary.BigEndian.Uint32(lengthBits)
	if payloadLength > uint32(s.MaxPayloadSize) {
		return nil, errors.New("extracted payload size exceeds maximum allowed size")
	}

	// 提取payload数据
	payload := make([]byte, payloadLength)
	bitIndex = 32 // 从长度信息之后开始

	for y := 0; y < height && bitIndex < int(payloadLength)*8+32; y++ {
		for x := 0; x < width && bitIndex < int(payloadLength)*8+32; x++ {
			pixel := img.At(x, y)
			r, g, b, _ := pixel.RGBA()

			// 从每个通道的最低位提取数据
			if bitIndex < int(payloadLength)*8+32 {
				byteIndex := (bitIndex - 32) / 8
				bitPos := 7 - ((bitIndex - 32) % 8)
				payload[byteIndex] |= ((byte(r) & 1) << bitPos)
				bitIndex++
			}
			if bitIndex < int(payloadLength)*8+32 {
				byteIndex := (bitIndex - 32) / 8
				bitPos := 7 - ((bitIndex - 32) % 8)
				payload[byteIndex] |= ((byte(g) & 1) << bitPos)
				bitIndex++
			}
			if bitIndex < int(payloadLength)*8+32 {
				byteIndex := (bitIndex - 32) / 8
				bitPos := 7 - ((bitIndex - 32) % 8)
				payload[byteIndex] |= ((byte(b) & 1) << bitPos)
				bitIndex++
			}
		}
	}

	return payload, nil
}

// ValidateImage 验证图片是否适合用于隐写
func (s *Steganography) ValidateImage(imgData []byte) error {
	// 检查是否为PNG格式
	if !bytes.HasPrefix(imgData, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return errors.New("image must be PNG format")
	}

	// 解码图片以验证完整性
	img, err := png.Decode(bytes.NewReader(imgData))
	if err != nil {
		return err
	}

	// 检查图片尺寸
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 检查最小尺寸要求
	if width*height < 100 { // 假设最小100像素
		return errors.New("image is too small for steganography")
	}

	return nil
}
