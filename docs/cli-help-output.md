=== Duty Module CLI Help Output ===

=== Main Help ===
```bash
$ duty --help
Duty module CLI

Usage:
  duty [command]

Available Commands:
  tx      Transaction commands for the duty module
  query   Query commands for the duty module (aliases: q)
  help    Help about any command

Flags:
  -h, --help      help for duty
      --version   version for duty

Use "duty [command] --help" for more information about a command.
```

=== Transaction Commands Help ===
```bash
$ duty tx --help
Transaction commands for the duty module

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
      --node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")
```

=== Set Duty Metadata Help ===
```bash
$ duty tx set-duty-metadata --help
Set duty metadata for a validator including checkpoint public key and storage URI

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
      --note string         Note to add a description to the transaction
```

=== Query Commands Help ===
```bash
$ duty query --help
Query commands for the duty module

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
      --node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")
```

=== Duty Set Query Help ===
```bash
$ duty query duty-set --help
Query the current duty set including validators and quorum parameters

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
      --node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")
```

=== Duty Metadata Query Help ===
```bash
$ duty query duty-metadata --help
Query duty metadata for a specific consensus validator

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
      --node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")
```

=== Example Command Execution ===
```bash
$ duty tx set-duty-metadata cosmosvaloper1abc123def456 0x1234567890abcdef s3://my-bucket/checkpoints/ --from my-validator --chain-id duty-testnet-1
Setting duty metadata:
  Signer: cosmosvaloper1abc123def456
  Checkpoint Pub Key: 0x1234567890abcdef
  Checkpoint Storage URI: s3://my-bucket/checkpoints/
```

=== Example Query Execution ===
```bash
$ duty query duty-set --chain-id duty-testnet-1 --output json
{
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
}
```
