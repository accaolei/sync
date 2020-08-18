package etcd

import (
	"code.clouderwork.com/clouderwork/sync/leader"
	"context"
	client "github.com/coreos/etcd/clientv3"
	cc "github.com/coreos/etcd/clientv3/concurrency"
	"log"
	"path"
	"strings"
)

type etcdLeader struct {
	opts leader.Options
	path string
	client *client.Client
}

type etcdElected struct{
	s *cc.Session
	e *cc.Election
	id string
}

func (e *etcdLeader) Elect(id string, opts ...leader.ElectOption)(leader.Elected,error){
	var options leader.ElectionOptions
	for _,o :=range opts{
		o(&options)
	}

	path := path.Join(e.path, strings.Replace(id, "/","-",-1))
	s, err := cc.NewSession(e.client)
	if err !=nil{
		return nil,err
	}
	l := cc.NewElection(s,path)
	if err:= l.Campaign(context.TODO(), id); err !=nil{
		return nil,err
	}
	return &etcdElected{
		e:l,id:id,
	},nil

}

func (e *etcdLeader) Follow() chan string{
	ch := make(chan string)
	s,err := cc.NewSession(e.client)
	if err !=nil{
		return ch
	}
	l := cc.NewElection(s,e.path)
	ech :=l.Observe(context.Background())
	go func(){
		for r :=range ech {
			ch <- string(r.Kvs[0].Value)
		}
	}()
	return ch
}

func (e *etcdElected) Reelect() error{
	return e.e.Campaign(context.TODO(),e.id)
}

func (e *etcdElected) Revoked() chan bool {
	ch := make(chan bool, 1)
	ech := e.e.Observe(context.Background())

	go func() {
		for r := range ech {
			if string(r.Kvs[0].Value) != e.id {
				ch <- true
				close(ch)
				return
			}
		}
	}()

	return ch
}

func (e *etcdElected) Resign() error {
	return e.e.Resign(context.Background())
}

func (e *etcdElected) Id() string {
	return e.id
}


func NewLeader(opts ...leader.Option) leader.Leader {
	var options leader.Options
	for _,o :=range opts{
		o(&options)
	}
	var endpoints []string
	for _, addr := range options.Nodes{
		if len(addr)>0{
			endpoints = append(endpoints,addr)
		}
	}
	if len(endpoints) == 0{
		endpoints =[]string{"http://127.0.0.1:2379"}
	}

	c, err := client.New(client.Config{
		Endpoints: endpoints,
	})
	if err !=nil{
		log.Fatal(err)
	}
	return &etcdLeader{
		path: "/sync/leader",
		client: c,
		opts: options,
	}

}