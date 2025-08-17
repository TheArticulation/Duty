# x/duty - Hyperlane Validator Duty Enshrinement

A Cosmos SDK module that enshrines Hyperlane validator duties into the consensus validator set, enabling seamless integration between Cosmos SDK staking and Hyperlane cross-chain messaging.

## Overview

The `x/duty` module provides a bridge between Cosmos SDK validators and Hyperlane validators by:

- **Enshrining validator duties** directly in the consensus validator set
- **Exposing validator metadata** (checkpoint signer keys, storage URIs) on-chain
- **Automatically updating** the duty set when validators join/leave the consensus set
- **Providing quorum management** for Hyperlane checkpoint verification
- **Creating deterministic key mappings** between consensus validators and Hyperlane checkpoint signers

This eliminates the need for off-chain governance of Hyperlane validator sets and ensures that the same validators securing your Cosmos chain are also securing your cross-chain messaging.

## Quick Start

### Setting Validator Metadata

```bash
# Set duty metadata for your validator
tx duty set-duty-metadata \
  --signer cosmosvaloper1... \
  --checkpoint-pub-key 0x1234... \
  --checkpoint-storage-uri s3://my-bucket/checkpoints/
```

### Querying the Duty Set

```bash
# Get the complete duty set with all validators and metadata
q duty duty-set

# Get metadata for a specific validator
q duty duty-metadata cosmosvalcons1...
```

### Example Output

```json
{
  "validators": [
    {
      "val_cons_addr": "cosmosvalcons1...",
      "voting_power": "1000000",
      "checkpoint_pub_key": "0x1234...",
      "checkpoint_storage_uri": "s3://my-bucket/checkpoints/"
    }
  ],
  "quorum_num": 2,
  "quorum_den": 3
}
```

## Command Line Interface

The duty module provides a comprehensive CLI for managing validator duties:

### Transaction Commands

```bash
# Set duty metadata for a validator
duty tx set-duty-metadata [signer] [checkpoint-pub-key] [checkpoint-storage-uri]

# Rotate checkpoint signing key
duty tx rotate-checkpoint-key [signer] [new-checkpoint-pub-key] [attestation-signature]

# Bind checkpoint key to consensus validator
duty tx bind-checkpoint-key [signer] [checkpoint-pub-key] [binding-signature] [consensus-address]
```

### Query Commands

```bash
# Query the current duty set
duty query duty-set

# Query metadata for a specific validator
duty query duty-metadata [consensus-address]
```

### Global Flags

```bash
--chain-id string     The network chain ID
--home string         Directory for config and data
--output string       Output format (text|json)
--node string         <host>:<port> to tendermint rpc interface
```

For complete CLI documentation with examples and detailed flag descriptions, see [CLI Documentation](docs/README.md).

## Core Features

- **Validator Metadata Management**: Validators can set their Hyperlane checkpoint signer keys and storage URIs
- **Automatic Duty Set Updates**: The duty set automatically updates when validators join/leave the consensus set
- **Quorum Configuration**: Configurable quorum fractions for Hyperlane checkpoint verification
- **Event Emission**: Events are emitted for validator lifecycle changes and metadata updates
- **gRPC Query Interface**: Clean API for querying duty information
- **Deterministic Key Mapping**: Canonical binding between consensus validators and Hyperlane checkpoint signers
- **Sidecar Integration**: Lightweight service for producing machine-readable Hyperlane manifests

## Integration with Hyperlane

The `x/duty` module is designed to work seamlessly with Hyperlane:

1. **Consensus Validators as Hyperlane Validators**: Cosmos validators can run Hyperlane validators using the same identity
2. **On-Chain Metadata**: Checkpoint signer keys and storage locations are stored on-chain
3. **Quorum Verification**: Relayers and ISMs can query the duty set to verify checkpoint signatures
4. **Automatic Updates**: No manual intervention needed when the validator set changes
5. **Deterministic Attestation**: Canonical binding between consensus validators and Hyperlane checkpoint signers
6. **Sidecar Manifest**: Machine-readable validator set manifest for Hyperlane components

## Documentation

- [Module Overview](docs/overview.md) - Detailed technical documentation
- [CLI Documentation](docs/README.md) - Complete command-line interface reference
- [Sidecar Setup](docs/sidecar.md) - Lightweight service for Hyperlane integration
- [API Reference](docs/api.md) - Complete API documentation
- [Integration Guide](docs/integration.md) - How to integrate with your app

## Development

### Building

```bash
# Build the module
go build ./x/duty

# Run tests
go test ./x/duty/...
```

### Module Structure

```
x/duty/
├── module.go          # Main module implementation
├── keeper/            # State management
│   ├── keeper.go      # Core keeper logic
│   ├── msg_server.go  # Message handlers
│   ├── query_server.go # Query handlers
│   └── hooks.go       # Staking hooks
├── types/             # Type definitions
│   ├── keys.go        # Store keys
│   ├── params.go      # Module parameters
│   ├── msgs.go        # Message types
│   ├── queries.proto  # gRPC query definitions
│   └── codec.go       # Codec registration
└── genesis/           # Genesis state management
    └── genesis.go     # Genesis functions
```

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
