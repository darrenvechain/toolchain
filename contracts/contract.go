package contracts

import (
	_ "embed"
	"errors"
	"log/slog"
	"strings"
	"sync"

	"github.com/darrenvechain/thor-go-sdk/thorgo"
	"github.com/darrenvechain/thor-go-sdk/thorgo/accounts"
	"github.com/darrenvechain/thor-go-sdk/txmanager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

//go:embed Toolchain.abi
var ABI string

//go:embed Toolchain.bin
var Bytecode string

func DeployContracts(thor *thorgo.Thor, managers []*txmanager.PKManager, amount int) ([]*accounts.Contract, error) {
	var (
		contracts []*accounts.Contract
		mu        sync.Mutex // mutex to protect concurrent writes
		wg        sync.WaitGroup
	)

	toolchainABI, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {
		return nil, err
	}
	deployer := thor.Deployer(common.Hex2Bytes(Bytecode), &toolchainABI)

	for i := range amount {
		manager := managers[i%len(managers)]
		wg.Add(1)
		go func(m *txmanager.PKManager) {
			defer wg.Done()

			contract, txID, err := deployer.Deploy(manager)
			if err != nil {
				slog.Error("failed to deploy toolchain contract", "error", err, "txID", txID)
				return
			}

			mu.Lock()
			contracts = append(contracts, contract)
			mu.Unlock()

			slog.Info("deployed toolchain contract", "contract", contract.Address.String())
		}(manager)
	}

	wg.Wait()

	if len(contracts) != amount {
		return nil, errors.New("failed to deploy all contracts")
	}

	return contracts, nil
}
