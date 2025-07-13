package scanner

import (
	"context"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zyberg/ethcli-tx-history/internal/types"
)

func weiToEth(wei *big.Int) *big.Float {
	f := new(big.Float).SetInt(wei)
	return new(big.Float).Quo(f, big.NewFloat(1e18))
}

func ScanNativeTxs(ctx context.Context, client *ethclient.Client, address common.Address, start, end int64, concurrency int) ([]types.TxRecord, error) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []types.TxRecord

	blockCh := make(chan int64, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for blockNum := range blockCh {
				block, err := client.BlockByNumber(ctx, big.NewInt(blockNum))
				if err != nil {
					log.Printf("Failed to fetch block %d: %v", blockNum, err)
					continue
				}

				for _, tx := range block.Transactions() {
					var signer gethtypes.Signer
					chainID := tx.ChainId()
					if chainID != nil && chainID.Cmp(big.NewInt(0)) != 0 {
						signer = gethtypes.LatestSignerForChainID(chainID)
					} else {
						signer = gethtypes.HomesteadSigner{}
					}
					from, err := signer.Sender(tx)
					if err != nil {
						continue
					}
					to := tx.To()
					value := tx.Value()

					// Skip 0 ETH txs and contract creations
					if value.Cmp(big.NewInt(0)) == 0 || to == nil {
						continue
					}

					txType := ""
					if strings.EqualFold(from.Hex(), address.Hex()) {
						txType = "outgoing"
					} else if strings.EqualFold(to.Hex(), address.Hex()) {
						txType = "incoming"
					} else {
						continue
					}

					mu.Lock()
					results = append(results, types.TxRecord{
						BlockNumber: block.NumberU64(),
						TxHash:      tx.Hash().Hex(),
						TxType:      txType,
						From:        from.Hex(),
						To:          to.Hex(),
						Value:       weiToEth(value),
						Asset:       "ETH",
					})
					mu.Unlock()
				}
			}
		}()
	}

	for i := start; i <= end; i++ {
		blockCh <- i
	}
	close(blockCh)
	wg.Wait()

	return results, nil
}

