package sillyhat_pool

import (
	"testing"
	"time"
	"fmt"
	"math/rand"
	"golang-cloud/tool/basic"
)

type TestData struct {

	Start int

	End int
}

func TestPool(t *testing.T) {
	pageCount := 50000
	totalRecord := 90000000
	totalPage := (totalRecord+pageCount-1)/pageCount;
	var testDataArray [] TestData
	for i := 0;i < totalPage;i++{
		start := i*pageCount
		end := basic.MinInt((i+1)*pageCount,totalRecord)
		testDataArray = append(testDataArray,*&TestData{Start:start,End:end})
	}
	pool := new(GoroutinePool)
	pool.Init(20, len(testDataArray))
	for i := range testDataArray {
		dataTest := testDataArray[i]
		pool.AddTask(func() error {
			return testInitialArithmeticData(dataTest)
		})
	}
	pool.SetFinishCallback(FinishCallback)
	t1 := time.Now() // get current time
	pool.Start()
	pool.Stop()
	fmt.Println("所有操作已完成！")
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)
}

func testInitialArithmeticData(test TestData) error {
	for i := test.Start; i < test.End; i++ {
		fmt.Println(i)
	}
	return nil
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
	pool.SetFinishCallback(FinishCallback)
	//isFinish := false
	//pool.SetFinishCallback(func() {
	//	func(isFinish *bool) {
	//		*isFinish = true
	//	}(&isFinish)
	//})
	pool.Start()
	//for pool.IsFinish {
	//	log.Println("---------- sleep ----------",pool.IsFinish)
	//	time.Sleep(time.Millisecond * 1000)
	//}
	pool.Stop()
	fmt.Println("所有操作已完成！")
}

func FinishCallback() {
	fmt.Println("---------- finishCallback ----------")
}

func download(url string) error {
	fmt.Println("开始下载... ", url)
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	fmt.Println("## 下载完成！ ", url)
	return nil
}