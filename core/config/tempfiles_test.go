package pkg

import (
	"io/ioutil"
	"os"
	"testing"
)

// 测试TempFileWithText函数
func TestTempFileWithText(t *testing.T) {
	text := "Hello, World!"
	tmpFile, err := TempFileWithText(text)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // 清理临时文件
	defer tmpFile.Close()           // 关闭文件句柄

	// 读取文件内容并验证
	content, err := ioutil.ReadAll(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read from temp file: %v", err)
	}

	if string(content) != text {
		t.Errorf("Content mismatch: expected %v, got %v", text, string(content))
	}
}

// 测试TempFilenameWithText函数
func TestTempFilenameWithText(t *testing.T) {
	text := "Hello, World!"
	filename, err := TempFilenameWithText(text)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(filename) // 清理临时文件

	// 读取文件内容并验证
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read from temp file: %v", err)
	}

	if string(content) != text {
		t.Errorf("Content mismatch: expected %v, got %v", text, string(content))
	}
}
