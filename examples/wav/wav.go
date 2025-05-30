package wav

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

// WAVWriter handles writing WAV files
type WAVWriter struct {
	file       *os.File
	sampleRate int
	channels   int
}

// NewWAVWriter creates a new WAV file writer
func NewWAVWriter(filename string, sampleRate, channels int) (*WAVWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}

	writer := &WAVWriter{
		file:       file,
		sampleRate: sampleRate,
		channels:   channels,
	}

	// Write WAV header
	if err := writer.writeHeader(); err != nil {
		file.Close()
		return nil, err
	}

	return writer, nil
}

// writeHeader writes the WAV file header
func (w *WAVWriter) writeHeader() error {
	// RIFF header
	if _, err := w.file.Write([]byte("RIFF")); err != nil {
		return err
	}

	// File size (to be filled later)
	if err := binary.Write(w.file, binary.LittleEndian, uint32(0)); err != nil {
		return err
	}

	// WAVE format
	if _, err := w.file.Write([]byte("WAVE")); err != nil {
		return err
	}

	// fmt chunk
	if _, err := w.file.Write([]byte("fmt ")); err != nil {
		return err
	}

	// fmt chunk size
	if err := binary.Write(w.file, binary.LittleEndian, uint32(16)); err != nil {
		return err
	}

	// Audio format (1 for PCM)
	if err := binary.Write(w.file, binary.LittleEndian, uint16(1)); err != nil {
		return err
	}

	// Number of channels
	if err := binary.Write(w.file, binary.LittleEndian, uint16(w.channels)); err != nil {
		return err
	}

	// Sample rate
	if err := binary.Write(w.file, binary.LittleEndian, uint32(w.sampleRate)); err != nil {
		return err
	}

	// Byte rate
	byteRate := w.sampleRate * w.channels * 2 // 2 bytes per sample
	if err := binary.Write(w.file, binary.LittleEndian, uint32(byteRate)); err != nil {
		return err
	}

	// Block align
	blockAlign := w.channels * 2 // 2 bytes per sample
	if err := binary.Write(w.file, binary.LittleEndian, uint16(blockAlign)); err != nil {
		return err
	}

	// Bits per sample
	if err := binary.Write(w.file, binary.LittleEndian, uint16(16)); err != nil {
		return err
	}

	// data chunk
	if _, err := w.file.Write([]byte("data")); err != nil {
		return err
	}

	// data chunk size (to be filled later)
	if err := binary.Write(w.file, binary.LittleEndian, uint32(0)); err != nil {
		return err
	}

	return nil
}

// WriteSamples writes audio samples to the WAV file
func (w *WAVWriter) WriteSamples(samples []int16) error {
	// Write samples
	if err := binary.Write(w.file, binary.LittleEndian, samples); err != nil {
		return fmt.Errorf("写入音频数据失败: %w", err)
	}

	// Update file size in header
	fileSize := uint32(36 + len(samples)*2) // 36 bytes header + samples
	if _, err := w.file.Seek(4, io.SeekStart); err != nil {
		return err
	}
	if err := binary.Write(w.file, binary.LittleEndian, fileSize); err != nil {
		return err
	}

	// Update data chunk size
	if _, err := w.file.Seek(40, io.SeekStart); err != nil {
		return err
	}
	if err := binary.Write(w.file, binary.LittleEndian, uint32(len(samples)*2)); err != nil {
		return err
	}

	return nil
}

// Close closes the WAV file
func (w *WAVWriter) Close() error {
	return w.file.Close()
}

// ValidateWAVFile validates a WAV file
func ValidateWAVFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// Read and validate RIFF header
	header := make([]byte, 12)
	if _, err := file.Read(header); err != nil {
		return fmt.Errorf("读取文件头失败: %w", err)
	}

	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return fmt.Errorf("无效的 WAV 文件格式")
	}

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}

	// Validate file size
	if fileInfo.Size() < 44 { // Minimum WAV file size
		return fmt.Errorf("文件大小异常")
	}

	fmt.Printf("✓ WAV 文件验证通过\n")
	fmt.Printf("  - 文件大小: %d 字节\n", fileInfo.Size())
	fmt.Printf("  - 修改时间: %s\n", fileInfo.ModTime().Format(time.RFC3339))

	return nil
}
