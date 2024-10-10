package main

import (
	"context"
	"flag"
	"log/slog"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/darrenvechain/thor-go-sdk/builtins"
	"github.com/darrenvechain/thor-go-sdk/crypto/hdwallet"
	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/darrenvechain/thor-go-sdk/thorgo"
	"github.com/darrenvechain/thor-go-sdk/txmanager"
	"github.com/darrenvechain/toolchain/contracts"
)

func main() {
	flag.Parse()

	thor, err := thorgo.FromURL(*thorUrlFlag)
	if err != nil {
		slog.Error("failed to create thor client", "error", err)
		os.Exit(1)
	}
	wallet, err := hdwallet.FromMnemonic(*mnemonicFlag)
	if err != nil {
		slog.Error("failed to create wallet from mnemonic", "error", err)
		os.Exit(1)
	}
	managers := createTxManagers(thor, wallet, *requireFundingFlag)

	slog.Info("funded accounts", "count", len(managers))

	toolchainContracts, err := contracts.DeployContracts(thor, managers, 30)
	if err != nil {
		slog.Error("failed to deploy contracts", "error", err)
		os.Exit(1)
	}

	slog.Info("contracts deployed", "count", len(toolchainContracts))

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)

	contracts.Bombard(ctx, managers, toolchainContracts)
}

func createTxManagers(thor *thorgo.Thor, wallet *hdwallet.Wallet, requireFunding bool) []*txmanager.PKManager {
	managers := make([]*txmanager.PKManager, 0)
	vtho := builtins.VTHO.Load(thor)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	fundAmount := big.NewInt(0)
	fundAmount.SetString("1000000000000000000000000", 10)

	for i := range *mnemonicAccountsFlag {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			funder := txmanager.FromPK(wallet.Child(uint32(i)).MustGetPrivateKey(), thor)
			mu.Lock()
			managers = append(managers, funder)
			mu.Unlock()
			clauses := make([]*transaction.Clause, 0)

			for j := 0; j < *extraAccountMultiplier; j++ {
				fundeeIndex := i**extraAccountMultiplier + j
				slog.Info("funding account", "index", fundeeIndex)
				fundee := txmanager.FromPK(wallet.Child(uint32(fundeeIndex)).MustGetPrivateKey(), thor)

				mu.Lock()
				managers = append(managers, fundee)
				mu.Unlock()

				clause, _ := vtho.AsClause("transfer", fundee.Address(), fundAmount)
				clauses = append(clauses, clause)
			}

			if requireFunding {
				tx, err := thor.Transactor(clauses, funder.Address()).Send(funder)

				receipt, err := tx.Wait()
				if err != nil {
					slog.Error("failed to wait for transaction", "error", err, "origin", funder.Address())
					return
				}

				if receipt.Reverted {
					slog.Error("transaction reverted", "txID", tx.ID(), "origin", funder.Address())
					return
				}
			}
		}(i)
	}

	wg.Wait()

	return managers
}
