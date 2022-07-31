package pool

import (
	"context"
	"fmt"
	"testing"
	"time"

	"errors"
	"errors"
	
	"github.com/stretchr/testify/assert"
)

type simpleTask struct {
	TaskBase
	ctx    context.Context
	myData int
}

func (p *simpleTask) process() error {
	time.Sleep(time.Second)
	fmt.Println(p.myData)
	return nil
}

func poolSimpleTask(data interface{}) {
	var err error
	defer func() {
		data.(ProcessTasker).SetResult(err)
	}()

	task, ok := data.(*simpleTask)
	if !ok {
		err = fmt.Errorf("data is must poolTask, data=%v", data)
		logs.Log.Error( err)
		return
	}

	select {
	case <-task.ctx.Done():
		err = task.ctx.Err()
	default:
		err = task.process()
	}
}

// TestPoolExecutor 正常 case
func TestPoolExecutor(t *testing.T) {
	taskCount := 20
	ctx := context.Background()
	tasks := make([]interface{}, taskCount)
	for i := 0; i < taskCount; i++ {
		tasks[i] = &simpleTask{
			ctx:    ctx,
			myData: i,
		}
	}
	concurrencyNum := 5
	err := ExecutorList(ctx, tasks, poolSimpleTask, concurrencyNum)
	assert.Nil(t, err)
}

type poolTask struct {
	TaskBase
	ctx    context.Context
	myData int
	result chan *poolTask
}

func (p *poolTask) process() error {
	time.Sleep(time.Duration(p.myData) * time.Millisecond)
	if p.result != nil {
		p.result <- p
	}
	return nil
}

func poolTaskHandle(data interface{}) {
	var err error
	defer func() {
		data.(ProcessTasker).SetResult(err)
	}()

	task, ok := data.(*poolTask)
	if !ok {
		err = fmt.Errorf("data is must poolTask, data=%v", data)
		logs.Log.Error( err)
		return
	}

	select {
	case <-task.ctx.Done():
		err = task.ctx.Err()
	default:
		err = task.process()
	}
}

// TestExecutorRet 带返回值的 case
func TestExecutorRet(t *testing.T) {
	taskCount := 20
	result := make(chan *poolTask, taskCount)
	defer close(result)
	go func() {
		for ret := range result {
			logs.Log.Debug(ret.myData)
		}
	}()

	tasks := make(chan interface{}, taskCount)
	for i := 0; i < taskCount; i++ {
		tasks <- &poolTask{
			ctx:    context.Background(),
			myData: i,
			result: result,
		}
	}
	close(tasks)
	concurrencyNum := 8
	err := Executor(context.Background(), tasks, poolTaskHandle, concurrencyNum)
	assert.Nil(t, err)
}

type poolTaskOther struct {
	TaskBase
	myData int
	result chan *poolTaskOther
}

// TestExecutorHandleTypeError 类型错误 case
func TestExecutorHandleTypeError(t *testing.T) {
	taskCount := 3
	result := make(chan *poolTaskOther, taskCount)
	defer close(result)
	go func() {
		for ret := range result {
			logs.Log.Debug(ret.myData)
		}
	}()

	tasks := make(chan interface{}, taskCount)
	for i := 0; i < taskCount; i++ {
		tasks <- &poolTaskOther{
			myData: i,
			result: result,
		}
	}
	close(tasks)
	concurrencyNum := 2
	err := Executor(context.Background(), tasks, poolTaskHandle, concurrencyNum)
	assert.NotNil(t, err)
}

func poolTaskHandlePanic(data interface{}) {
	var err error
	defer func() {
		data.(ProcessTasker).SetResult(err)
	}()

	task, ok := data.(*poolTask)
	if !ok {
		err = fmt.Errorf("data is must poolTask, data=%v", data)
		logs.Log.Error( err)
		return
	}
	time.Sleep(time.Duration(task.myData) * time.Millisecond)

	errList := []error{nil}
	err = errList[1]
}

func TestPoolExecutorPanic(t *testing.T) {
	taskCount := 3

	tasks := make([]interface{}, taskCount)
	for i := 0; i < taskCount; i++ {
		tasks[i] = &poolTask{
			ctx:    context.Background(),
			myData: i * 10,
		}
	}
	concurrencyNum := 2
	err := ExecutorList(context.Background(), tasks, poolTaskHandlePanic, concurrencyNum)
	assert.Equal(t, errors.RetPanic, errors.Code(err))
}
