package types

const (
	ModuleName = "duty"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	// KV prefixes
	// DutyMeta: validator-consensus-address -> DutyMetadata
	DutyMetaPrefix = []byte{0x01}
)

func DutyMetaKey(valConsAddr []byte) []byte {
	return append(DutyMetaPrefix, valConsAddr...)
}
