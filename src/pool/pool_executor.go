// Package pool 对协程池进行的简单封装
package pool

import (
	"context"
	"runtime"
	"time"

	"errors"
	"errors"
	
	"github.com/panjf2000/ants/v2"
)

// ProcessTasker 并发操作结构。调用者只需包含下边的 TaskBase 即可
type ProcessTasker interface {
	setErrorChan(result chan error) // 设置错误输出管道
	SetResult(err error)            // 将错误塞入管道，通知调用者
}

// TaskBase 实现 ProcessTasker 接口，调用者包含此结构，将错误执行结果通过 SetResult 输出到 ErrChan 中
type TaskBase struct {
	errChan chan error
}

func (t *TaskBase) setErrorChan(result chan error) {
	t.errChan = result
}

// SetResult 设置错误结果到输出通道中，用于告知调用者
func (t *TaskBase) SetResult(err error) {
	if err != nil {
		t.errChan <- err
	}
}

// ExecutorList 封装 PoolExecutor 函数，方便list调用
func ExecutorList(ctx context.Context, tasks []interface{}, handle func(interface{}), concurrencyNum int) error {
	if len(tasks) == 0 {
		return nil
	}
	tasksChan := make(chan interface{}, len(tasks))
	for i := 0; i < len(tasks); i++ {
		tasksChan <- tasks[i]
	}
	close(tasksChan)
	return Executor(context.Background(), tasksChan, handle, concurrencyNum)
}

// processTask 正在调度协程进行任务处理
func processTask(
	ctx context.Context,
	pool *ants.PoolWithFunc,
	tasks chan interface{},
	taskErrChan chan error,
	chanClose chan struct{},
) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasks:
			if !ok { // chan is close
				time.Sleep(10 * time.Millisecond)
				chanClose <- struct{}{}
				logs.Log.Debugf("tasks loop out, count:%v", pool.Running())
				return
			}
			task.(ProcessTasker).setErrorChan(taskErrChan)

			if err := pool.Invoke(task); err != nil {
				logs.Log.Errorf("pool.Invoke fail, count:%v, err:%+v", pool.Running(), err)
				taskErrChan <- err
				return
			}
		}
	}
}

func newPool(
	ctx context.Context,
	handle func(interface{}),
	concurrencyNum int,
	errPanicChan chan error,
) (*ants.PoolWithFunc, error) {
	return ants.NewPoolWithFunc(
		concurrencyNum,
		handle,
		ants.WithPanicHandler(func(err interface{}) {
			logs.Log.Errorf("并发触发 panic %+v", err)
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			logs.Log.Errorf("调度栈: %v", string(buf[:n]))
			errPanicChan <- errors.New(fmt.Sprintf("pool panic %+v", err))
		}),
	)
}

// Executor 限制 goroutine 个数的并发处理, tasks 结构需要内含 struct TaskBase, 参考如下
// type taskData struct {
// 	  TaskBase
//    otherData int
// }
// func handle(data interface{}){
// 	   var err error  //handle执行结果
// 	   defer func() {
// 	   	    data.(ProcessTasker).SetResult(err)
// 	   }()
//     task, ok := data.(*taskData)
//     if !ok {
//     	   err = fmt.Errorf("data is must taskData, data=%v", data)
//     	   ErrorContext(context.Background(), err)
//     	   return
//     }
// 	   fmt.Println(task.otherData)
// }
func Executor(ctx context.Context, tasks chan interface{}, handle func(interface{}), concurrencyNum int) error {
	var (
		taskErrChan  = make(chan error, concurrencyNum)
		chanClose    = make(chan struct{}, 1)
		errPanicChan = make(chan error)
	)

	pool, err := newPool(ctx, handle, concurrencyNum, errPanicChan)
	if err != nil {
		return err
	}
	defer pool.Release()

	go processTask(ctx, pool, tasks, taskErrChan, chanClose)

	for {
		select {
		case <-ctx.Done(): // 调用者主动触发结束
			return ctx.Err()
		case err := <-errPanicChan:
			logs.Log.Errorf("task process panic, count:%v, err:%+v", pool.Running(), err)
			return err
		case err := <-taskErrChan:
			if err != nil { // 任务执行异常结束
				logs.Log.Errorf("task process fail, count:%v, err:%+v", pool.Running(), err)
				return err
			}
		case data := <-chanClose: // 输入通道关闭，等待任务执行结束
			if pool.Running() == 0 {
				return getCloseErr(ctx, errPanicChan, taskErrChan)
			}
			chanClose <- data // 防止taskErrChan没有数据，陷入死循环
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// getCloseErr 获取错误，特别是对panic的处理
func getCloseErr(ctx context.Context, errPanicChan chan error, taskErrChan chan error) error {
	var err error
	if len(taskErrChan) > 0 {
		logs.Log.Errorf("task process fail, err:%+v", err)
		err = <-taskErrChan
	}
	if len(errPanicChan) > 0 {
		logs.Log.Errorf("task process panic, err:%+v", err)
		err = <-errPanicChan
	}
	if err != nil {
		logs.Log.Errorf("task process over, err:%+v", err)
	}
	return err
}
