# x/duty Sidecar Service

A lightweight sidecar service that watches on-chain events, queries the `x/duty` module, and produces a canonical, machine-readable manifest for Hyperlane components.

## Overview

The sidecar service bridges the gap between the Cosmos SDK `x/duty` module and Hyperlane infrastructure by:

- **Watching on-chain events** for validator set changes
- **Polling the duty set** to maintain current validator information
- **Producing standardized manifests** for Hyperlane components
- **Managing deterministic storage layouts** for checkpoint signatures
- **Handling one-time Hyperlane announcements** for validators

## Architecture

```
Cosmos Chain (x/duty module)
         ↓
    Sidecar Service
         ↓
    Manifest Generation
         ↓
    Hyperlane Components
    (Relayers, ISMs, Validators)
```

## Core Responsibilities

### 1. Event Monitoring
- Polls `/duty.DutySet` every N blocks (configurable)
- Listens for `duty_validator_bonded` and `duty_validator_removed` events
- Triggers manifest updates on validator set changes

### 2. Manifest Generation
Produces a standardized JSON manifest:

```json
{
  "chain_id": "cosmoshub-4",
  "quorum": {"num": 2, "den": 3},
  "validators": [
    {
      "consensus_address": "cosmosvalcons1abcdef...",
      "voting_power": "1000000",
      "checkpoint_pub_key": "0x1234567890abcdef...",
      "checkpoint_storage_uri": "s3://bucket/prefix/cosmosvalcons1abcdef.../"
    }
  ],
  "version": 1,
  "asof_height": 12345678
}
```

### 3. HTTP Endpoint
- Exposes `/manifest` endpoint for read-only access
- Returns the latest manifest JSON
- Includes cache headers for efficient polling

### 4. Deterministic Storage Layout
Recommends standardized storage structure:

```
s3://<bucket>/hyperlane/<chain-id>/validators/<consensus-address>/checkpoints/{block}-{epoch}.json
```

## Checkpoint Signature Schema

Each checkpoint signature file follows this schema:

```json
{
  "chain_id": "cosmoshub-4",
  "origin_block": 12345678,
  "checkpoint_root": "0x1234567890abcdef...",
  "signer_pubkey": "0xabcdef1234567890...",
  "signature": "0x9876543210fedcba...",
  "signer_consensus_address": "cosmosvalcons1abcdef...",
  "signed_at_unix": 1640995200
}
```

## Deployment

### Docker Compose Example

```yaml
version: '3.8'
services:
  duty-sidecar:
    build: .
    environment:
      - DUTY_GRPC=localhost:9090
      - DUTY_RPC=http://localhost:26657
      - CHAIN_ID=cosmoshub-4
      - POLL_INTERVAL=30
      - OUTPUT_PATH=/app/manifest.json
      - HYP_MAILBOX=0x1234...  # Optional
      - HYP_RPC=https://eth-mainnet.alchemyapi.io/v2/...  # Optional
      - ANNOUNCE_KEY=0xabcd...  # Optional
    volumes:
      - ./manifests:/app/manifests
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o duty-sidecar .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /app/duty-sidecar .
EXPOSE 8080
CMD ["./duty-sidecar"]
```

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DUTY_GRPC` | gRPC endpoint for duty queries | Yes | - |
| `DUTY_RPC` | RPC endpoint for event streaming | Yes | - |
| `CHAIN_ID` | Cosmos chain identifier | Yes | - |
| `POLL_INTERVAL` | Polling interval in seconds | No | 30 |
| `OUTPUT_PATH` | Path for manifest file | No | ./manifest.json |
| `HYP_MAILBOX` | Hyperlane mailbox address | No | - |
| `HYP_RPC` | Hyperlane RPC endpoint | No | - |
| `ANNOUNCE_KEY` | Private key for announcements | No | - |

## API Endpoints

### GET /manifest
Returns the latest duty set manifest.

**Response:**
```json
{
  "chain_id": "cosmoshub-4",
  "quorum": {"num": 2, "den": 3},
  "validators": [...],
  "version": 1,
  "asof_height": 12345678
}
```

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "last_update": "2024-01-01T00:00:00Z",
  "block_height": 12345678
}
```

## Deterministic Key Mapping

### Attestation Flow

The sidecar supports deterministic key mapping between Cosmos consensus validators and Hyperlane checkpoint signers:

1. **Validator Registration**: Validator sets duty metadata on-chain
2. **Key Attestation**: Sidecar verifies the binding between consensus address and checkpoint key
3. **Manifest Generation**: Sidecar includes attested keys in the manifest
4. **Hyperlane Integration**: Hyperlane components use attested keys for verification

### Verification Process

```go
// Example verification logic
func verifyKeyAttestation(consAddr string, checkpointKey string) bool {
    // 1. Verify validator exists in consensus set
    // 2. Verify checkpoint key is properly formatted (secp256k1)
    // 3. Verify storage URI is accessible
    // 4. Generate attestation signature
    return true
}
```

## Integration with Hyperlane Validators

### One-Time Announcement

The sidecar can optionally handle one-time Hyperlane announcements:

```bash
# Announce validator to Hyperlane
curl -X POST http://localhost:8080/announce \
  -H "Content-Type: application/json" \
  -d '{
    "validator_address": "cosmosvalcons1...",
    "checkpoint_pub_key": "0x1234..."
  }'
```

### Validator Configuration

Validators can configure their Hyperlane software to use the sidecar manifest:

```yaml
# hyperlane-validator.yaml
validator:
  manifest_url: "http://sidecar:8080/manifest"
  checkpoint_signer:
    type: "local"
    key: "checkpoint_key.pem"
  storage:
    type: "s3"
    bucket: "my-hyperlane-checkpoints"
    prefix: "hyperlane/cosmoshub-4/validators/"
```

## Monitoring and Logging

### Log Format

```
2024-01-01T00:00:00Z INFO [SIDECAR] Manifest updated at height 12345678
2024-01-01T00:00:01Z INFO [SIDECAR] Validator cosmosvalcons1... bonded
2024-01-01T00:00:02Z INFO [SIDECAR] Validator cosmosvalcons1... set metadata
```

### Metrics

The sidecar exposes Prometheus metrics:

- `duty_manifest_updates_total`: Total manifest updates
- `duty_validator_changes_total`: Total validator set changes
- `duty_poll_duration_seconds`: Time to poll duty set
- `duty_manifest_size_bytes`: Size of current manifest

## Security Considerations

### Key Management
- Never store private keys in the sidecar
- Use environment variables or secure key management
- Rotate keys regularly

### Access Control
- Restrict access to the sidecar API
- Use HTTPS in production
- Implement rate limiting

### Data Validation
- Validate all on-chain data before including in manifest
- Verify checkpoint signatures before storing
- Check storage accessibility

## Troubleshooting

### Common Issues

1. **Manifest not updating**
   - Check gRPC connection to duty module
   - Verify event subscription
   - Check logs for errors

2. **Storage access issues**
   - Verify S3 bucket permissions
   - Check network connectivity
   - Validate storage URIs

3. **Validator metadata missing**
   - Ensure validators have set duty metadata
   - Check validator bonding status
   - Verify consensus address format

### Debug Mode

Enable debug logging:

```bash
export LOG_LEVEL=debug
./duty-sidecar
```

## Future Enhancements

### Planned Features

1. **Governance Integration**: Support for governance-controlled parameter changes
2. **Multi-Chain Support**: Support for multiple Cosmos chains
3. **Advanced Monitoring**: Dashboard for duty set health
4. **Automated Recovery**: Automatic recovery from failures
5. **Caching Layer**: Redis-based caching for improved performance

### Extension Points

- **Custom Storage Providers**: Support for IPFS, Filecoin, etc.
- **Advanced Attestation**: Multi-signature attestation schemes
- **Performance Optimization**: Parallel processing of validator updates
- **Integration APIs**: REST APIs for external systems
