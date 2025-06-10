// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package version

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/gosuri/uitable"
)

// 保存版本信息

var (
	// 语义化的版本号默认值，实际使用中会通过 -ldflags参数在编译时赋值为实际的版本号
	gitVersion = "v0.0.0-master+&Format:%h$"
	// ISO8601格式的构建时间 $(date -u + '%Y-%m-%dT%H:%M:%SZ')命令的输出
	// 实际使用中会通过 -ldflags参数在编译时赋值为构建时的时间戳
	buildDate = "1970-01-01T00:00:00Z"
	// git的SHA1值，$(git rev-parse HEAD)命令的输出
	gitCommit = "$Format:%H$"
	// 代表构建时git仓库的状态
	gitTreeState = ""
)

type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

func (info Info) String() string {
	return info.GitVersion
}

func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)
	return string(s)
}

// Text将版本信息编码为UTF-8
func (info Info) Text() string {
	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow("gitVersion:", info.GitVersion)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("gitTreeState:", info.GitTreeState)
	table.AddRow("buildDate:", info.BuildDate)
	table.AddRow("goVersion:", info.GoVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)
	return table.String()
}

// 返回详尽的代码库版本信息，用来标明二进制文件由哪个版本的代码构建
func Get() Info {
	return Info{
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
