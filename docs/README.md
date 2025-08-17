# Duty Module CLI Documentation

The Duty module provides a comprehensive command-line interface for managing Hyperlane validator duties within the Cosmos SDK consensus validator set.

## Overview

The Duty module CLI enables validators to:
- Set and manage duty metadata (checkpoint signing keys and storage URIs)
- Query the current duty set and validator information
- Rotate checkpoint signing keys securely
- Create canonical bindings between consensus and checkpoint keys

## Command Structure

```bash
duty [command] [subcommand] [flags] [args]
```

## Transaction Commands (`tx`)

### Set Duty Metadata

Set duty metadata for a validator including checkpoint public key and storage URI.

```bash
duty tx set-duty-metadata [signer] [checkpoint-pub-key] [checkpoint-storage-uri] [flags]
```

**Arguments:**
- `signer`: Validator operator address (valoper...)
- `checkpoint-pub-key`: ECDSA secp256k1 public key for checkpoint signing (hex format)
- `checkpoint-storage-uri`: Public location for checkpoint signatures (e.g., s3://bucket/prefix/)

**Example:**
```bash
duty tx set-duty-metadata \
  cosmosvaloper1... \
  0x1234567890abcdef \
  s3://my-bucket/hyperlane/checkpoints/ \
  --from my-validator \
  --chain-id duty-testnet-1 \
  --gas auto \
  --gas-adjustment 1.3
```

**Example Output:**
```json
{
  "height": "12345",
  "txhash": "ABC123...",
  "codespace": "",
  "code": 0,
  "data": "0A0A...",
  "raw_log": "[{\"events\":[{\"type\":\"duty_metadata_set\",\"attributes\":[{\"key\":\"cons_addr\",\"value\":\"cosmosvalcons1...\"},{\"key\":\"storage_uri\",\"value\":\"s3://my-bucket/hyperlane/checkpoints/\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "duty_metadata_set",
          "attributes": [
            {
              "key": "cons_addr",
              "value": "cosmosvalcons1..."
            },
            {
              "key": "storage_uri", 
              "value": "s3://my-bucket/hyperlane/checkpoints/"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "150000",
  "tx": null,
  "timestamp": ""
}
```

### Rotate Checkpoint Key

Rotate the checkpoint signing key for a validator with attestation signature.

```bash
duty tx rotate-checkpoint-key [signer] [new-checkpoint-pub-key] [attestation-signature] [flags]
```

**Arguments:**
- `signer`: Validator operator address (valoper...)
- `new-checkpoint-pub-key`: New ECDSA secp256k1 public key for checkpoint signing
- `attestation-signature`: Cryptographic attestation proving key ownership

**Example:**
```bash
duty tx rotate-checkpoint-key \
  cosmosvaloper1... \
  0xabcdef1234567890 \
  0x1f2e3d4c5b6a7980... \
  --from my-validator \
  --chain-id duty-testnet-1
```

**Example Output:**
```json
{
  "height": "12346",
  "txhash": "DEF456...",
  "codespace": "",
  "code": 0,
  "data": "0A0A...",
  "raw_log": "[{\"events\":[{\"type\":\"duty_checkpoint_key_rotated\",\"attributes\":[{\"key\":\"cons_addr\",\"value\":\"cosmosvalcons1...\"},{\"key\":\"new_pub_key\",\"value\":\"0xabcdef1234567890\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "duty_checkpoint_key_rotated",
          "attributes": [
            {
              "key": "cons_addr",
              "value": "cosmosvalcons1..."
            },
            {
              "key": "new_pub_key",
              "value": "0xabcdef1234567890"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "180000",
  "tx": null,
  "timestamp": ""
}
```

### Bind Checkpoint Key

Create a canonical binding between consensus validator and checkpoint key.

```bash
duty tx bind-checkpoint-key [signer] [checkpoint-pub-key] [binding-signature] [consensus-address] [flags]
```

**Arguments:**
- `signer`: Validator operator address (valoper...)
- `checkpoint-pub-key`: ECDSA secp256k1 public key to bind
- `binding-signature`: Cryptographic proof of binding
- `consensus-address`: Consensus validator address (valcons...)

**Example:**
```bash
duty tx bind-checkpoint-key \
  cosmosvaloper1... \
  0x1234567890abcdef \
  0x9e8d7c6b5a493827... \
  cosmosvalcons1... \
  --from my-validator \
  --chain-id duty-testnet-1
```

**Example Output:**
```json
{
  "height": "12347",
  "txhash": "GHI789...",
  "codespace": "",
  "code": 0,
  "data": "0A0A...",
  "raw_log": "[{\"events\":[{\"type\":\"duty_checkpoint_key_bound\",\"attributes\":[{\"key\":\"cons_addr\",\"value\":\"cosmosvalcons1...\"},{\"key\":\"checkpoint_pub_key\",\"value\":\"0x1234567890abcdef\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "duty_checkpoint_key_bound",
          "attributes": [
            {
              "key": "cons_addr",
              "value": "cosmosvalcons1..."
            },
            {
              "key": "checkpoint_pub_key",
              "value": "0x1234567890abcdef"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "160000",
  "tx": null,
  "timestamp": ""
}
```

## Query Commands (`query` or `q`)

### Query Duty Set

Query the current duty set including all validators and quorum parameters.

```bash
duty query duty-set [flags]
```

**Example:**
```bash
duty query duty-set --chain-id duty-testnet-1
```

**Example Output:**
```json
{
  "validators": [
    {
      "val_cons_addr": "cosmosvalcons1abc123def456",
      "voting_power": "1000000",
      "checkpoint_pub_key": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
      "checkpoint_storage_uri": "s3://my-bucket/hyperlane/duty-testnet-1/validators/cosmosvalcons1abc123def456/checkpoints/"
    },
    {
      "val_cons_addr": "cosmosvalcons1ghi789jkl012",
      "voting_power": "800000",
      "checkpoint_pub_key": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
      "checkpoint_storage_uri": "s3://validator2-bucket/hyperlane/duty-testnet-1/validators/cosmosvalcons1ghi789jkl012/checkpoints/"
    },
    {
      "val_cons_addr": "cosmosvalcons1mno345pqr678",
      "voting_power": "600000",
      "checkpoint_pub_key": "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
      "checkpoint_storage_uri": "https://validator3.example.com/hyperlane/checkpoints/"
    }
  ],
  "quorum_num": 2,
  "quorum_den": 3
}
```

### Query Duty Metadata

Query duty metadata for a specific consensus validator.

```bash
duty query duty-metadata [consensus-address] [flags]
```

**Arguments:**
- `consensus-address`: Consensus validator address (valcons...)

**Example:**
```bash
duty query duty-metadata cosmosvalcons1abc123def456 --chain-id duty-testnet-1
```

**Example Output:**
```json
{
  "metadata": {
    "checkpoint_pub_key": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "checkpoint_storage_uri": "s3://my-bucket/hyperlane/duty-testnet-1/validators/cosmosvalcons1abc123def456/checkpoints/"
  }
}
```

**Example Output (No Metadata Found):**
```json
{
  "metadata": null
}
```

## Global Flags

All commands support the following global flags:

```bash
--chain-id string     The network chain ID
--home string         Directory for config and data (default "$HOME/.duty")
--log_level string    Log level (default "info")
--trace              Print out full stack trace on errors
--output string       Output format (text|json) (default "text")
--node string         <host>:<port> to tendermint rpc interface (default "tcp://localhost:26657")
```

## Transaction-Specific Flags

Transaction commands support additional flags:

```bash
--from string         Name or address of private key with which to sign
--fees string         Fees to pay along with transaction (e.g., "10uatom")
--gas string          Gas limit to set per-block (default "200000")
--gas-prices string   Gas prices in decimal format (e.g., "0.1uatom")
--broadcast-mode string Transaction broadcasting mode (sync|async|block) (default "sync")
--yes                 Skip tx broadcasting prompt confirmation
--note string         Note to add a description to the transaction
```

## Examples

### Complete Workflow Example

1. **Set initial duty metadata:**
```bash
duty tx set-duty-metadata \
  cosmosvaloper1abc123def456 \
  0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef \
  s3://my-bucket/hyperlane/duty-testnet-1/validators/cosmosvalcons1abc123def456/checkpoints/ \
  --from my-validator \
  --chain-id duty-testnet-1 \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes
```

2. **Query the duty set to verify:**
```bash
duty query duty-set --chain-id duty-testnet-1 --output json
```

3. **Query specific validator metadata:**
```bash
duty query duty-metadata cosmosvalcons1abc123def456 --chain-id duty-testnet-1 --output json
```

4. **Rotate checkpoint key:**
```bash
duty tx rotate-checkpoint-key \
  cosmosvaloper1abc123def456 \
  0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 \
  0x1f2e3d4c5b6a7980abcdef1234567890abcdef1234567890abcdef1234567890 \
  --from my-validator \
  --chain-id duty-testnet-1 \
  --yes
```

5. **Bind checkpoint key to consensus validator:**
```bash
duty tx bind-checkpoint-key \
  cosmosvaloper1abc123def456 \
  0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef \
  0x9e8d7c6b5a493827fedcba0987654321fedcba0987654321fedcba0987654321 \
  cosmosvalcons1abc123def456 \
  --from my-validator \
  --chain-id duty-testnet-1 \
  --yes
```

## Error Handling

The CLI provides clear error messages for common issues:

- **Invalid validator address**: Returns validation error with format details
- **Missing metadata**: Returns empty response for queries
- **Insufficient permissions**: Returns authorization error
- **Invalid signatures**: Returns cryptographic validation error
- **Network issues**: Returns connection error with retry suggestions

## Integration with Hyperlane

The duty module CLI integrates seamlessly with Hyperlane validator operations:

1. **Validator Setup**: Use `set-duty-metadata` to configure checkpoint signing keys
2. **Key Rotation**: Use `rotate-checkpoint-key` for secure key updates
3. **Binding Verification**: Use `bind-checkpoint-key` for canonical validator-key mappings
4. **Monitoring**: Use query commands to monitor validator set changes
5. **Relayer Integration**: Query duty set for checkpoint verification parameters

## Security Considerations

- Always verify transaction details before signing
- Use secure key management for validator keys
- Regularly rotate checkpoint signing keys
- Monitor duty set changes for unexpected validator modifications
- Validate all cryptographic signatures and attestations
