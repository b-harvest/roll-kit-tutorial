package amm

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/version"

	ammapi "github.com/b-harvest/roll-kit-tutorial/api/simapp/amm/v1beta1"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: ammapi.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current coinswap module parameters information",
				},
				{
					RpcMethod: "Pairs",
					Use:       "pairs",
					Short:     "Query all pairs",
				},
				{
					RpcMethod: "Pair",
					Use:       "pair <pair_id>",
					Short:     "Query a pair",
					Example:   fmt.Sprintf("$ %s query amm pair 1", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              ammapi.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions:    []*autocliv1.RpcCommandOptions{},
			EnhanceCustomCommand: true,
		},
	}
}
