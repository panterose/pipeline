package main

import (
	"testing"

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

	sg1 := stage1.Init()
	sg2 := stage2.Init()

  go func() {
    for i := 0; i<10000; i+1 {
       in1 <- i
    }
    close(in1)
  }()

  sg1.

}
