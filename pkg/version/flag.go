// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package version

// 自定义的版本标志系统
// 具体见pflag用法
import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

// 定义版本标识的类型
type versionValue int

const (
	// 未设置版本
	VersionNotSet versionValue = 0
	// 启用版本
	VersionEnabled versionValue = 1
	// 详细版本
	VersionRaw versionValue = 2
)

const (
	strRawVersion   = "raw"
	versionFlagName = "version"
)

// 定义版本标志
var versionFlag = Version(versionFlagName, VersionNotSet, "Print version information and quit.")

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

// 实现pflag.Value接口中的String方法
func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}
	return strconv.FormatBool(bool(*v == VersionEnabled))
}

// 实现pflag.Value接口中的Set方法
func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionEnabled
	} else {
		*v = VersionNotSet
	}
	return err
}

// 实现pflag.Value接口中的Type方法
func (v *versionValue) Type() string {
	return "version"
}

// VersionVar定义了一个具有指定名称和用法的标志
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	pflag.Var(p, name, usage)
	// 当只使用--version不带参数时，默认值为true
	pflag.Lookup(name).NoOptDefVal = "true"
}

func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

// AddFlags在任意FlagSet上注册这个包的标志，使他们指向与全局标志相同的值
func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(versionFlagName))
}

// 检查是否传递 --version 若是，打印版本并退出
func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%s\n", Get().Text())
		os.Exit(0)
	} else if *versionFlag == VersionEnabled {
		fmt.Printf("%s\n", Get().String())
		os.Exit(0)
	}
}
