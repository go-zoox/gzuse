package core

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-zoox/logger"
	"github.com/pbnjay/memory"
)

type Config struct {
	MemoryPercent uint   `json:"memory_percent"`
	CPUPercent    uint   `json:"cpu_percent"`
	MemorySize    string `json:"memory_size"`
	CPUCore       uint   `json:"cpu_cores"`
}

func Run(cfg *Config) error {
	isInUse := false

	if cfg.MemorySize != "" || cfg.MemoryPercent != 0 {
		isInUse = true

		if err := useMemory(cfg.MemoryPercent, cfg.MemorySize); err != nil {
			return err
		}
	}

	if cfg.CPUPercent > 0 || cfg.CPUCore > 0 {
		isInUse = true

		useCPU(cfg.CPUPercent, cfg.CPUCore)
	}

	if !isInUse {
		return fmt.Errorf("memory or cpu are required to set")
	}

	// 持续输出信息
	for {
		time.Sleep(time.Second)
	}

	// return nil
}

func useMemory(percent uint, size string) (err error) {
	var sizeX uint64 = 0
	if size != "" {
		sizeX, err = humanize.ParseBytes(size)
		if err != nil {
			return fmt.Errorf("invalid memory(%s): %s", size, err)
		}

		logger.Infof("[memory][size] %s", size)
	} else if percent != 0 {
		// 获取系统总内存大小
		totalMem := memory.TotalMemory()

		// 计算要占用的内存大小
		memBytes := float64(totalMem) * float64(percent) / 100
		memBytes = math.Floor(memBytes/1024/1024) * 1024 * 1024 // 四舍五入到最接近的 MB 数量
		sizeX = uint64(memBytes)

		logger.Infof("[memory][percent] %d%% (size: %s)", percent, humanize.Bytes(sizeX))
	}

	data := make([]byte, sizeX)
	for i := uint64(0); i < sizeX; i++ {
		data[i] = byte(i % 256)
	}

	return nil
}

func useCPU(percent uint, cores uint) error {
	if cores == 0 {
		cores = uint(runtime.NumCPU())
	}

	if percent == 0 {
		percent = 1
	}

	logger.Infof("[cpu][cores] %d", cores)
	logger.Infof("[cpu][percent] %d", percent)

	// 创建 WaitGroup 等待所有 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(int(cores))

	// 启动每个 CPU 核心的占用操作
	for i := uint(0); i < cores; i++ {
		go func(num uint) {
			defer wg.Done()

			consumeCPU(int(percent))
		}(i)
	}

	wg.Wait()

	return nil
}

// 占用单个 CPU 核心指定百分比的函数
func consumeCPU(percentage int) {
	for {
		// 获取当前时间
		start := time.Now()

		// 消耗 CPU 时间
		for time.Since(start).Seconds() < float64(percentage)/100 {
			runtime.LockOSThread()
		}

		// 休眠一段时间以降低 CPU 使用率
		time.Sleep(time.Duration(100-percentage) * time.Millisecond)
	}
}
