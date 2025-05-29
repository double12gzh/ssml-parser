package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// WAVWriter WAV 文件写入器
type WAVWriter struct {
	file       *os.File
	sampleRate int
	channels   int
	dataSize   int
}

// NewWAVWriter 创建新的 WAV 写入器
func NewWAVWriter(filename string, sampleRate, channels int) (*WAVWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	writer := &WAVWriter{
		file:       file,
		sampleRate: sampleRate,
		channels:   channels,
		dataSize:   0,
	}

	// 写入 WAV 头部（先写入占位符，稍后更新）
	err = writer.writeHeader()
	if err != nil {
		file.Close()
		return nil, err
	}

	return writer, nil
}

// WriteSamples 写入音频样本数据
func (w *WAVWriter) WriteSamples(samples []int16) error {
	for _, sample := range samples {
		err := binary.Write(w.file, binary.LittleEndian, sample)
		if err != nil {
			return err
		}
		w.dataSize += 2 // 每个样本 2 字节
	}
	return nil
}

// Close 关闭文件并更新头部信息
func (w *WAVWriter) Close() error {
	// 更新文件头部的大小信息
	_, err := w.file.Seek(4, 0) // 跳到文件大小位置
	if err != nil {
		return err
	}

	fileSize := 36 + w.dataSize
	err = binary.Write(w.file, binary.LittleEndian, uint32(fileSize))
	if err != nil {
		return err
	}

	// 更新数据块大小
	_, err = w.file.Seek(40, 0) // 跳到数据大小位置
	if err != nil {
		return err
	}

	err = binary.Write(w.file, binary.LittleEndian, uint32(w.dataSize))
	if err != nil {
		return err
	}

	return w.file.Close()
}

// writeHeader 写入 WAV 文件头部
func (w *WAVWriter) writeHeader() error {
	// RIFF 标识符
	_, err := w.file.WriteString("RIFF")
	if err != nil {
		return err
	}

	// 文件大小（稍后更新）
	err = binary.Write(w.file, binary.LittleEndian, uint32(0))
	if err != nil {
		return err
	}

	// WAVE 标识符
	_, err = w.file.WriteString("WAVE")
	if err != nil {
		return err
	}

	// fmt 子块
	_, err = w.file.WriteString("fmt ")
	if err != nil {
		return err
	}

	// fmt 子块大小
	err = binary.Write(w.file, binary.LittleEndian, uint32(16))
	if err != nil {
		return err
	}

	// 音频格式（PCM = 1）
	err = binary.Write(w.file, binary.LittleEndian, uint16(1))
	if err != nil {
		return err
	}

	// 声道数
	err = binary.Write(w.file, binary.LittleEndian, uint16(w.channels))
	if err != nil {
		return err
	}

	// 采样率
	err = binary.Write(w.file, binary.LittleEndian, uint32(w.sampleRate))
	if err != nil {
		return err
	}

	// 字节率 (采样率 * 声道数 * 每样本位数 / 8)
	byteRate := w.sampleRate * w.channels * 16 / 8
	err = binary.Write(w.file, binary.LittleEndian, uint32(byteRate))
	if err != nil {
		return err
	}

	// 块对齐 (声道数 * 每样本位数 / 8)
	blockAlign := w.channels * 16 / 8
	err = binary.Write(w.file, binary.LittleEndian, uint16(blockAlign))
	if err != nil {
		return err
	}

	// 每样本位数
	err = binary.Write(w.file, binary.LittleEndian, uint16(16))
	if err != nil {
		return err
	}

	// data 子块标识符
	_, err = w.file.WriteString("data")
	if err != nil {
		return err
	}

	// data 子块大小（稍后更新）
	err = binary.Write(w.file, binary.LittleEndian, uint32(0))
	if err != nil {
		return err
	}

	return nil
}

// 辅助函数：验证生成的 WAV 文件
func ValidateWAVFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 读取并验证 RIFF 头部
	riff := make([]byte, 4)
	_, err = file.Read(riff)
	if err != nil {
		return fmt.Errorf("读取 RIFF 头部失败: %w", err)
	}

	if string(riff) != "RIFF" {
		return fmt.Errorf("无效的 RIFF 头部: %s", string(riff))
	}

	// 跳过文件大小
	file.Seek(4, 1)

	// 读取并验证 WAVE 标识符
	wave := make([]byte, 4)
	_, err = file.Read(wave)
	if err != nil {
		return fmt.Errorf("读取 WAVE 标识符失败: %w", err)
	}

	if string(wave) != "WAVE" {
		return fmt.Errorf("无效的 WAVE 标识符: %s", string(wave))
	}

	fmt.Printf("✓ WAV 文件格式验证通过: %s\n", filename)
	return nil
}
