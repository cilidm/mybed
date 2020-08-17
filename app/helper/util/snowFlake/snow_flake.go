package snowFlake

import (
	"sync"
	"time"
)

func GetSnow() uint64 {
	NewAlgorithmSnowFlake(1, 1)
	newId := GetAlgorithmSnowFlake().GetID()
	return newId
}

type AlgorithmSnowFlake struct {
	sync.Mutex
	machineId     int64
	dataCenterId  int64
	lastTimeStamp int64
	sn            int64
}

var algorithmSnowFlake *AlgorithmSnowFlake = nil

func GetAlgorithmSnowFlake() *AlgorithmSnowFlake {
	return algorithmSnowFlake
}

func (sf *AlgorithmSnowFlake) GetID() uint64 {
	sf.Lock()
	defer sf.Unlock()

	curTimeStamp := time.Now().UnixNano() / 1000000 // 13位

	if curTimeStamp == sf.lastTimeStamp {
		if sf.sn > 4095 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano() / 1000000
			sf.sn = 0
		}
	} else {
		sf.sn = 0
	}
	sf.sn++
	sf.lastTimeStamp = curTimeStamp

	// 时间戳向左移动22位
	curTimeStamp = curTimeStamp << 22 // 35

	// 合并机器id
	machineId := sf.machineId << 17 // 18

	// 合并数据中心id
	dataCenterId := sf.dataCenterId << 12 // 13

	// 通过与运算把各个部位连接在一起
	id := curTimeStamp | machineId | dataCenterId | sf.sn
	//return id
	return uint64(id)
}

func NewAlgorithmSnowFlake(machineId int64, dataCenterId int64) {
	algorithmSnowFlake = &AlgorithmSnowFlake{
		machineId:    machineId,
		dataCenterId: dataCenterId,
	}
}
