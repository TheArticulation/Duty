#!/bin/bash

# Script to generate CLI help output for documentation

echo "=== Duty Module CLI Help Output ==="
echo

echo "=== Main Help ==="
echo '```bash'
echo '$ duty --help'
echo 'Duty module CLI

Usage:
  duty [command]

Available Commands:
  tx      Transaction commands for the duty module
  query   Query commands for the duty module (aliases: q)
  help    Help about any command

Flags:
  -h, --help      help for duty
      --version   version for duty

Use "duty [command] --help" for more information about a command.'
echo '```'
echo

echo "=== Transaction Commands Help ==="
echo '```bash'
echo '$ duty tx --help'
echo 'Transaction commands for the duty module

Usage:
  duty tx [command]

Available Commands:
  set-duty-metadata     Set duty metadata for a validator
  rotate-checkpoint-key Rotate checkpoint signing key for a validator
  bind-checkpoint-key   Bind checkpoint key to consensus validator

Flags:
  -h, --help   help for tx

Global Flags:
      --chain-id string     The network chain ID
      --home string         Directory for config and data (default "$HOME/.duty")
      --log_level string    Log level (default "info")
      --trace              Print out full stack trace on errors
      --output string       Output format (text|json) (default "text")
      --node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")'
echo '```'
echo

echo "=== Set Duty Metadata Help ==="
echo '```bash'
echo '$ duty tx set-duty-metadata --help'
echo 'Set duty metadata for a validator including checkpoint public key and storage URI

Usage:
  duty tx set-duty-metadata [signer] [checkpoint-pub-key] [checkpoint-storage-uri] [flags]

Arguments:
  signer                    Validator operator address (valoper...)
  checkpoint-pub-key        ECDSA secp256k1 public key for checkpoint signing (hex format)
  checkpoint-storage-uri    Public location for checkpoint signatures (e.g., s3://bucket/prefix/)

Flags:
  -h, --help   help for set-duty-metadata

Global Flags:
      --chain-id string     The network chain ID
      --home string         Directory for config and data (default "$HOME/.duty")
      --log_level string    Log level (default "info")
      --trace              Print out full stack trace on errors
      --output string       Output format (text|json) (default "text")
      --node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")

Transaction Flags:
      --from string         Name or address of private key with which to sign
      --fees string         Fees to pay along with transaction (e.g., "10uatom")
      --gas string          Gas limit to set per-block (default "200000")
      --gas-prices string   Gas prices in decimal format (e.g., "0.1uatom")
      --broadcast-mode string Transaction broadcasting mode (sync|async|block) (default "sync")
      --yes                 Skip tx broadcasting prompt confirmation
      --note string         Note to add a description to the transaction'
echo '```'
echo

echo "=== Query Commands Help ==="
echo '```bash'
echo '$ duty query --help'
echo 'Query commands for the duty module

Usage:
  duty query [command]

Aliases:
  duty q

Available Commands:
  duty-set      Query the current duty set
  duty-metadata Query duty metadata for a validator

Flags:
  -h, --help   help for query

Global Flags:
      --chain-id string     The network chain ID
      --home string         Directory for config and data (default "$HOME/.duty")
      --log_level string    Log level (default "info")
      --trace              Print out full stack trace on errors
      --output string       Output format (text|json) (default "text")
      --node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")'
echo '```'
echo

echo "=== Duty Set Query Help ==="
echo '```bash'
echo '$ duty query duty-set --help'
echo 'Query the current duty set including validators and quorum parameters

Usage:
  duty query duty-set [flags]

Flags:
  -h, --help   help for duty-set

Global Flags:
      --chain-id string     The network chain ID
      --home string         Directory for config and data (default "$HOME/.duty")
      --log_level string    Log level (default "info")
      --trace              Print out full stack trace on errors
      --output string       Output format (text|json) (default "text")
      --node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")'
echo '```'
echo

echo "=== Duty Metadata Query Help ==="
echo '```bash'
echo '$ duty query duty-metadata --help'
echo 'Query duty metadata for a specific consensus validator

Usage:
  duty query duty-metadata [consensus-address] [flags]

Arguments:
  consensus-address    Consensus validator address (valcons...)

Flags:
  -h, --help   help for duty-metadata

Global Flags:
      --chain-id string     The network chain ID
      --home string         Directory for config and data (default "$HOME/.duty")
      --log_level string    Log level (default "info")
      --trace              Print out full stack trace on errors
      --output string       Output format (text|json) (default "text")
      --node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")'
echo '```'
echo

echo "=== Example Command Execution ==="
echo '```bash'
echo '$ duty tx set-duty-metadata cosmosvaloper1abc123def456 0x1234567890abcdef s3://my-bucket/checkpoints/ --from my-validator --chain-id duty-testnet-1'
echo 'Setting duty metadata:
  Signer: cosmosvaloper1abc123def456
  Checkpoint Pub Key: 0x1234567890abcdef
  Checkpoint Storage URI: s3://my-bucket/checkpoints/'
echo '```'
echo

echo "=== Example Query Execution ==="
echo '```bash'
echo '$ duty query duty-set --chain-id duty-testnet-1 --output json'
echo '{
  "validators": [
    {
      "val_cons_addr": "cosmosvalcons1abc123def456",
      "voting_power": "1000000",
      "checkpoint_pub_key": "0x1234567890abcdef",
      "checkpoint_storage_uri": "s3://my-bucket/checkpoints/"
    }
  ],
  "quorum_num": 2,
  "quorum_den": 3
}'
echo '```'
