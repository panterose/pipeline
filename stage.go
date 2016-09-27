package main

import (
	"context"
	"fmt"
	"sync"
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
		go func() {
			for i := range s.in {
				s.out <- i + 1
			}
		}()
	}
	return &grp
}

func (s Stage1) Accept(item int, ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("Cancelled")
	default:
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
		go func() {
			for i := range s.in {
				maxi = max(maxi, i)
				s.out <- maxi
			}
		}()
	}
	return &grp
}

func (s Stage2) Accept(item int, ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("Cancelled")
	default:
		s.in <- item
		return nil
	}
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
