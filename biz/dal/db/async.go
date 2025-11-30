package db

import (
	"LearnShare/pkg/logger"
	"context"
	"fmt"
	"runtime/debug"
	"sync"
)

// AsyncTask 异步任务结构
type AsyncTask struct {
	Fn  func() error
	Err chan error
}

// AsyncWorkerPool 异步工作池
type AsyncWorkerPool struct {
	taskChan chan AsyncTask
	wg       sync.WaitGroup
	once     sync.Once
}

var (
	globalAsyncPool *AsyncWorkerPool
	poolInitOnce    sync.Once
)

// GetAsyncPool 获取全局异步任务池
func GetAsyncPool() *AsyncWorkerPool {
	poolInitOnce.Do(func() {
		globalAsyncPool = NewAsyncWorkerPool(10) // 10个worker
		globalAsyncPool.Start()
	})
	return globalAsyncPool
}

// NewAsyncWorkerPool 创建异步工作池
func NewAsyncWorkerPool(workerCount int) *AsyncWorkerPool {
	return &AsyncWorkerPool{
		taskChan: make(chan AsyncTask, 100), // 缓冲队列
	}
}

// Start 启动工作池
func (p *AsyncWorkerPool) Start() {
	p.once.Do(func() {
		for i := 0; i < 10; i++ {
			p.wg.Add(1)
			go p.worker()
		}
	})
}

// worker 工作协程
func (p *AsyncWorkerPool) worker() {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("异步工作池 worker panic: %v\nstack: %s", r, string(debug.Stack()))
		}
		p.wg.Done()
	}()

	for task := range p.taskChan {
		func() {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("异步任务 panic: %v", r)
					logger.Errorf("%v\nstack: %s", err, string(debug.Stack()))
					if task.Err != nil {
						task.Err <- err
						close(task.Err)
					}
				}
			}()

			err := task.Fn()
			// 如果是 SubmitNoWait 提交的任务且发生错误，记录日志
			if err != nil && task.Err == nil {
				logger.Errorf("异步任务执行失败（SubmitNoWait）: %v", err)
			}
			if task.Err != nil {
				task.Err <- err
				close(task.Err)
			}
		}()
	}
}

// Submit 提交异步任务
func (p *AsyncWorkerPool) Submit(fn func() error) chan error {
	errChan := make(chan error, 1)
	task := AsyncTask{
		Fn:  fn,
		Err: errChan,
	}
	p.taskChan <- task
	return errChan
}

// SubmitNoWait 提交异步任务(不等待结果)
func (p *AsyncWorkerPool) SubmitNoWait(fn func() error) {
	task := AsyncTask{
		Fn:  fn,
		Err: nil,
	}
	p.taskChan <- task
}

// Shutdown 关闭工作池
func (p *AsyncWorkerPool) Shutdown() {
	close(p.taskChan)
	p.wg.Wait()
}

// AsyncBatch 批量异步执行
func AsyncBatch(ctx context.Context, fns []func() error) []error {
	results := make([]error, len(fns))
	var wg sync.WaitGroup

	for i, fn := range fns {
		wg.Add(1)
		go func(idx int, f func() error) {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("批量异步任务 [%d] panic: %v", idx, r)
					logger.Errorf("%v\nstack: %s", err, string(debug.Stack()))
					results[idx] = err
				}
				wg.Done()
			}()
			results[idx] = f()
		}(i, fn)
	}

	wg.Wait()
	return results
}
