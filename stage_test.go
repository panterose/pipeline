package main

import (
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

	sg1 := s1.Init()
	sg2 := s2.Init()

	go func() {
		for i := 0; i < 10000; i = i + 1 {
			in1 <- i
		}
		close(in1)
	}()

	for out := range out2 {
		fmt.Printf("Result %v: %v\n", out, time.Now())
	}

	sg1.Wait()
	s1.Complete()
	sg2.Wait()
	s2.Complete()
}
