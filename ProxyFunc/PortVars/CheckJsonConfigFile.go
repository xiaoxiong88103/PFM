package PortVars

import (
	"encoding/json"
	"fmt"
	"os"
)

func CheckAndCreateFile(fileName string) error {
	// 检查文件是否存在
	_, err := os.Stat(fileName)
	if err == nil {
		// 文件已存在
		fmt.Printf("File %s already exists.\n", fileName)
		return nil
	}

	if os.IsNotExist(err) {
		// 文件不存在，创建文件并写入空的 JSON 对象
		file, createErr := os.Create(fileName)
		if createErr != nil {
			return fmt.Errorf("failed to create file: %w", createErr)
		}
		defer file.Close()

		// 写入空 JSON 对象 {}
		emptyJSON := make(map[string]interface{})
		jsonData, marshalErr := json.MarshalIndent(emptyJSON, "", "  ")
		if marshalErr != nil {
			return fmt.Errorf("failed to marshal JSON: %w", marshalErr)
		}

		_, writeErr := file.Write(jsonData)
		if writeErr != nil {
			return fmt.Errorf("failed to write to file: %w", writeErr)
		}

		fmt.Printf("File %s created successfully.\n", fileName)
		return nil
	}

	// 其他错误
	return fmt.Errorf("failed to check file: %w", err)
}
