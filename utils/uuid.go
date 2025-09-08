package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// GenerateOfflinePlayerUUID 根据角色名称生成与离线验证系统兼容的UUID
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E5%85%BC%E5%AE%B9%E7%A6%BB%E7%BA%BF%E9%AA%8C%E8%AF%81
func GenerateOfflinePlayerUUID(playerName string) (string, error) {
	hashStr := "OfflinePlayer:" + playerName
	data := []byte(hashStr)

	hash := md5.Sum(data)

	// Set the version number (3) in the 7th byte (index 6)
	hash[6] &= 0x0f // Clear the version bits
	hash[6] |= 0x30 // Set version to 3

	// Set the variant to IETF in the 9th byte (index 8)
	hash[8] &= 0x3f // Clear the variant bits
	hash[8] |= 0x80 // Set to IETF variant
	// 转换为UUID对象
	parsedUUID, err := uuid.FromBytes(hash[:])
	if err != nil {
		return "", err
	}

	// 返回无符号UUID字符串（去掉-）
	return strings.ReplaceAll(parsedUUID.String(), "-", ""), nil
}

// FormatUUID 将无符号UUID字符串转换为标准格式的UUID字符串
func FormatUUID(undashedUUID string) (string, error) {
	// 检查输入长度
	if len(undashedUUID) != 32 {
		return "", errors.New("invalid undashed UUID length")
	}

	// 格式化UUID
	formatted := undashedUUID[0:8] + "-" +
		undashedUUID[8:12] + "-" +
		undashedUUID[12:16] + "-" +
		undashedUUID[16:20] + "-" +
		undashedUUID[20:32]

	// 验证格式是否正确
	if _, err := uuid.Parse(formatted); err != nil {
		return "", err
	}

	return formatted, nil
}

// ParseUndashedUUID 将无符号UUID字符串解析为UUID对象
func ParseUndashedUUID(undashedUUID string) (uuid.UUID, error) {
	formatted, err := FormatUUID(undashedUUID)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(formatted)
}

// GenerateUUID 生成一个随机的无符号UUID字符串（Version 4）
func GenerateUUID() string {
	newUUID := uuid.New()
	return strings.ReplaceAll(newUUID.String(), "-", "")
}

// ValidateUndashedUUID 验证无符号UUID字符串是否有效
func ValidateUndashedUUID(undashedUUID string) bool {
	if len(undashedUUID) != 32 {
		return false
	}

	// 检查是否都是十六进制字符
	_, err := hex.DecodeString(undashedUUID)
	if err != nil {
		return false
	}

	// 检查是否可以解析为UUID
	_, err = ParseUndashedUUID(undashedUUID)
	return err == nil
}
