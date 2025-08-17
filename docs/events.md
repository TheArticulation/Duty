# Duty Module Events

The Duty module emits a comprehensive set of events that enable off-chain indexers, sidecars, and monitoring systems to track validator duty changes in real-time.

## Event Overview

Events are emitted for all major state changes in the duty module:

- **Validator lifecycle events**: When validators join, leave, or begin unbonding
- **Metadata management events**: When validators set or update their duty metadata
- **Key management events**: When checkpoint keys are rotated or bound

## Event Types

### 1. Validator Lifecycle Events

#### `duty_validator_bonded`

Emitted when a validator joins the active validator set.

**Attributes:**
- `cons_addr`: Consensus validator address (bech32)
- `val_addr`: Validator operator address (bech32)
- `voting_power`: Validator's voting power (string)
- `moniker`: Validator's moniker/name

**Example:**
```json
{
  "type": "duty_validator_bonded",
  "attributes": [
    {
      "key": "cons_addr",
      "value": "cosmosvalcons1abc123def456"
    },
    {
      "key": "val_addr", 
      "value": "cosmosvaloper1abc123def456"
    },
    {
      "key": "voting_power",
      "value": "1000000"
    },
    {
      "key": "moniker",
      "value": "My Validator"
    }
  ]
}
```

#### `duty_validator_removed`

Emitted when a validator is removed from the active validator set.

**Attributes:**
- `cons_addr`: Consensus validator address (bech32)
- `val_addr`: Validator operator address (bech32)

**Example:**
```json
{
  "type": "duty_validator_removed",
  "attributes": [
    {
      "key": "cons_addr",
      "value": "cosmosvalcons1abc123def456"
    },
    {
      "key": "val_addr",
      "value": "cosmosvaloper1abc123def456"
    }
  ]
}
```

#### `duty_validator_unbonding`

Emitted when a validator begins the unbonding process.

**Attributes:**
- `cons_addr`: Consensus validator address (bech32)
- `val_addr`: Validator operator address (bech32)

**Example:**
```json
{
  "type": "duty_validator_unbonding",
  "attributes": [
    {
      "key": "cons_addr",
      "value": "cosmosvalcons1abc123def456"
    },
    {
      "key": "val_addr",
      "value": "cosmosvaloper1abc123def456"
    }
  ]
}
```

### 2. Metadata Management Events

#### `duty_metadata_set`

Emitted when a validator sets or updates their duty metadata.

**Attributes:**
- `cons_addr`: Consensus validator address (bech32)
- `val_addr`: Validator operator address (bech32)
- `checkpoint_pub_key`: ECDSA secp256k1 public key for checkpoint signing (hex)
- `storage_uri`: Public location for checkpoint signatures
- `block_height`: Block height when the event was emitted

**Example:**
```json
{
  "type": "duty_metadata_set",
  "attributes": [
    {
      "key": "cons_addr",
      "value": "cosmosvalcons1abc123def456"
    },
    {
      "key": "val_addr",
      "value": "cosmosvaloper1abc123def456"
    },
    {
      "key": "checkpoint_pub_key",
      "value": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    },
    {
      "key": "storage_uri",
      "value": "s3://my-bucket/hyperlane/duty-testnet-1/validators/cosmosvalcons1abc123def456/checkpoints/"
    },
    {
      "key": "block_height",
      "value": "12345"
    }
  ]
}
```

### 3. Key Management Events

#### `duty_checkpoint_key_rotated`

Emitted when a validator rotates their checkpoint signing key.

**Attributes:**
- `cons_addr`: Consensus validator address (bech32)
- `val_addr`: Validator operator address (bech32)
- `old_checkpoint_pub_key`: Previous ECDSA secp256k1 public key (hex)
- `new_checkpoint_pub_key`: New ECDSA secp256k1 public key (hex)
- `block_height`: Block height when the event was emitted

**Example:**
```json
{
  "type": "duty_checkpoint_key_rotated",
  "attributes": [
    {
      "key": "cons_addr",
      "value": "cosmosvalcons1abc123def456"
    },
    {
      "key": "val_addr",
      "value": "cosmosvaloper1abc123def456"
    },
    {
      "key": "old_checkpoint_pub_key",
      "value": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    },
    {
      "key": "new_checkpoint_pub_key",
      "value": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
    },
    {
      "key": "block_height",
      "value": "12346"
    }
  ]
}
```

#### `duty_checkpoint_key_bound`

Emitted when a checkpoint key is bound to a consensus validator.

**Attributes:**
- `cons_addr`: Consensus validator address (bech32)
- `val_addr`: Validator operator address (bech32)
- `checkpoint_pub_key`: ECDSA secp256k1 public key being bound (hex)
- `binding_signature`: Cryptographic proof of binding (hex)
- `block_height`: Block height when the event was emitted

**Example:**
```json
{
  "type": "duty_checkpoint_key_bound",
  "attributes": [
    {
      "key": "cons_addr",
      "value": "cosmosvalcons1abc123def456"
    },
    {
      "key": "val_addr",
      "value": "cosmosvaloper1abc123def456"
    },
    {
      "key": "checkpoint_pub_key",
      "value": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    },
    {
      "key": "binding_signature",
      "value": "0x9e8d7c6b5a493827fedcba0987654321fedcba0987654321fedcba0987654321"
    },
    {
      "key": "block_height",
      "value": "12347"
    }
  ]
}
```

## Event Indexing and Monitoring

### Real-time Event Processing

Events are emitted immediately when state changes occur, enabling real-time monitoring:

```bash
# Subscribe to duty events via Tendermint RPC
curl -X GET "http://localhost:26657/subscribe?query=duty_validator_bonded"

# Subscribe to all duty events
curl -X GET "http://localhost:26657/subscribe?query=duty"
```

### Event Filtering

Indexers can filter events by type and attributes:

```bash
# Filter by event type
query="tm.event='Tx' AND duty_validator_bonded"

# Filter by validator address
query="tm.event='Tx' AND duty_metadata_set.cons_addr='cosmosvalcons1abc123def456'"

# Filter by block height range
query="tm.event='Tx' AND duty_metadata_set.block_height>'12340'"
```

## Sidecar/Indexer Integration

### Event-Driven Manifest Updates

Sidecars and indexers can use these events to maintain up-to-date Hyperlane manifests:

1. **Subscribe to Events**: Listen for all duty-related events
2. **Query Current State**: Use `/duty.DutySet` to get the complete validator set
3. **Update Manifest**: Generate new manifest when events are received
4. **Publish Changes**: Make updated manifest available to Hyperlane components

### Example Event Processing Flow

```python
import asyncio
import json
from typing import Dict, List

class DutyEventProcessor:
    def __init__(self, rpc_url: str, grpc_url: str):
        self.rpc_url = rpc_url
        self.grpc_url = grpc_url
        self.manifest = {}
    
    async def process_event(self, event: Dict):
        """Process a duty event and update the manifest"""
        event_type = event["type"]
        
        if event_type == "duty_validator_bonded":
            await self.handle_validator_bonded(event)
        elif event_type == "duty_validator_removed":
            await self.handle_validator_removed(event)
        elif event_type == "duty_metadata_set":
            await self.handle_metadata_set(event)
        elif event_type == "duty_checkpoint_key_rotated":
            await self.handle_key_rotation(event)
        
        # Update manifest after processing
        await self.update_manifest()
    
    async def handle_validator_bonded(self, event: Dict):
        """Handle validator joining the set"""
        cons_addr = self.get_attribute(event, "cons_addr")
        val_addr = self.get_attribute(event, "val_addr")
        voting_power = self.get_attribute(event, "voting_power")
        
        # Add validator to manifest
        self.manifest["validators"][cons_addr] = {
            "val_addr": val_addr,
            "voting_power": voting_power,
            "checkpoint_pub_key": None,
            "checkpoint_storage_uri": None
        }
    
    async def handle_metadata_set(self, event: Dict):
        """Handle metadata updates"""
        cons_addr = self.get_attribute(event, "cons_addr")
        checkpoint_pub_key = self.get_attribute(event, "checkpoint_pub_key")
        storage_uri = self.get_attribute(event, "storage_uri")
        
        if cons_addr in self.manifest["validators"]:
            self.manifest["validators"][cons_addr].update({
                "checkpoint_pub_key": checkpoint_pub_key,
                "checkpoint_storage_uri": storage_uri
            })
    
    async def update_manifest(self):
        """Query current state and update manifest"""
        # Query current duty set for complete state
        duty_set = await self.query_duty_set()
        
        # Update manifest with current state
        self.manifest.update({
            "chain_id": duty_set["chain_id"],
            "quorum": {
                "num": duty_set["quorum_num"],
                "den": duty_set["quorum_den"]
            },
            "validators": duty_set["validators"],
            "asof_height": duty_set["block_height"],
            "version": 1
        })
        
        # Publish updated manifest
        await self.publish_manifest()
    
    def get_attribute(self, event: Dict, key: str) -> str:
        """Extract attribute value from event"""
        for attr in event["attributes"]:
            if attr["key"] == key:
                return attr["value"]
        return None
```

### Manifest Generation

The sidecar generates a canonical manifest from events and queries:

```json
{
  "chain_id": "duty-testnet-1",
  "quorum": {
    "num": 2,
    "den": 3
  },
  "validators": [
    {
      "consensus_address": "cosmosvalcons1abc123def456",
      "voting_power": "1000000",
      "checkpoint_pub_key": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
      "checkpoint_storage_uri": "s3://my-bucket/hyperlane/duty-testnet-1/validators/cosmosvalcons1abc123def456/checkpoints/"
    }
  ],
  "version": 1,
  "asof_height": 12345,
  "last_updated": "2024-01-15T10:30:00Z"
}
```

## Best Practices

### Event Processing

1. **Idempotency**: Events may be replayed, so ensure idempotent processing
2. **Ordering**: Process events in block order to maintain consistency
3. **Error Handling**: Implement robust error handling for malformed events
4. **Monitoring**: Track event processing metrics and failures

### Manifest Management

1. **Versioning**: Include version numbers in manifests for compatibility
2. **Validation**: Validate manifest integrity before publishing
3. **Backup**: Maintain backup copies of manifests for recovery
4. **Notifications**: Notify downstream systems of manifest updates

### Performance Considerations

1. **Batching**: Batch event processing for high-volume chains
2. **Caching**: Cache frequently accessed validator data
3. **Compression**: Compress manifests for efficient storage and transmission
4. **CDN**: Use CDN for global manifest distribution

## Integration Examples

### Hyperlane Validator Setup

```bash
# 1. Validator sets duty metadata
duty tx set-duty-metadata \
  cosmosvaloper1abc123def456 \
  0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef \
  s3://my-bucket/hyperlane/checkpoints/ \
  --from validator

# 2. Sidecar detects duty_metadata_set event
# 3. Sidecar queries /duty.DutySet for complete state
# 4. Sidecar generates updated manifest
# 5. Hyperlane validator reads manifest and configures checkpoint signing
```

### Monitoring Dashboard

```python
# Monitor validator duty status
async def monitor_duty_status():
    events = await subscribe_to_duty_events()
    
    for event in events:
        if event["type"] == "duty_validator_bonded":
            print(f"✅ Validator {event['attributes']['val_addr']} joined duty set")
        elif event["type"] == "duty_validator_removed":
            print(f"❌ Validator {event['attributes']['val_addr']} left duty set")
        elif event["type"] == "duty_metadata_set":
            print(f"📝 Validator {event['attributes']['val_addr']} updated metadata")
```

This event system provides a robust foundation for building sophisticated duty monitoring and management systems.
