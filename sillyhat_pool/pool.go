package sillyhat_pool

import "fmt"

type GoroutinePool struct {
	Queue  chan func() error
	Number int
	Total  int
	IsFinish       bool

	result         chan error
	finishCallback func()
}

// initial
func (self *GoroutinePool) Init(number int, total int) {
	self.Queue = make(chan func() error, total)
	self.Number = number
	self.Total = total
	self.result = make(chan error, total)
}

// thread start
func (self *GoroutinePool) Start() {
	self.IsFinish = false
	// 开启Number个goroutine
	for i := 0; i < self.Number; i++ {
		go func() {
			for {
				task, ok := <-self.Queue
				if !ok {
					break
				}
				err := task()
				self.result <- err
			}
		}()
	}

	// 获得每个work的执行结果
	for j := 0; j < self.Total; j++ {
		res, ok := <-self.result
		if !ok {
			break
		}

		if res != nil {
			fmt.Println(res)
		}
	}

	// 所有任务都执行完成，回调函数
	if self.finishCallback != nil {
		func(isFinish *bool) {
			*isFinish = true
		}(&self.IsFinish)
		self.finishCallback()
	}
}

// thread stop
func (self *GoroutinePool) Stop() {
	close(self.Queue)
	close(self.result)
}

// add task
func (self *GoroutinePool) AddTask(task func() error) {
	self.Queue <- task
}

// 设置结束回调
func (self *GoroutinePool) SetFinishCallback(callback func()) {
	self.finishCallback = callback
}

//package main
//
//import (
//"fmt"
//"time"
//)
//
//type Pool struct {
//	Queue chan func() error;
//	RuntineNumber int;
//	Total int;
//	Result chan error;
//	FinishCallback func();
//}
//
////初始化
//func (self *Pool) Init(runtineNumber int,total int)  {
//	self.RuntineNumber = runtineNumber;
//	self.Total = total;
//	self.Queue = make(chan func() error, total);
//	self.Result = make(chan error, total);
//}
//
//func (self *Pool) Start()  {
//	//开启 number 个goruntine
//	for i:=0;i<self.RuntineNumber;i++ {
//		go func() {
//			for {
//				task,ok := <-self.Queue
//				if !ok {
//					break;
//				}
//				err := task();
//				self.Result <- err;
//			}
//		}();
//	}
//
//	//获取每个任务的处理结果
//	for j:=0;j<self.RuntineNumber;j++ {
//		res,ok := <-self.Result;
//		if !ok {
//			break;
//		}
//		if res != nil {
//			fmt.Println(res);
//		}
//	}
//
//	//结束回调函数
//	if self.FinishCallback != nil {
//		self.FinishCallback();
//	}
//}
//
////关闭
//func (self *Pool) Stop()  {
//	close(self.Queue);
//	close(self.Result);
//}
//
//func (self *Pool) AddTask(task func() error)  {
//	self.Queue <- task;
//}
//
//func (self *Pool) SetFinishCallback(fun func())  {
//	self.FinishCallback = fun;
//}
//
//func Download(url string) error {
//	time.Sleep(10*time.Second);
//	fmt.Println("Download " + url);
//	return nil;
//}
//
//func DownloadFinish()  {
//	fmt.Println("Download finsh");
//}
//
//func main()  {
//	var p Pool;
//	url := []string{"11111","22222","33333","444444","55555","66666","77777","88888","999999"};
//	p.Init(9, len(url));
//
//	for i := range url {
//		u := url[i];
//		p.AddTask(func() error {
//			return Download(u);
//		});
//	}
//
//	p.SetFinishCallback(DownloadFinish);
//	p.Start();
//	p.Stop();
//}
