package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zyberg/ethcli-tx-history/internal/scanner"
)

func main() {
	addressArg := flag.String("address", "", "Ethereum address to scan")
	start := flag.Int64("start", 18000000, "Start block")
	end := flag.Int64("end", 18000100, "End block. Set to 0 to scan all blocks.")
	rpcURL := flag.String("rpc", "https://mainnet.infura.io/v3/YOUR_KEY", "Ethereum RPC URL")
	jsonOutput := flag.Bool("json", false, "Output JSON")
	concurrency := flag.Int("workers", 5, "Concurrency level")
	flag.Parse()

	if *addressArg == "" {
		log.Fatal("Address is required")
	}
	addr := common.HexToAddress(*addressArg)

	client, err := ethclient.Dial(*rpcURL)
	if err != nil {
		log.Fatalf("Ethclient error: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if *end == 0 {
		latest, err := client.BlockNumber(ctx)
		if err != nil {
			log.Fatalf("Failed to fetch latest block: %v", err)
		}
		*end = int64(latest)
	}

	results, err := scanner.ScanNativeTxs(ctx, client, addr, *start, *end, *concurrency)
	if err != nil {
		log.Fatalf("Scan error: %v", err)
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].BlockNumber < results[j].BlockNumber
	})

	// Output
	if *jsonOutput {
		json.NewEncoder(os.Stdout).Encode(results)
		return
	}
	for _, tx := range results {
		fmt.Printf("Block %d | %s | %s | From: %s -> To: %s | %s ETH\n",
			tx.BlockNumber, tx.TxHash, tx.TxType, tx.From, tx.To, tx.Value.Text('f', 6))
	}
}
