package amm

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

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
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              ammapi.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions:    []*autocliv1.RpcCommandOptions{},
			EnhanceCustomCommand: true,
		},
	}
}
