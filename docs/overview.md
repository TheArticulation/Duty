# x/duty Module Overview

## Introduction

The `x/duty` module enshrines Hyperlane validator duties into the Cosmos SDK consensus validator set. This creates a direct connection between the validators securing your Cosmos chain and the validators securing your cross-chain messaging infrastructure.

### What it Does

In simple terms, the `x/duty` module:

1. **Stores validator metadata** on-chain (checkpoint signer keys, storage URIs)
2. **Exposes the consensus validator set** with their Hyperlane duty information
3. **Automatically updates** when validators join/leave the consensus set
4. **Provides quorum configuration** for Hyperlane checkpoint verification

### How it Works

```
Cosmos SDK Staking Module
         ↓
    Validator Set Changes
         ↓
    x/duty Hooks (Events)
         ↓
    Duty Set Updates
         ↓
    Hyperlane Validators
```

When a validator becomes bonded in the Cosmos SDK staking module, the `x/duty` module automatically includes them in the duty set. Validators can then set their Hyperlane metadata (checkpoint signer keys, storage locations) which becomes part of the on-chain duty set.

## Motivation

### Why Enshrinement?

Traditional Hyperlane validator sets require off-chain governance to manage validator membership. This creates several problems:

- **Coordination overhead**: Manual processes to add/remove validators
- **Security risks**: Off-chain governance can be compromised
- **Operational complexity**: Separate validator sets to manage
- **Inconsistency**: Validators securing consensus may not be securing cross-chain messaging

The `x/duty` module solves these problems by:

- **Automatic updates**: Validator set changes are automatically reflected in the duty set
- **On-chain governance**: All validator metadata is stored and managed on-chain
- **Unified identity**: The same validators secure both consensus and cross-chain messaging
- **Reduced complexity**: Single source of truth for validator information

### Benefits

1. **Security**: Validators with the most stake automatically become Hyperlane validators
2. **Simplicity**: No separate governance processes for Hyperlane validator management
3. **Consistency**: Ensures the same validators securing consensus also secure cross-chain messaging
4. **Transparency**: All validator metadata is publicly visible on-chain

## Core Features

### 1. Validator Metadata Management

Validators can set their Hyperlane duty metadata on-chain:

```go
type DutyMetadata struct {
    CheckpointPubKey     string `json:"checkpoint_pub_key"`     // ECDSA secp256k1 public key
    CheckpointStorageURI string `json:"checkpoint_storage_uri"` // Storage location for signatures
}
```

**Key Features:**
- Only validators can set their own metadata
- Metadata is stored by consensus address (not operator address)
- Automatic validation of address formats and metadata completeness

### 2. Duty Set Queries

The module provides comprehensive querying capabilities:

```bash
# Get complete duty set with all validators and metadata
q duty duty-set

# Get metadata for specific validator
q duty duty-metadata cosmosvalcons1...
```

**Response includes:**
- All bonded validators with their voting power
- Hyperlane metadata (if set)
- Current quorum configuration
- Flattened structure for easy consumption

### 3. Automatic Updates via Hooks

The module integrates with the Cosmos SDK staking module through hooks:

```go
// When validator becomes bonded
func (h DutyHooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, _ sdk.ValAddress) {
    ctx.EventManager().EmitEvent(
        sdk.NewEvent("duty_validator_bonded", sdk.NewAttribute("cons_addr", consAddr.String())),
    )
}

// When validator is removed
func (h DutyHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, _ sdk.ValAddress) {
    ctx.EventManager().EmitEvent(
        sdk.NewEvent("duty_validator_removed", sdk.NewAttribute("cons_addr", consAddr.String())),
    )
}
```

**Benefits:**
- Real-time updates when validator set changes
- Events for off-chain indexers and agents
- No manual intervention required

### 4. Quorum Management

The module manages quorum configuration for Hyperlane checkpoint verification:

```go
type Params struct {
    QuorumNumerator   uint32 `json:"quorum_num"`
    QuorumDenominator uint32 `json:"quorum_den"`
}
```

**Default Configuration:**
- Quorum: 2/3 (66.67%)
- Configurable through governance
- Used by relayers and ISMs for checkpoint verification

## Integration with Hyperlane

### How Consensus Validators Become Hyperlane Validators

1. **Validator Bonding**: When a validator becomes bonded in the Cosmos SDK staking module
2. **Metadata Setting**: Validator sets their Hyperlane metadata on-chain
3. **Duty Set Inclusion**: Validator appears in the duty set with their metadata
4. **Hyperlane Integration**: Validator runs Hyperlane validator software using the same identity

### Checkpoint Verification Flow

```
1. Hyperlane Validator Signs Checkpoint
         ↓
2. Checkpoint Stored at Storage URI
         ↓
3. Relayer Queries Duty Set
         ↓
4. Verifies Signatures Against Duty Set
         ↓
5. Submits to Destination Chain
```

**Key Integration Points:**
- **Checkpoint Signer Keys**: Stored on-chain for verification
- **Storage URIs**: Public locations where signatures are stored
- **Quorum Fraction**: Determines how many signatures are required
- **Validator Set**: Always matches the consensus validator set

### Example Integration

```javascript
// Query duty set
const dutySet = await queryClient.duty.DutySet({});

// Verify checkpoint signatures
const requiredSignatures = Math.ceil(
  dutySet.validators.length * dutySet.quorum_num / dutySet.quorum_den
);

// Check signatures against validator public keys
for (const validator of dutySet.validators) {
  if (validator.checkpoint_pub_key) {
    // Verify signature using validator's public key
    verifySignature(checkpoint, signature, validator.checkpoint_pub_key);
  }
}
```

## Usage

### Setting Validator Metadata

```bash
# Set duty metadata for your validator
tx duty set-duty-metadata \
  --signer cosmosvaloper1... \
  --checkpoint-pub-key 0x1234567890abcdef... \
  --checkpoint-storage-uri s3://my-bucket/checkpoints/
```

**Parameters:**
- `--signer`: Your validator operator address
- `--checkpoint-pub-key`: ECDSA secp256k1 public key for signing checkpoints
- `--checkpoint-storage-uri`: Public storage location for checkpoint signatures

### Querying Duty Information

```bash
# Get complete duty set
q duty duty-set

# Get metadata for specific validator
q duty duty-metadata cosmosvalcons1abcdef...

# Example response
{
  "validators": [
    {
      "val_cons_addr": "cosmosvalcons1abcdef...",
      "voting_power": "1000000",
      "checkpoint_pub_key": "0x1234567890abcdef...",
      "checkpoint_storage_uri": "s3://my-bucket/checkpoints/"
    },
    {
      "val_cons_addr": "cosmosvalcons1fedcba...",
      "voting_power": "500000",
      "checkpoint_pub_key": "0xfedcba0987654321...",
      "checkpoint_storage_uri": "https://my-storage.com/checkpoints/"
    }
  ],
  "quorum_num": 2,
  "quorum_den": 3
}
```

### Events

The module emits events for important state changes:

```json
// Validator bonded
{
  "type": "duty_validator_bonded",
  "attributes": [
    {"key": "cons_addr", "value": "cosmosvalcons1..."}
  ]
}

// Validator removed
{
  "type": "duty_validator_removed",
  "attributes": [
    {"key": "cons_addr", "value": "cosmosvalcons1..."}
  ]
}

// Metadata set
{
  "type": "duty_metadata_set",
  "attributes": [
    {"key": "cons_addr", "value": "cosmosvalcons1..."},
    {"key": "storage_uri", "value": "s3://my-bucket/checkpoints/"}
  ]
}
```

## Running a Hyperlane Validator

### Prerequisites

1. **Cosmos SDK Validator**: You must be a bonded validator in the Cosmos SDK staking module
2. **Hyperlane Validator Software**: Install and configure Hyperlane validator software
3. **Storage Solution**: Set up public storage for checkpoint signatures (S3, IPFS, etc.)

### Configuration Steps

1. **Generate Checkpoint Signer Key**:
   ```bash
   # Generate ECDSA secp256k1 key pair
   openssl ecparam -genkey -name secp256k1 -out checkpoint_key.pem
   openssl ec -in checkpoint_key.pem -pubout -out checkpoint_pub.pem
   ```

2. **Set Up Storage**:
   ```bash
   # Configure S3 bucket for checkpoint storage
   aws s3 mb s3://my-hyperlane-checkpoints
   aws s3api put-bucket-public-access-block \
     --bucket my-hyperlane-checkpoints \
     --public-access-block-configuration "BlockPublicAcls=false,IgnorePublicAcls=false,BlockPublicPolicy=false,RestrictPublicBuckets=false"
   ```

3. **Set Duty Metadata**:
   ```bash
   # Extract public key
   PUBKEY=$(openssl ec -in checkpoint_pub.pem -pubin -text -noout | grep -A 5 "pub:" | tail -n +2 | tr -d ' :\n' | sed 's/^04//')
   
   # Set metadata on-chain
   tx duty set-duty-metadata \
     --signer cosmosvaloper1... \
     --checkpoint-pub-key 0x$PUBKEY \
     --checkpoint-storage-uri s3://my-hyperlane-checkpoints/
   ```

4. **Configure Hyperlane Validator**:
   ```yaml
   # hyperlane-validator.yaml
   validator:
     checkpoint_signer:
       type: "local"
       key: "checkpoint_key.pem"
     storage:
       type: "s3"
       bucket: "my-hyperlane-checkpoints"
       region: "us-east-1"
   ```

### Benefits of Unified Identity

- **Single Key Management**: Use the same validator identity for both consensus and Hyperlane
- **Unified Governance**: Validator set changes automatically propagate to Hyperlane
- **Reduced Complexity**: No separate validator management processes
- **Enhanced Security**: Validators with the most stake automatically secure cross-chain messaging

## For Developers

### Module Structure

```
x/duty/
├── module.go              # Main module implementation
├── keeper/                # State management
│   ├── keeper.go          # Core keeper logic and duty set management
│   ├── msg_server.go      # Message handlers (SetDutyMetadata)
│   ├── query_server.go    # Query handlers (DutySet, DutyMetadata)
│   └── hooks.go           # Staking hooks for automatic updates
├── types/                 # Type definitions
│   ├── keys.go            # Store keys and key generation
│   ├── params.go          # Module parameters (quorum configuration)
│   ├── msgs.go            # Message types and validation
│   ├── queries.proto      # gRPC query service definitions
│   └── codec.go           # Codec registration
└── genesis/               # Genesis state management
    └── genesis.go         # Genesis initialization and export
```

### Key Components

#### Keeper (`keeper/keeper.go`)
- **Duty Metadata CRUD**: `SetDutyMetadata`, `GetDutyMetadata`
- **Duty Set Management**: `GetDutySet` returns validators with metadata
- **Parameter Management**: `GetParams`, `SetParams`
- **Staking Integration**: Interface with staking module for validator information

#### Message Server (`keeper/msg_server.go`)
- **SetDutyMetadata**: Validates and stores validator metadata
- **Authorization**: Only validators can set their own metadata
- **Event Emission**: Emits events for metadata updates

#### Query Server (`keeper/query_server.go`)
- **DutySet**: Returns complete duty set with quorum information
- **DutyMetadata**: Returns metadata for specific validator
- **Error Handling**: Graceful handling of missing metadata

#### Hooks (`keeper/hooks.go`)
- **Validator Lifecycle**: Responds to validator bonding/unbonding
- **Event Emission**: Emits events for off-chain systems
- **Automatic Updates**: No manual intervention required

### Integration into app.go

```go
// In your app.go
import (
    "yourapp/x/duty"
    dutykeeper "yourapp/x/duty/keeper"
    dutymodule "yourapp/x/duty/module"
)

// Add to your app struct
type App struct {
    // ... other fields
    DutyKeeper dutykeeper.Keeper
    DutyModule dutymodule.AppModule
}

// In your app constructor
func NewApp(...) *App {
    // ... other initialization
    
    // Create duty keeper
    app.DutyKeeper = dutykeeper.NewKeeper(
        appCodec,
        keys[dutytypes.StoreKey],
        app.ParamsKeeper.Subspace(dutytypes.ModuleName),
        app.StakingKeeper,
    )
    
    // Create duty module
    app.DutyModule = dutymodule.NewAppModule(app.DutyKeeper)
    
    // Register hooks
    app.StakingKeeper.SetHooks(
        stakingtypes.NewMultiStakingHooks(
            app.DutyKeeper.Hooks(),
            // ... other hooks
        ),
    )
    
    return app
}

// In your module manager
func (app *App) registerAPIRoutes() {
    // ... other routes
    
    // Register duty module
    dutymodule.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCWebRouter)
}
```

### Testing

```bash
# Run all tests
go test ./x/duty/...

# Run specific test
go test ./x/duty/keeper -v

# Run with coverage
go test ./x/duty/... -cover
```

## Future Extensions

### Governance-Controlled Quorum Changes

```go
// Future: Governance proposal for quorum changes
type MsgUpdateQuorum struct {
    Authority string `json:"authority"`
    QuorumNum uint32 `json:"quorum_num"`
    QuorumDen uint32 `json:"quorum_den"`
}
```

### Enhanced Metadata Fields

```go
// Future: Additional metadata fields
type DutyMetadata struct {
    CheckpointPubKey     string `json:"checkpoint_pub_key"`
    CheckpointStorageURI string `json:"checkpoint_storage_uri"`
    StorageProofs        bool   `json:"storage_proofs"`        // Enable storage proofs
    KMSConfig           string `json:"kms_config"`            // KMS configuration
    BackupValidators    []string `json:"backup_validators"`   // Backup validator addresses
}
```

### Indexer/Relayer Integration

- **Event Indexing**: Index duty events for off-chain systems
- **Relayer Integration**: Direct integration with Hyperlane relayers
- **Monitoring**: Dashboard for duty set health and validator participation
- **Alerting**: Notifications for validator metadata issues

### Advanced Features

- **Validator Rotation**: Automatic rotation of checkpoint signer keys
- **Storage Verification**: On-chain verification of storage availability
- **Performance Metrics**: Tracking of validator participation and performance
- **Multi-Chain Support**: Support for multiple Hyperlane domains

## Conclusion

The `x/duty` module provides a robust foundation for integrating Cosmos SDK validators with Hyperlane cross-chain messaging. By enshrining validator duties in the consensus validator set, it eliminates the need for separate governance processes and ensures that the same validators securing your chain are also securing your cross-chain infrastructure.

The module's clean API, automatic updates, and comprehensive querying capabilities make it easy to integrate with existing Hyperlane deployments while providing the security and transparency benefits of on-chain validator management.
