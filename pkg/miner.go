package simplechain;

import (
	"errors";
	"context";
	"bytes";
	"hash";
	"math/rand";
	"sync/atomic";
)

var (
	ErrInterrupted error = errors.New("Mining task interrupted")
)

type mineBlockTask struct {
	ctx context.Context
	block *Block
	result chan []byte
	err chan error
	done uint32
}

type Miner struct {
	ctx context.Context
	ctxCancel context.CancelFunc
	tasks []chan *mineBlockTask
	stop uint32
}

func NewMiner(nWorkers int) *Miner {
	ctx, ctxCancel := context.WithCancel(context.Background())
	m := &Miner{
		ctx: ctx,
		ctxCancel: ctxCancel,
		tasks: make([]chan *mineBlockTask, nWorkers),
	}
	for i := 0; i < nWorkers; i++ {
		m.tasks[i] = make(chan *mineBlockTask)
		go m.worker(m.tasks[i], rand.Int63())
	}
	return m
}

func (this *Miner) Stop() {
	atomic.StoreUint32(&this.stop, 1)
	this.ctxCancel()
}

func (this *Miner) MineBlock(ctx context.Context, block *Block) error {
	task := &mineBlockTask{
		ctx: ctx,
		block: block,
		result: make(chan []byte),
		err: make(chan error),
		done: 0,
	}
	defer func() {
		atomic.StoreUint32(&task.done, 1)
		close(task.result)
		close(task.err)
	}()

	for _, tch := range this.tasks {
		select {
		case tch <- task:
		case <- task.ctx.Done():
			return task.ctx.Err()
		case <- this.ctx.Done():
			return this.ctx.Err()
		}
	}

	select {
	case r := <- task.result:
		block.setNonce(r)
		return nil
	case err := <- task.err:
		return err
	case <- task.ctx.Done():
		return task.ctx.Err()
	case <- this.ctx.Done():
		return this.ctx.Err()
	}
}

func (this *Miner) worker(tasks chan *mineBlockTask, seed int64) {
	generator := rand.New(rand.NewSource(seed))
	for {
		select {
		case task := <- tasks:
			this.mineBlock(task, generator)
		case <- this.ctx.Done():
			return
		}
	}
}

func (this *Miner) mineBlock(task *mineBlockTask, generator *rand.Rand) {
	hash := sha256HashPool.Get().(hash.Hash)
	defer sha256HashPool.Put(hash)
	
	buffer := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buffer)
	buffer.Reset()

	if err := task.block.write(buffer); err != nil {
		this.writeTaskError(task, err)
		return
	}

	blockBytes := buffer.Bytes()
	n := len(blockBytes)

	for (atomic.LoadUint32(&task.done) == 0) && (atomic.LoadUint32(&this.stop) == 0) {
		if _, err := generator.Read(blockBytes[n-NONCE_SIZE:n]); err != nil {
			this.writeTaskError(task, err)
			return
		}

		hash.Reset()
		if _, err := hash.Write(blockBytes); err != nil {
			this.writeTaskError(task, err)
			return
		}
		if this.isValidHashSum(hash.Sum(nil)) {
			this.writeTaskResult(task, blockBytes[n-NONCE_SIZE:n])
			return
		}
	}
	if atomic.LoadUint32(&this.stop) > 0 {
		this.writeTaskError(task, ErrInterrupted)
	}
}

func (this *Miner) writeTaskError(task *mineBlockTask, err error) {
	if atomic.LoadUint32(&task.done) > 0 {
		return
	}
	select {
	case <- this.ctx.Done():
	case <- task.ctx.Done():
	case task.err <- err:
	}
}

func (this *Miner) writeTaskResult(task *mineBlockTask, result []byte) {
	if atomic.LoadUint32(&task.done) > 0 {
		return
	}
	select {
	case <- this.ctx.Done():
	case <- task.ctx.Done():
	case task.result <- result:
	}
}

func (this *Miner) isValidHashSum(hashSum []byte) bool {
	for i := 0; i < LEADING_ZERO_BYTES; i++ {
		if hashSum[i] != 0 {
			return false
		}
	}
	return true
}
