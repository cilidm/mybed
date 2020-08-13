package govatar

import (
	"github.com/o1egl/govatar"
	"path/filepath"
	"strconv"
)

// 随机生成头像
func CreateAvatar(name, path string, num int) {
	for i := 1; i <= num; i++ {
		govatar.GenerateFile(govatar.FEMALE, filepath.Join(path, name+"_"+strconv.Itoa(i)+"_0.jpg"))
		govatar.GenerateFileForUsername(govatar.MALE, "username"+strconv.Itoa(i), filepath.Join(path, strconv.Itoa(i)+"_1.jpg"))
	}
}
