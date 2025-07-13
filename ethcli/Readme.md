# ethcli

`ethcli` is a small command‑line application that lists all ETH transfers for a given address.
It connects to an Ethereum node over RPC, scans a block range and prints each incoming or outgoing
transaction. Incoming transfers are shown in green and outgoing ones in red for easy reading.

## Requirements

- Go 1.24 or newer
- An Ethereum RPC endpoint (for example from Infura or Ankr)

## Building from source

Run the following command inside this directory to compile the binary:

```bash
go build ./cmd
```

The resulting executable will be named `ethcli`.

## Running the tool

Execute the binary with your desired parameters:

```bash
./ethcli --address 0xDEADBEEF... --start 18000000 --end 18000100 --rpc https://rpc.ankr.com/eth/YOUR_API_KEY
```

Key options:

- `--address` – Ethereum address you want to inspect
- `--start` and `--end` – block range to scan (inclusive)
- `--rpc` – full URL of your Ethereum RPC node
- `--json` – print JSON output instead of the colorized format
- `--workers` – number of concurrent block fetchers (defaults to 5)

## Example output

```
🧱 Block: 18000000
🔗 TxHash:  0x...
📤 Type:    outgoing
📥 From:    0x...
📤 To:      0x...
💰 Value:   1.234567 ETH
────────────────────────────────────────────────────────────────────────────────
```

Each transaction in the specified block range is summarized in a similar way.
