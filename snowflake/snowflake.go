package snowflake

import (
	"sync"
	"time"
)

// 基本配置
var (
	// 偏移时间 2018-01-01 00:00:00 的时间戳-秒
	epoch int64 = 1514764800

	// 时间戳占位
	timeBits uint8 = 41
	// worker占位
	workerBits uint8 = 4
	// 序列器占位
	sequenceBits uint8 = 18

	// shift
	timeShift uint8 = workerBits + sequenceBits
	workerShift uint8 = sequenceBits

	// max
	sequenceMax uint32 = 1 << sequenceBits

)

var SnowFlakeInstance *SnowFlake

func init(){
	SnowFlakeInstance = &SnowFlake{
		Epoch:epoch,
		TimeBits:timeBits,
		WorkerBits:workerBits,
		SequenceBits:sequenceBits,
		TimeShift:timeShift,
		WorkerShift:workerShift,
		CurrentWorkerID: 1,  // 机器id, 可从配置文件读取
		CurrentSequence:0,
		CurrentTimeStamp:time.Now().Unix() - epoch,
	}
}

type SnowFlake struct {
	// 实例锁
	mt sync.RWMutex

	Epoch int64  // 偏移时间
	TimeBits uint8
	WorkerBits uint8
	SequenceBits uint8

	TimeShift uint8
	WorkerShift uint8

	// 当前时间戳
	CurrentTimeStamp int64
	// 当前worker号
	CurrentWorkerID uint32
	// 当前序列器
	CurrentSequence uint32
}

// 生成ID
func (sf *SnowFlake)GenerateID()int64{
	sf.mt.Lock()
	defer sf.mt.Unlock()

	// 序列器超了
	if sf.CurrentSequence >= sequenceMax {
		time.Sleep(1 * time.Second)
	}

	ctt := sf.getTimeStamp()
	var id int64

	// 如何在同一秒内 自增序列器
	if sf.CurrentTimeStamp == ctt {
		sf.CurrentSequence++
		id = sf.generateID()
	}else{
		// 更新&复位
		sf.CurrentTimeStamp = ctt
		sf.CurrentSequence = 0
		id = sf.generateID()
	}
	return id
}

func (sf *SnowFlake)generateID()int64{
	var id int64 = 0
	// 时间戳
	id = id | (sf.CurrentTimeStamp << sf.TimeShift)
	// worker_id
	id = id | int64(sf.CurrentWorkerID << sf.WorkerShift)
	// 序列器
	id = id | int64(sf.CurrentSequence)
	return id
}

// 获取时间 - 秒级别
func (sf *SnowFlake)getTimeStamp()int64{
	return time.Now().Unix() - sf.Epoch
}