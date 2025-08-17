package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers the necessary x/duty interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSetDutyMetadata{}, "duty/SetDutyMetadata", nil)
	cdc.RegisterConcrete(&MsgRotateCheckpointKey{}, "duty/RotateCheckpointKey", nil)
	cdc.RegisterConcrete(&MsgBindCheckpointKey{}, "duty/BindCheckpointKey", nil)
}

// RegisterInterfaces registers the x/duty interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetDutyMetadata{},
		&MsgRotateCheckpointKey{},
		&MsgBindCheckpointKey{},
	)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
