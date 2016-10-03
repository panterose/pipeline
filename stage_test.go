package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPricingEngine(t *testing.T) {
	in1 := make(chan int)
	out1 := make(chan int)
	s1 := Stage1{0, in1, out1}

	assert.Equal(t, s1.total, 0)

	in2 := out1
	out2 := make(chan int)
	s2 := Stage2{0, in2, out2}

	assert.Equal(t, s2.max, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 1050*time.Millisecond)
	defer cancel()
	s1.Init(ctx)
	s2.Init(ctx)

	go func() {
		for i := 0; i < 10000; i = i + 1 {
			in1 <- i
		}
		close(in1)
	}()

	for item := range out2 {
		select {
		case <-ctx.Done():
			fmt.Printf("Done : %v\n", time.Now())
		default:
			fmt.Printf("Result %v: %v\n", item, time.Now())
		}
	}

}
