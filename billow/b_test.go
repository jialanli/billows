package billow

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestInitBillow(t *testing.T) {
	bw := InitBillow(3)
	go func() {
		for i := 1; i <= 30; i++ {
			fmt.Println("准备执行：", i)
			bw.jobChan <- &Job{f: fun, Tran: JobTrans{Name: "iii-" + strconv.Itoa(i)}}
		}
	}()
	bw.Start()
	log.Println("*********")
	//time.Sleep(time.Second * 6)
}

func fun() error {
	fmt.Println("哈哈哈")
	time.Sleep(time.Second * 1)
	return nil
}
