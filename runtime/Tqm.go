package runtime

import (
	"TmdtServer/common"
	"sync"
	"time"
)

// TemperatureQueue represents the queue to store temperature data
type TemperatureQueue struct {
	Queue     []DeviceData
	MaxLength int
	mutex     sync.Mutex
}

// NewTemperatureQueue creates a new TemperatureQueue with a given length
func NewTemperatureQueue(length int) *TemperatureQueue {
	return &TemperatureQueue{
		Queue:     make([]DeviceData, 0, length),
		MaxLength: length,
	}
}

// Add adds a new DeviceData to the queue
func (tq *TemperatureQueue) Add(data DeviceData) {
	tq.mutex.Lock()
	defer tq.mutex.Unlock()

	if len(tq.Queue) >= tq.MaxLength {
		// Remove the oldest data
		tq.Queue = tq.Queue[1:]
	}
	// Add new data to the end of the queue
	tq.Queue = append(tq.Queue, data)
}

// TemperatureQueueManager manages multiple TemperatureQueues
type TemperatureQueueManager struct {
	queues      map[uint32]*TemperatureQueue
	queueLength int
	dataMap     *map[uint32]DeviceData
	dataMutex   *sync.Mutex
}

// NewTemperatureQueueManager creates a new TemperatureQueueManager
func NewTemperatureQueueManager(queueLength int, dataMap *map[uint32]DeviceData, dataMutex *sync.Mutex) *TemperatureQueueManager {
	return &TemperatureQueueManager{
		queues:      make(map[uint32]*TemperatureQueue),
		queueLength: queueLength,
		dataMap:     dataMap,
		dataMutex:   dataMutex,
	}
}

// Start periodically stores data from the Server's dataMap into queues
func (tqm *TemperatureQueueManager) Start() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			tqm.storeData()
		}
	}()
}

func (tqm *TemperatureQueueManager) storeData() {
	tqm.dataMutex.Lock()
	defer tqm.dataMutex.Unlock()

	for deviceID, data := range *tqm.dataMap {
		queue, exists := tqm.queues[deviceID]
		if !exists {
			queue = NewTemperatureQueue(tqm.queueLength)
			tqm.queues[deviceID] = queue
		}
		queue.Add(data)
		if *common.TestMode {
			//fmt.Printf("Stored data for device %d: %+v\n", deviceID, data)
			//fmt.Println(queue)
		}
	}
}
