package sync

import  (
	"code.clouderwork.com/clouderwork/sync/leader"
	"code.clouderwork.com/clouderwork/sync/task"
	"time"
)

//Cron 是一种使用了leader选举的分布式调度程序和分布式任务管理器，依赖etcd
type Cron interface {
	Schedule(task.Schedule, task.Command) error
}

type Options struct {
	Leader leader.Leader
	Task task.Task
	Time time.Time
}

type Option func(o *Options)