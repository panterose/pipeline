package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Stage interface {
	Init() *sync.WaitGroup
	Accept(item int, ctx context.Context)
	Complete()
}

type Stage1 struct {
	total int
	in    chan int
	out   chan int
}

func (s *Stage1) Init() *sync.WaitGroup {
	var grp sync.WaitGroup
	grp.Add(5)
	for i := 0; i < 5; i++ {
		go func(name int) {
			for item := range s.in {
				fmt.Printf("Stage1 %v item %v: %v\n", name, item, time.Now())
				s.out <- i + 1
			}
			grp.Done()
		}(i)
	}
	return &grp
}

func (s Stage1) Accept(item int, ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("Cancelled")
	default:
		fmt.Printf("Stage1 Accept %v: %v\n", item, time.Now())
		s.in <- item
		return nil
	}
}

func (s Stage1) Complete() {
	close(s.out)
}

type Stage2 struct {
	max int
	in  chan int
	out chan int
}

func (s *Stage2) Init() *sync.WaitGroup {
	var maxi = 0
	var grp sync.WaitGroup
	grp.Add(5)
	for i := 0; i < 5; i++ {
		go func(name int) {
			for item := range s.in {
				maxi = max(maxi, item)
				fmt.Printf("Stage2 %v max %v: %v\n", name, maxi, time.Now())
				s.out <- maxi
			}
			grp.Done()
		}(i)
	}
	return &grp
}

func (s Stage2) Accept(item int, ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("Cancelled")
	default:
		fmt.Printf("Stage2 Accept %v: %v\n", item, time.Now())
		s.in <- item
		return nil
	}
}

func (s Stage2) Complete() {
	close(s.out)
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
