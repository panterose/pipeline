package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPricingEngine(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1050*time.Millisecond)
	defer cancel()
	in1 := make(chan int)
	out1 := make(chan int)
	errors := make(chan error)
	var wait1 sync.WaitGroup
	adder := Adder{
		ConcurrentStage{
			Stage{
				Task{errors, ctx},
				&wait1,
			},
			10,
		},
		in1,
		out1,
	}

	//assert.Equal(t, adder..total, 0)
	var wait2 sync.WaitGroup
	in2 := out1
	out2 := make(chan int)
	maxer := Maxer{
		ConcurrentStage{
			Stage{
				Task{errors, ctx},
				&wait2,
			},
			10,
		},
		in2,
		out2,
	}

	adder.Start()
	maxer.Start()

	go func() {
		for i := 0; i < 10000; i = i + 1 {
			in1 <- i
		}
		close(in1)
		fmt.Printf("Closing in1 : %v\n", time.Now())
	}()

	go func() {
		for item := range out2 {
			select {
			case <-ctx.Done():
				fmt.Printf("Done : %v\n", time.Now())
			default:
				//fmt.Printf("Result %v: \n", item)
				if item == 0 {
					fmt.Printf("Result %v: \n", item)
				}
			}
		}
	}()

	adder.Wait()
	maxer.Wait()

}
