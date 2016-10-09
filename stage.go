package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	nbworkers = 1
)

type routine func() error

type Process interface {
	Start()
	Wait()
}

type Worker interface {
	Run()
}

type Task struct {
	errors chan error
	ctx    context.Context
}

type Stage struct {
	Task
	wait *sync.WaitGroup
}

type ConcurrentStage struct {
	Stage
	concurrency int
}

type Pipeline struct {
	Task
	stages []Stage
}

type Adder struct {
	ConcurrentStage
	in  chan int
	out chan int
}

type AdderWorker struct {
	Stage
	name int
	in   chan int
	out  chan int
}

func (a *AdderWorker) Run() {
	a.wait.Add(1)
	for item := range a.in {
		v := a.name + item + 1

		select {
		case <-a.ctx.Done():
			return
		case a.out <- v:
			if item%100 == 0 {
				fmt.Printf("AdderWorker of %v doing %v : %v\n", a.name, v, time.Now())
			}
		}
	}
	fmt.Printf("AdderWorker of %v Done: %v\n", a.name, time.Now())
	a.wait.Done()
}

func (a *Adder) Start() {
	a.wait.Add(1)
	for i := 0; i < a.concurrency; i++ {
		worker := AdderWorker{
			Stage: Stage{
				Task: Task{errors: a.errors, ctx: a.ctx},
				wait: a.wait},
			name: i,
			in:   a.in,
			out:  a.out}

		go worker.Run()
	}
	a.wait.Done()
	fmt.Printf("Leave Start() : %v\n", time.Now())
}

func (a *Adder) Wait() {
	a.wait.Wait()
	close(a.out)
}

type Maxer struct {
	ConcurrentStage
	in  chan int
	out chan int
}

type MaxerWorker struct {
	Stage
	name int
	in   chan int
	out  chan int
}

func (m *MaxerWorker) Run() {
	m.wait.Add(1)
	var maxi = 0
	for item := range m.in {
		maxi = max(maxi, m.name+(item%((m.name+10)*10)))

		select {
		case <-m.ctx.Done():
			return
		case m.out <- maxi:
			if item%100 == 0 {
				fmt.Printf("MaxerWorker of %v doing %v : %v\n", m.name, maxi, time.Now())
			}
		}
	}
	m.wait.Done()
}

func (m *Maxer) Start() {
	m.wait.Add(1)
	for i := 0; i < m.concurrency; i++ {
		worker := MaxerWorker{
			Stage: Stage{
				Task: Task{errors: m.errors, ctx: m.ctx},
				wait: m.wait},
			name: i,
			in:   m.in,
			out:  m.out}

		go worker.Run()
	}
	m.wait.Done()
}

func (m *Maxer) Wait() {
	m.wait.Wait()
	close(m.out)
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
