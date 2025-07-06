/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
		"context"
    "fmt"
    "log"
    "math/big"
		"time"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/spf13/cobra"
)

var (
    rpcURL string
		startB  uint64
    endB    uint64
		timeout time.Duration
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history [address]",
	Short: "Fetch transaction history for an address",
	Args: cobra.ExactArgs(1),
	Run: runHistory,
}

func runHistory(cmd *cobra.Command, args []string) {
	address := args[0]
	target := common.HexToAddress(address)

	// Timeout context for connection and calls
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Connect to node
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to rpc: %v", err)
	}
	defer client.Close()


	netID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}
	signer := types.NewEIP155Signer(netID)

	// Determine block range
	ctxHdr, cancelHdr := context.WithTimeout(context.Background(), timeout)
	header, err := client.HeaderByNumber(ctxHdr, nil)
	cancelHdr()
	if err != nil {
		log.Fatalf("Failed to get latest block header: %v", err)
	}
	latest := header.Number.Uint64()
	if endB == 0 || endB > latest {
		endB = latest
	}

	fmt.Printf("Scanning blocks %d to %d for %s...\n", startB, endB, target.Hex())

	// Iterate and collect
	for blk := startB; blk <= endB; blk++ {
		fmt.Printf("Fetching block %d\r", blk)
		ctxBlk, cancelBlk := context.WithTimeout(context.Background(), timeout)
		block, err := client.BlockByNumber(ctxBlk, big.NewInt(int64(blk)))
		cancelBlk()
		if err != nil || block == nil {
			log.Printf("Error fetching block %d: %v", blk, err)
			continue
		}

		for _, tx := range block.Transactions() {
			if tx == nil {
				continue
			}
			// Determine sender
			from, err := types.Sender(signer, tx)
			if err != nil {
				continue
			}
			to := tx.To()

			var kind string
			switch {
			case from == target:
				kind = "OUT"
			case to != nil && *to == target:
				kind = "IN"
			default:
				continue
			}

			ethVal := new(big.Float).Quo(new(big.Float).SetInt(tx.Value()), big.NewFloat(1e18))
			fmt.Printf("#%d %s from: %s to: %s %f ETH\n", block.NumberU64(), kind, from.Hex(), to.Hex(), ethVal)
		}
	}
}


func init() {
	rootCmd.AddCommand(historyCmd)

	rootCmd.PersistentFlags().StringVar(&rpcURL, "rpc", "https://mainnet.infura.io/v3/YOUR_KEY", "Ethereum RPC endpoint URL")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 10*time.Second, "RPC call timeout")

	historyCmd.Flags().Uint64Var(&startB, "start", 0, "Start block number")
	historyCmd.Flags().Uint64Var(&endB, "end", 0, "End block number (default latest)")
}
