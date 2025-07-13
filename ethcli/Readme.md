# ethcli

A command-line tool for scanning an Ethereum address for native (ETH) transfers.

## Build

To build the binary run:

```bash
go build ./cmd
```

## Usage

Run the compiled binary with your desired parameters:

```bash
./ethcli --address 0xDEADBEEF... --start 18000000 --end 18000100 --rpc https://rpc.ankr.com/eth/YOUR_API_KEY
```

- `--address` specifies the target Ethereum address.
- `--start` and `--end` are the starting and ending block numbers to scan.
- `--rpc` must contain a valid Ethereum node RPC URL.

The tool prints the transaction history within the given block range.

