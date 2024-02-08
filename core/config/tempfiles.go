package pkg

import (
	"github.com/bhmy-shm/gofks/core/utils/hash"
	"os"
)

// TempFileWithText 传入text 创建文件并写入内容
// 注意：调用者需要在外部关闭文件句柄，并通过 file.Stat()获取文件信息，再通过filename 手动删除临时文件 os.Remove
func TempFileWithText(text string) (*os.File, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), hash.Md5Hex([]byte(text)))
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(tmpFile.Name(), []byte(text), os.ModeTemporary); err != nil {
		return nil, err
	}

	return tmpFile, nil
}

// TempFilenameWithText 创建文件并写入指定内容, 返回文件完整的绝对路径。
// 注意：调用者无需手动关闭文件句柄，但需要调用 os.Remove(filename) 函数来临时删除文件
func TempFilenameWithText(text string) (string, error) {
	tmpFile, err := TempFileWithText(text)
	if err != nil {
		return "", err
	}

	filename := tmpFile.Name()
	if err = tmpFile.Close(); err != nil {
		return "", err
	}

	return filename, nil
}
