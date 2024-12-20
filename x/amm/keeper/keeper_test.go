package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	simapp "github.com/b-harvest/roll-kit-tutorial/app"
	"github.com/b-harvest/roll-kit-tutorial/x/amm/keeper"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	app    *simapp.SimApp
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = simapp.Setup(s.T(), false)
	s.ctx = s.app.BaseApp.NewContextLegacy(false, tmproto.Header{})
	s.keeper = s.app.AmmKeeper
}
