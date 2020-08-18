package task

import (
	"context"
	"fmt"
	"time"
)

type Task interface {
	Run(Command) error
	Status() string
}

type Command struct {
	Name string
	Func func() error
}

type Schedule struct {
	Time time.Time
	Interval time.Duration
}

type Options struct{
	Pool int
	Context context.Context
}

type Option func(o *Options)

func (c Command) Execute() error{
	return c.Func()
}

func (c Command) String() string{
	return c.Name
}

func (s Schedule) Run() <-chan time.Time{

	//计算时间差
	d := s.Time.Sub(time.Now())
	ch := make(chan time.Time,1)
	go func(){

		//等待开始时间
		<-time.After(d)

		// 没有间隔
		if s.Interval == time.Duration(0){
			ch <- time.Now()
			close(ch)
			return
		}

		// 启动滴答器
		ticker := time.NewTicker(s.Interval)
		defer ticker.Stop()
		for t:= range ticker.C{
			ch <- t
		}
	}()
	return ch
}

func (s Schedule) String() string{
	return fmt.Sprintf("%d-%d", s.Time.Unix(), s.Interval)
}

func WithPool(i int) Option{
	return func(o *Options){
		o.Pool = i
	}
}



