// Copyright 2022 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package workerpool

import (
	"context"
	"log"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task := NewTask(ctx, func(c context.Context) (interface{}, error) {
		if c != ctx {
			t.Fatal()
		}
		return nil, nil
	})

	// execute task
	task.Execute()

	if r := <-task.Result(); r.Err != nil || r.Result != nil {
		t.Fatal()
	}
}

func TestNewPool(t *testing.T) {
	var ctx context.Context
	pool := NewPool(ctx, Option{ExpandableLimit: -1})
	if pool.opt.NumberWorker != numCPU || pool.opt.ExpandableLimit != 0 || pool.opt.ExpandedLifetime != time.Minute {
		t.Fatal()
	}
}

//nolint
func TestPool(t *testing.T) {
	var taskCtx context.Context // nil context

	pool := NewPool(context.Background(), Option{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tasks := make([]*Task, 1024)
		for i := range tasks {
			if i&1 == 0 {
				tasks[i] = pool.Execute(func(context.Context) (interface{}, error) {
					time.Sleep(2 * time.Millisecond)
					return nil, nil
				})
				pool.TryExecute(func(context.Context) (interface{}, error) {
					time.Sleep(time.Millisecond)
					return nil, nil
				})
			} else {
				tasks[i] = pool.ExecuteWithCtx(taskCtx, func(context.Context) (interface{}, error) {
					time.Sleep(2 * time.Millisecond)
					return nil, nil
				})
				pool.TryExecuteWithCtx(taskCtx, func(context.Context) (interface{}, error) {
					time.Sleep(time.Millisecond)
					return nil, nil
				})
			}
		}
	}()
	wg.Wait()
	pool.Stop()
}

//nolint
func TestTaskCtxCanceled(t *testing.T) {
	pool := NewPool(context.Background(), Option{})

	{
		task := NewTask(context.Background(), func(c context.Context) (interface{}, error) {
			return nil, nil
		})
		pool.Do(task)
	}
	time.Sleep(5 * time.Millisecond)

	{
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		for i := 0; i < 1000; i++ {
			task := NewTask(ctx, func(c context.Context) (interface{}, error) {
				time.Sleep(2 * time.Millisecond)
				return nil, nil
			})
			pool.Do(task)
		}
	}

	{
		oldctx := pool.ctx

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		pool.ctx = ctx

		for i := 0; i < 1000; i++ {
			task := NewTask(nil, func(c context.Context) (interface{}, error) {
				time.Sleep(2 * time.Millisecond)
				return nil, nil
			})
			pool.Do(task)
		}

		pool.ctx = oldctx
	}

	{
		task := NewTask(nil, func(c context.Context) (interface{}, error) {
			return nil, nil
		})
		pool.Do(task)
	}

	{
		task := NewTask(nil, func(c context.Context) (interface{}, error) {
			return nil, nil
		})
		pool.TryDo(task)
	}

	{
		oldctx := pool.ctx

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		pool.ctx = ctx

		for i := 0; i < 1000; i++ {
			task := NewTask(nil, func(c context.Context) (interface{}, error) {
				time.Sleep(2 * time.Millisecond)
				return nil, nil
			})
			pool.TryDo(task)
		}

		pool.ctx = oldctx
	}
	time.Sleep(5 * time.Millisecond)

	{
		task := NewTask(context.Background(), func(c context.Context) (interface{}, error) {
			return nil, nil
		})
		pool.TryDo(task)
	}
	time.Sleep(5 * time.Millisecond)

	{
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		task := NewTask(ctx, func(c context.Context) (interface{}, error) {
			return nil, nil
		})
		pool.Do(task)
	}
	time.Sleep(5 * time.Millisecond)

	{
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		task := NewTask(ctx, func(c context.Context) (interface{}, error) {
			return nil, nil
		})
		pool.TryDo(task)
	}
	time.Sleep(5 * time.Millisecond)

	pool.cancel()
	{
		task := NewTask(context.Background(), func(c context.Context) (interface{}, error) {
			return nil, nil
		})
		pool.Do(task)
	}
	{
		task := NewTask(context.Background(), func(c context.Context) (interface{}, error) {
			return nil, nil
		})
		pool.TryDo(task)
	}

	pool.Stop()
}

func TestPoolWithExpandable(t *testing.T) {
	var taskCtx context.Context // nil context

	exec := func(shouldSleepBeforeClose bool) {
		pool := NewPool(context.Background(), Option{ExpandableLimit: 2, ExpandedLifetime: 10 * time.Millisecond})
		pool.Start()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			tasks := make([]*Task, 256)
			for i := range tasks {
				if i&1 == 0 {
					tasks[i] = pool.Execute(func(context.Context) (interface{}, error) {
						time.Sleep(5 * time.Millisecond)
						return nil, nil
					})
					pool.TryExecute(func(context.Context) (interface{}, error) {
						time.Sleep(5 * time.Millisecond)
						return nil, nil
					})
				} else {
					tasks[i] = pool.ExecuteWithCtx(context.Background(), func(context.Context) (interface{}, error) {
						time.Sleep(5 * time.Millisecond)
						return nil, nil
					})
					pool.TryExecuteWithCtx(taskCtx, func(context.Context) (interface{}, error) {
						time.Sleep(5 * time.Millisecond)
						return nil, nil
					})
				}
			}
		}()
		wg.Wait()

		if shouldSleepBeforeClose {
			time.Sleep(200 * time.Millisecond)
		}

		pool.Stop()
	}

	t.Run("HappyCase", func(*testing.T) {
		exec(false)
	})

	t.Run("Expiration", func(*testing.T) {
		exec(true)
	})
}

//nolint
func TestCorrectness(t *testing.T) {
	pool := NewPool(context.Background(), Option{NumberWorker: runtime.NumCPU()})
	pool.Start()

	// Calculate (1^1 + 2^2 + 3^3 + ... + 1000000^1000000) modulo 1234567
	tasks := make([]*Task, 0, 1000000)
	for i := 1; i <= 1000000; i++ {
		task := moduloTask(context.Background(), uint(i), uint(i), 1234567)
		pool.Do(task)
		tasks = append(tasks, task)
	}

	// collect task results
	var s1, s2 uint
	for i := range tasks {
		if result := <-tasks[i].Result(); result.Err != nil {
			log.Fatal(result.Err)
		} else {
			s1 = uint((uint64(s1) + uint64(result.Result.(uint))) % 1234567)
		}
	}

	// sequential computation
	for i := 1; i <= 1000000; i++ {
		s2 = uint((uint64(s2) + uint64(modulo(uint(i), uint(i), 1234567))) % 1234567)
	}
	if s1 != s2 {
		log.Fatal(s1, s2)
	}

	pool.Stop()
}

func moduloTask(ctx context.Context, a, b, N uint) *Task {
	return NewTask(ctx, func(ctx context.Context) (interface{}, error) {
		return modulo(a, b, N), nil
	})
}

// calculate a^b MODULO N
func modulo(a, b uint, N uint) uint {
	switch b {
	case 0:
		return 1 % N
	case 1:
		return a % N
	default:
		if b&1 == 0 {
			t := modulo(a, b>>1, N)
			return uint(uint64(t) * uint64(t) % uint64(N))
		}

		t := modulo(a, b>>1, N)
		t = uint(uint64(t) * uint64(t) % uint64(N))
		return uint(uint64(a) * uint64(t) % uint64(N))
	}
}

func TestStopTimer(t *testing.T) {
	ti := time.NewTimer(time.Millisecond)
	time.Sleep(2 * time.Millisecond)
	stopTimer(ti)
}
