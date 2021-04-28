package billow

import (
	"context"
	"log"
	"sync"
	"time"
)

type JobFunc func() error

type BillowsHome interface {
	Stop()
	Add(Billow)
	Start()
	Monitor()
	ApplyJob()
}

type CancelCtx struct {
	Ctx context.Context
	f   context.CancelFunc
}

type ErrorJob struct {
	Err     error
	JobInfo JobTrans
}

// 池
type Billow struct {
	//BillowChan chan *Job
	errorChan  chan ErrorJob
	jobChan    chan *Job
	stopBillow chan struct{}
	workers    chan *Worker
	stopT      int64
	startT     int64
	ctx        CancelCtx
	jobSize    int
	rg         *int32 // 当前已经在运行的任务
	mu         sync.Mutex
	size       int
}

type Worker struct {
	i int
}

func (worker *Worker) Do(job *Job) error {
	job.Tran.StartT = time.Now().Unix()
	return job.f()
}

type JobHome interface {
	Exec()
}

type JobTrans struct {
	Id     int64
	Name   string
	StartT int64
}

type Job struct {
	Tran        JobTrans
	f           JobFunc
	CanStopChan chan int
	NoStopChan  chan int
	Size        int
	ExecT       int64
}

func InitBillow(size int) (bw *Billow) {
	ctx, f := context.WithCancel(context.Background())
	initIdStorage(size)
	return &Billow{
		errorChan:  make(chan ErrorJob),
		stopBillow: make(chan struct{}),
		startT:     time.Now().Unix(),
		jobChan:    make(chan *Job, size),
		workers:    make(chan *Worker, size),
		ctx:        CancelCtx{ctx, f},
		size:       size,
	}
}

func (job *Job) Exec(bw *Billow) {
	log.Println("--------")
	for {
		select {
		case worker := <-bw.workers:
			//go func() {
				dressI(worker.i, false)
				log.Println("任务[", job.Tran.Name, "]开始执行，当前正在运行的worker数量A=", getI())
				if err := worker.Do(job); err != nil {
					log.Fatalf("%+v exec err:%v", job, err.Error())
					bw.errorChan <- ErrorJob{Err: err, JobInfo: job.Tran} // 错误日志及时消费避免阻塞正常任务执行
				}
				job.ExecT = time.Now().Unix() - job.Tran.StartT
			//}()
			dressI(worker.i, true)
			log.Println(">>> [", job.Tran.Name, "]执行完毕，", "当前worker数量1=", len(bw.workers))
			if len(bw.workers) < bw.size {
				bw.workers <- worker
				log.Println("当前正在运行的worker数量B=", getI())
				log.Println(">>> 当前worker数量2=", len(bw.workers))
			}
		}
	}
	return
}
