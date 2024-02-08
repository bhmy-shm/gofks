package snowflake

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"strings"
)

// SnowflakeUUid 仅生成一个uuid
func SnowflakeUUid(nodeId int64) string {
	// 创建一个新的雪花算法实例
	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		logx.Error("Failed to create snowflake node:", err)
		return ""
	}

	return node.Generate().String()
}

func GenerateUniqueID(length int) string {
	// 生成随机的UUID
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		errorx.Fatal(err)
	}

	// 将UUID转换为字符串形式，去掉分隔符和空格
	uuidString := strings.ReplaceAll(randomUUID.String(), "-", "")

	// 截取指定长度的字符串
	if length > len(uuidString) {
		length = len(uuidString)
	}
	uuidString = uuidString[:length]

	return uuidString
}
