package chanUtil

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"MIT6.824-6.5840/util/color"
)

func TestBroadcast_Example(t *testing.T) {
	t.Log(color.Cyan(fmt.Sprintf("begin test,number of goroutine: %d", runtime.NumGoroutine())))
	var callbacks []func()
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("listener_%d", i)
		callbacks = append(
			callbacks, func() {
				t.Log(color.Green(name), "running")
				time.Sleep(time.Second + time.Duration(rand.Int()%1000)*time.Millisecond)
				t.Log(color.Green(name), "run end")
			},
		)
	}
	broadcaster := NewBroadcaster(callbacks)

	t.Log(color.Cyan(fmt.Sprintf("after start broadcast,number of goroutine: %d", runtime.NumGoroutine())))

	t.Log(color.Magenta("First broadcast start"))
	broadcaster.Broadcast()
	t.Log(color.Magenta("First broadcast end"))
	time.Sleep(800 * time.Millisecond)

	t.Log(color.Magenta("Second broadcast start"))
	broadcaster.Broadcast()
	t.Log(color.Magenta("Second broadcast end"))
	time.Sleep(4 * time.Second)
	broadcaster.Close()
	time.Sleep(100 * time.Millisecond)
	t.Log(color.Cyan(fmt.Sprintf("after close broadcast,number of goroutine: %d", runtime.NumGoroutine())))

}

func TestBlockingThrottler_DoubleRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	var count atomic.Int32
	throttler := NewBlockingThrottler(
		func() {
			fmt.Println("Executing task count:", count.Add(1))
			time.Sleep(2 * time.Second)
			fmt.Println("Executing task end")
		}, ctx,
	)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		throttler.Run()
	}()

	go func() {
		defer wg.Done()
		throttler.Run()
	}()

	go func() {
		defer wg.Done()
		throttler.Run()
	}()

	wg.Wait()
	if count.Load() != 2 {
		t.Fail()
	}
}
