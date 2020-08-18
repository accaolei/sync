package leader

type Leader interface{
	Elect(id string, opts ...ElectOption) (Elected,error)
	Follow() chan string
}

type Elected interface {
	Id() string
	Reelect() error
	Resign() error
	Revoked() chan bool
}

type Option func(o *Options)



type ElectOption func(o *ElectionOptions)