package sillyhat_pool

import (
	"testing"
	"time"
	"fmt"
	"log"
)



func TestPool(t *testing.T) {
	urls := []string{"11111","22222","33333","444444","55555","66666","77777","88888","999999"}
	pool := new(GoroutinePool)
	pool.Init(3, len(urls))

	for i := range urls {
		url := urls[i]
		pool.AddTask(func() error {
			return download(url)
		})
	}

	isFinish := false

	pool.SetFinishCallback(func() {
		func(isFinish *bool) {
			*isFinish = true
		}(&isFinish)
	})

	pool.Start()

	for !isFinish {
		log.Println("---------- sleep ----------")
		time.Sleep(time.Millisecond * 1000)
	}

	pool.Stop()
	fmt.Println("所有操作已完成！")
}

func TestPool2(t *testing.T) {
	urls := []string{"11111","22222","33333","444444","55555","66666","77777","88888","999999"}
	pool := new(GoroutinePool)
	pool.Init(3, len(urls))

	for i := range urls {
		url := urls[i]
		pool.AddTask(func() error {
			return download(url)
		})
	}

	pool.Start()


	for !pool.IsFinish {
		log.Println("---------- sleep ----------",pool.IsFinish)
		time.Sleep(time.Millisecond * 1000)
	}

	pool.Stop()
	fmt.Println("所有操作已完成！")
}

func download(url string) error {
	fmt.Println("开始下载... ", url)
	time.Sleep(5 * time.Second)
	fmt.Println("## 下载完成！ ", url)
	return nil
}