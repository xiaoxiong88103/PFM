package Vars

import (
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
)

func CheckAndCreateFileJson(fileName string) error {
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

// CheckAndCreateINI 检查配置文件是否存在，如果不存在则创建；如果存在但缺少字段，则补充字段
func CheckAndCreateINI(filePath string) error {
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 文件不存在，创建新的配置文件
		cfg := ini.Empty()

		// 添加默认的 [white_list] 和 [black_list] 字段
		cfg.NewSection("white_list")
		cfg.NewSection("black_list")

		// 保存到文件
		if err := cfg.SaveTo(filePath); err != nil {
			return err
		}
		log.Printf("配置文件 %s 不存在，已创建默认配置", filePath)
		return nil
	}

	// 文件存在，加载配置文件
	cfg, err := ini.Load(filePath)
	if err != nil {
		return err
	}

	// 检查 [white_list] 是否存在
	if !cfg.HasSection("white_list") {
		cfg.NewSection("white_list")
		log.Println("配置文件中缺少 [white_list] 字段，已补充")
	}

	// 检查 [black_list] 是否存在
	if !cfg.HasSection("black_list") {
		cfg.NewSection("black_list")
		log.Println("配置文件中缺少 [black_list] 字段，已补充")
	}

	// 保存修改后的配置文件
	if err := cfg.SaveTo(filePath); err != nil {
		return err
	}

	log.Printf("配置文件 %s 检查完成", filePath)
	return nil
}
