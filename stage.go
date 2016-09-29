package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	nbworkers = 5
)

type routine func() error

type Stage interface {
	Init()
	Accept(item int, ctx context.Context)
	Complete()
}

type Stage1 struct {
	total int
	in    chan int
	out   chan int
}

func (s *Stage1) Init(parent context.Context) {
	g, ctx := errgroup.WithContext(parent)
	for i := 0; i < nbworkers; i++ {
		g.Go(s.engine(i, ctx))
	}
	go s.cleaner(g)
}

func (s *Stage1) engine(name int, ctx context.Context) routine {
	return func() error {
		for item := range s.in {
			v := name + item + 1

			select {
			case <-ctx.Done():
				return fmt.Errorf("Cancelled1 %v", name)
			case s.out <- v:
				if item%100 == 0 {
					fmt.Printf("Stage1 of %v doing %v : %v\n", name, v, time.Now())
				}
			}
		}
		return nil
	}
}

func (s *Stage1) cleaner(g *errgroup.Group) {
	error := g.Wait()
	fmt.Printf("Stage1 error %v : %v\n", error, time.Now())
	close(s.out)
}

type Stage2 struct {
	max int
	in  chan int
	out chan int
}

func (s *Stage2) Init(parent context.Context) {

	g, ctx := errgroup.WithContext(parent)
	for i := 0; i < nbworkers; i++ {
		g.Go(s.engine(i, ctx))
	}

	go s.cleaner(g)
}

func (s *Stage2) cleaner(g *errgroup.Group) {
	error := g.Wait()
	fmt.Printf("Stage2 error %v : %v\n", error, time.Now())
	close(s.out)
}

func (s *Stage2) engine(name int, ctx context.Context) routine {
	return func() error {
		var maxi = 0
		for item := range s.in {
			maxi = max(maxi, (name+item)%100)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case s.out <- maxi:
				if item%100 == 0 {
					fmt.Printf("Stage2 of %v doing %v : %v\n", name, maxi, time.Now())
				}
			}
		}
		return nil
	}
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
