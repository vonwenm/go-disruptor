package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)

	producerSequence := NewSequence()
	consumerSequence := NewSequence()

	producerBarrier := NewBarrier(producerSequence)
	consumerBarrier := NewBarrier(consumerSequence)

	sequencer := NewSingleProducerSequencer(producerSequence, RingSize, consumerBarrier)
	go consume(producerBarrier, producerSequence, consumerSequence)

	started := time.Now()
	for i := int64(0); i < MaxSequenceValue; i++ {
		ticket := sequencer.Next(1)
		ringBuffer[ticket&RingMask] = ticket
		sequencer.Publish(ticket)
		if ticket%Mod == 0 {
			finished := time.Now()
			elapsed := finished.Sub(started)
			fmt.Println(ticket, elapsed)
			started = time.Now()
		}
	}
}

func consume(barrier Barrier, source, sequence *Sequence) {
	worker := NewWorker(barrier, TestHandler{}, source, sequence)

	for {
		worker.Process()
	}
}

const Mod = 1000000 * 10 // 1 million * 10
const RingSize = 1024 * 256
const RingMask = RingSize - 1

// RingMask = RingSize - 1 // slightly faster than a mod operation
// 0&3 = 0
// 1&3 = 1
// 2&3 = 2
// 3&3 = 3
// 4&3 = 0

var ringBuffer [RingSize]int64

type TestHandler struct{}

func (this TestHandler) Consume(sequence, remaining int64) {
	message := ringBuffer[sequence&RingMask]
	if message != sequence {
		panic(fmt.Sprintf("\n", message, sequence))
	} else if sequence%Mod == 0 {
		// fmt.Println("Current Sequence:", sequence)
	}
}