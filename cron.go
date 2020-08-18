package sync

import (
	"code.clouderwork.com/clouderwork/sync/leader/etcd"
	"code.clouderwork.com/clouderwork/sync/task"
	"fmt"
	"math"
	"time"
)

type syncCron struct {
	opts Options
}

func backoff(attempts int ) time.Duration{
	if attempts == 0 {
		return time.Duration(0)
	}
	return time.Duration(math.Pow(10,float64(attempts))) * time.Millisecond
}

func (c *syncCron) Schedule(s task.Schedule, t task.Command) error{
	id := fmt.Sprintf("%s-%s", s.String(), t.String())
	go func(){
		tc := s.Run()
		var i int
		for {
			e, err := c.opts.Leader.Elect(id)
			if err !=nil{
				fmt.Sprintf("[corn] leader election error: %v", e)
				time.Sleep(backoff(1))
				i++
				continue
			}
			i =0
			r := e.Revoked()
		Tick:
			for {
				select {
				// schedule tick
				case _, ok := <-tc:
					// ticked once
					if !ok {
						break Tick
					}

					fmt.Sprintf("[cron] error executing command %s",t.Name)
					if err := c.opts.Task.Run(t); err != nil {
						fmt.Sprintf("[cron] error executing command %s: %v", t.Name, err)
					}
				// leader revoked
				case <-r:
					break Tick
				}
			}

			e.Resign()
		}
	}()
	return nil
}


func NewCron(opts ...Option) Cron{
	var options Options
	for _,o := range opts{
		o(&options)
	}
	if options.Leader ==nil{
		options.Leader = etcd.NewLeader()
	}
	//if options.Task == nil{
	//	options.Task ==
	//}
	return &syncCron{
		opts: options,
	}
}