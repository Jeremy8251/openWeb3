package config

import (
	"path"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// init 函数通过确定配置文件的绝对路径并将其内容解码到 Config 结构体中，初始化应用程序的配置。
// 它使用 `getCurrentAbPathByCaller` 函数来确定当前的绝对路径，并将配置文件名附加到该路径。
// 如果在解析文件路径或解码 TOML 文件时发生任何错误，该函数将通过 panic 抛出相应的错误信息。
func init() {
	// 通过 runtime.Caller(0) 获取当前执行文件的路径信息
	currentAbPath := getCurrentAbPathByCaller()
	// 获取当前执行文件的目录路径
	tomlFile, err := filepath.Abs(currentAbPath + "/configV21.toml")
	//tomlFile, err := filepath.Abs(currentAbPath + "/configV22.toml")
	if err != nil {
		panic("read toml file err: " + err.Error())
		return
	}
	// 解析 toml 文件到 Config 结构体
	if _, err := toml.DecodeFile(tomlFile, &Config); err != nil {
		panic("read toml file err: " + err.Error())
		return
	}
}

func getCurrentAbPathByCaller() string {
	var abPath string
	// 通过 runtime.Caller(0) 获取当前执行文件的路径信息
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		// 获取当前执行文件的目录路径
		abPath = path.Dir(filename)
	}
	return abPath
}
