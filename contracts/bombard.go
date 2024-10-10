package contracts

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/darrenvechain/thor-go-sdk/thorgo/accounts"
	"github.com/darrenvechain/thor-go-sdk/txmanager"
	"github.com/darrenvechain/toolchain/random"
)

type bombard struct {
	managers  []*txmanager.PKManager
	contracts []*accounts.Contract
	sendChan  chan struct{}
	ctx       context.Context
	wg        sync.WaitGroup
}

func Bombard(ctx context.Context, managers []*txmanager.PKManager, contracts []*accounts.Contract) {
	b := bombard{
		managers:  managers,
		contracts: contracts,
		sendChan:  make(chan struct{}),
		ctx:       ctx,
	}

	go b.produce()
	go b.consume()

	<-ctx.Done()
	b.wg.Wait()
}

func (b *bombard) produce() {
	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-b.ctx.Done():
			b.wg.Wait()
			return
		case <-ticker.C:
			b.sendChan <- struct{}{}
		}
	}
}

func (b *bombard) consume() {
	for range b.sendChan {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			manager := random.Element(b.managers)
			contract := random.Element(b.contracts)

			clauseAmount := 40
			clauses := make([]*transaction.Clause, clauseAmount)
			for i := 0; i < clauseAmount; i++ {
				a := random.Uint8()
				b := [32]byte(random.Bytes(32))
				c := [32]byte(random.Bytes(32))
				clause, err := contract.AsClause("setBytes32", a, b, c)
				if err != nil {
					slog.Error("failed to create clause", "error", err)
					return
				}
				clauses[i] = clause
			}

			tx, err := manager.SendClauses(clauses)
			if err != nil {
				slog.Error("failed to send transaction", "error", err)
				return
			}

			slog.Info("transaction sent", "txID", tx.String())
		}()
	}
}
