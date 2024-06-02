package runtime

import (
	"container/list"
	"go.uber.org/zap"
	"os/exec"
	"time"
)

type STRUCT_SOUND_QUEUE struct {
	sem  chan int
	list list.List
}

type Sound_Play struct {
	sound_queue STRUCT_SOUND_QUEUE
}

func (sp *Sound_Play) Run() {
	sp.sound_queue = STRUCT_SOUND_QUEUE{make(chan int, 1), list.List{}}
	go sp.Sound_start()
}

func (sp *Sound_Play) Sound_start() {
	zap.S().Infoln("Sound_Play Running......")
	for {
		if sp.sound_queue.list.Len() > 0 {
			sp.sound_queue.sem <- 1
			str := sp.sound_queue.list.Back()
			sp.sound_queue.list.Remove(str)
			<-sp.sound_queue.sem

			exec.Command("play", str.Value.([]string)...).Run()
			continue
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (sp *Sound_Play) Add_Sound_Queue(str []string) {
	sp.sound_queue.sem <- 1
	_ = sp.sound_queue.list.PushFront(str)
	<-sp.sound_queue.sem
}
