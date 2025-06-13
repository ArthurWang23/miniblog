package rid

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"os"
)

// 计算机器ID的哈希值并返回一个uint64类型的盐值
func Salt() uint64 {
	hasher := fnv.New64a()
	hasher.Write(ReadMachineID())

	hashValue := hasher.Sum64()
	return hashValue
}

// 获取机器ID，若无法获取则生成随机ID
func ReadMachineID() []byte {
	id := make([]byte, 3)
	machineID, err := readPlatformMachineID()
	if err != nil || len(machineID) == 0 {
		machineID, err = os.Hostname()
	}
	if err == nil && len(machineID) != 0 {
		hasher := sha256.New()
		hasher.Write([]byte(machineID))
		copy(id, hasher.Sum(nil))
	} else {
		// 无法收集
		if _, randErr := rand.Reader.Read(id); randErr != nil {
			panic(fmt.Errorf("id:cannot get hostname nor generate a random number: %w;%w", err, randErr))
		}
	}
	return id
}

func readPlatformMachineID() (string, error) {
	data, err := os.ReadFile("/etc/machine-id")
	if err != nil || len(data) == 0 {
		data, err = os.ReadFile("/sys/class/dmi/id/product_uuid")
	}
	return string(data), err
}
