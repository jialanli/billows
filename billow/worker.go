package billow

import (
	"fmt"
	lacia "github.com/jialanli/lacia/utils"
	"log"
	"time"
)

// 停止系统工作  客户发起调用 stop the world
func (bw *Billow) Stop() {
	bw.stopT = time.Now().Unix()
	bw.ctx.f()
	defer close(bw.stopBillow)
}

func (bw *Billow) Start() {
	fmt.Println("启动...")
	defer func() {
		if err := recover(); err != nil {
			bw.Stop()
		}
	}()

	n := bw.size
	if len(bw.workers) == 0 {
		fmt.Println("添加worker")
		for i := 0; i < n; i++ {
			bw.workers <- &Worker{i: generateI(n)}
		}
		log.Println("worker添加完毕=",len(bw.workers))
	}

	for {
		select {
		case job := <-bw.jobChan:
			fmt.Println("读取到任务")
			job.Exec(bw)
		case <-bw.ctx.Ctx.Done():
			return
		default:
		}
	}
}

// 监控整个系统工作是否正常
func (bw *Billow) Monitor() {
	for {
		select {
		case <-bw.stopBillow:
			return
		case errObj := <-bw.errorChan:
			fmt.Printf("[%s][Monitor]:job[%+v] exec error: %v",
				lacia.GetTimeStrOfDayTimeByTs(time.Now().Unix()), errObj.JobInfo, errObj.Err.Error())
		}
	}
}

func CalcExecTimeCall(bw *Billow) {

}


//func (bw *Billow) ApplyJob() {
//	go bw.Monitor()
//	go func() {
//		for {
//			select {
//			// 接收到任务
//			case job := <-bw.JobChan:
//				go func(j *Job) {
//					worker := <-bw.Workers
//					jobChan <- j
//				}(job)
//			case <-bw.Ctx.Ctx.Done():
//				return
//			}
//		}
//	}()
//}
