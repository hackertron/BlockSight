package processor

import (
	"log"
	"sync"
	"time"

	"github.com/hackertron/blocksight/internal/models"
)

type BatchProcessor struct {
	batchSize    int
	batchTimeout time.Duration
	processor    func([]*models.Block) error
	batch        []*models.Block
	mu           sync.Mutex
	timer        *time.Timer
}

func NewBatchProcessor(batchSize int, batchTimeout time.Duration, processor func([]*models.Block) error) *BatchProcessor {
	bp := &BatchProcessor{
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
		processor:    processor,
		batch:        make([]*models.Block, 0, batchSize),
	}

	bp.timer = time.NewTimer(batchTimeout)
	go bp.timeoutHandler()

	return bp
}

func (bp *BatchProcessor) Add(block *models.Block) error {
	bp.mu.Lock()
	bp.batch = append(bp.batch, block)

	if len(bp.batch) >= bp.batchSize {
		batch := bp.batch
		bp.batch = make([]*models.Block, 0, bp.batchSize)
		bp.mu.Unlock()
		return bp.processor(batch)
	}

	bp.mu.Unlock()
	return nil
}

func (bp *BatchProcessor) timeoutHandler() {
	for range bp.timer.C {
		bp.mu.Lock()
		if len(bp.batch) > 0 {
			batch := bp.batch
			bp.batch = make([]*models.Block, 0, bp.batchSize)
			bp.mu.Unlock()

			if err := bp.processor(batch); err != nil {
				log.Printf("Batch processing error: %v", err)
			}
		} else {
			bp.mu.Unlock()
		}
		bp.timer.Reset(bp.batchTimeout)
	}
}
