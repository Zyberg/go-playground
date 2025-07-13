package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
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
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(results)
		return
	}
	for _, tx := range results {
		txTypeColor := color.New(color.FgGreen)
		if tx.TxType == "outgoing" {
			txTypeColor = color.New(color.FgRed)
		}

		fmt.Printf("%s Block: %s\n", color.New(color.FgCyan, color.Bold).Sprint("🧱"), color.YellowString("%d", tx.BlockNumber))
		fmt.Printf("🔗 TxHash:  %s\n", color.BlueString("%s", tx.TxHash))
		fmt.Printf("📤 Type:    %s\n", txTypeColor.Sprintf("%s", strings.Title(tx.TxType)))
		fmt.Printf("📥 From:    %s\n", color.HiWhiteString("%s", tx.From))
		fmt.Printf("📤 To:      %s\n", color.HiWhiteString("%s", tx.To))
		fmt.Printf("💰 Value:   %s %s\n", color.HiGreenString(tx.Value.Text('f', 6)), tx.Asset)
		fmt.Println(strings.Repeat("─", 80))
	}
}
