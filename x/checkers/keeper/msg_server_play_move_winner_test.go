package keeper_test

import (
	"context"
	"testing"

	"github.com/b9lab/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type GameMoveTest struct {
	player string
	fromX  uint64
	fromY  uint64
	toX    uint64
	toY    uint64
}

var (
	game1Moves = []GameMoveTest{
		{"b", 1, 2, 2, 3}, // "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"
		{"r", 0, 5, 1, 4}, // "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"
		{"b", 2, 3, 0, 5}, // "*b*b*b*b|b*b*b*b*|***b*b*b|********|********|b*r*r*r*|*r*r*r*r|r*r*r*r*"
		{"r", 4, 5, 3, 4}, // "*b*b*b*b|b*b*b*b*|***b*b*b|********|***r****|b*r***r*|*r*r*r*r|r*r*r*r*"
		{"b", 3, 2, 2, 3}, // "*b*b*b*b|b*b*b*b*|*****b*b|**b*****|***r****|b*r***r*|*r*r*r*r|r*r*r*r*"
		{"r", 3, 4, 1, 2}, // "*b*b*b*b|b*b*b*b*|*r***b*b|********|********|b*r***r*|*r*r*r*r|r*r*r*r*"
		{"b", 0, 1, 2, 3}, // "*b*b*b*b|**b*b*b*|*****b*b|**b*****|********|b*r***r*|*r*r*r*r|r*r*r*r*"
		{"r", 2, 5, 3, 4}, // "*b*b*b*b|**b*b*b*|*****b*b|**b*****|***r****|b*****r*|*r*r*r*r|r*r*r*r*"
		{"b", 2, 3, 4, 5}, // "*b*b*b*b|**b*b*b*|*****b*b|********|********|b***b*r*|*r*r*r*r|r*r*r*r*"
		{"r", 5, 6, 3, 4}, // "*b*b*b*b|**b*b*b*|*****b*b|********|***r****|b*****r*|*r*r***r|r*r*r*r*"
		{"b", 5, 2, 4, 3}, // "*b*b*b*b|**b*b*b*|*******b|****b***|***r****|b*****r*|*r*r***r|r*r*r*r*"
		{"r", 3, 4, 5, 2}, // "*b*b*b*b|**b*b*b*|*****r*b|********|********|b*****r*|*r*r***r|r*r*r*r*"
		{"b", 6, 1, 4, 3}, // "*b*b*b*b|**b*b***|*******b|****b***|********|b*****r*|*r*r***r|r*r*r*r*"
		{"r", 6, 5, 5, 4}, // "*b*b*b*b|**b*b***|*******b|****b***|*****r**|b*******|*r*r***r|r*r*r*r*"
		{"b", 4, 3, 6, 5}, // "*b*b*b*b|**b*b***|*******b|********|********|b*****b*|*r*r***r|r*r*r*r*"
		{"r", 7, 6, 5, 4}, // "*b*b*b*b|**b*b***|*******b|********|*****r**|b*******|*r*r****|r*r*r*r*"
		{"b", 7, 2, 6, 3}, // "*b*b*b*b|**b*b***|********|******b*|*****r**|b*******|*r*r****|r*r*r*r*"
		{"r", 5, 4, 7, 2}, // "*b*b*b*b|**b*b***|*******r|********|********|b*******|*r*r****|r*r*r*r*"
		{"b", 4, 1, 3, 2}, // "*b*b*b*b|**b*****|***b***r|********|********|b*******|*r*r****|r*r*r*r*"
		{"r", 3, 6, 4, 5}, // "*b*b*b*b|**b*****|***b***r|********|********|b***r***|*r******|r*r*r*r*"
		{"b", 5, 0, 4, 1}, // "*b*b***b|**b*b***|***b***r|********|********|b***r***|*r******|r*r*r*r*"
		{"r", 2, 7, 3, 6}, // "*b*b***b|**b*b***|***b***r|********|********|b***r***|*r*r****|r***r*r*"
		{"b", 0, 5, 2, 7}, // "*b*b***b|**b*b***|***b***r|********|********|****r***|***r****|r*B*r*r*"
		{"r", 4, 5, 3, 4}, // "*b*b***b|**b*b***|***b***r|********|***r****|********|***r****|r*B*r*r*"
		{"b", 2, 7, 4, 5}, // "*b*b***b|**b*b***|***b***r|********|***r****|****B***|********|r***r*r*"
		// Captures again
		{"b", 4, 5, 2, 3}, // "*b*b***b|**b*b***|***b***r|**B*****|********|********|********|r***r*r*"
		{"r", 6, 7, 5, 6}, // "*b*b***b|**b*b***|***b***r|**B*****|********|********|*****r**|r***r***"
		{"b", 2, 3, 3, 4}, // "*b*b***b|**b*b***|***b***r|********|***B****|********|*****r**|r***r***"
		{"r", 0, 7, 1, 6}, // "*b*b***b|**b*b***|***b***r|********|***B****|********|*r***r**|****r***"
		{"b", 3, 2, 4, 3}, // "*b*b***b|**b*b***|*******r|****b***|***B****|********|*r***r**|****r***"
		{"r", 7, 2, 6, 1}, // "*b*b***b|**b*b*r*|********|****b***|***B****|********|*r***r**|****r***"
		{"b", 7, 0, 5, 2}, // "*b*b****|**b*b***|*****b**|****b***|***B****|********|*r***r**|****r***"
		{"r", 1, 6, 2, 5}, // "*b*b****|**b*b***|*****b**|****b***|***B****|**r*****|*****r**|****r***"
		{"b", 3, 4, 1, 6}, // "*b*b****|**b*b***|*****b**|****b***|********|********|*B***r**|****r***"
		{"r", 4, 7, 3, 6}, // "*b*b****|**b*b***|*****b**|****b***|********|********|*B*r*r**|********"
		{"b", 4, 3, 3, 4}, // "*b*b****|**b*b***|*****b**|********|***b****|********|*B*r*r**|********"
		{"r", 5, 6, 4, 5}, // "*b*b****|**b*b***|*****b**|********|***b****|****r***|*B*r****|********"
		{"b", 3, 4, 5, 6}, // "*b*b****|**b*b***|*****b**|********|********|********|*B*r*b**|********"
		{"r", 3, 6, 2, 5}, // "*b*b****|**b*b***|*****b**|********|********|**r*****|*B***b**|********"
		{"b", 1, 6, 3, 4}, // "*b*b****|**b*b***|*****b**|********|***B****|********|*****b**|********"
	}
)

func getPlayer(color string) string {
	if color == "b" {
		return bob
	}
	return carol
}

func playAllMoves(t *testing.T, msgServer types.MsgServer, context context.Context, gameIndex string, moves []GameMoveTest) {
	for _, move := range game1Moves {
		_, err := msgServer.PlayMove(context, &types.MsgPlayMove{
			Creator:   getPlayer(move.player),
			GameIndex: gameIndex,
			FromX:     move.fromX,
			FromY:     move.fromY,
			ToX:       move.toX,
			ToY:       move.toY,
		})
		require.Nil(t, err)
	}
}

func TestPlayMoveUpToWinner(t *testing.T) {
	msgServer, keeper, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)

	playAllMoves(t, msgServer, context, "1", game1Moves)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)

	game, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		MoveCount:   uint64(len(game1Moves)),
		BeforeIndex: "-1",
		AfterIndex:  "-1",
		Deadline:    types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:      "b",
		Wager:       45,
		Denom:       "stake",
	}, game)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.Equal(t, event.Type, "move-played")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "creator", Value: bob},
		{Key: "game-index", Value: "1"},
		{Key: "captured-x", Value: "2"},
		{Key: "captured-y", Value: "5"},
		{Key: "winner", Value: "b"},
		{Key: "board", Value: "*b*b****|**b*b***|*****b**|********|***B****|********|*****b**|********"},
	}, event.Attributes[(len(game1Moves)-1)*6:])
}

func TestPlayMoveUpToWinnerCalledBank(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	payBob := escrow.ExpectPay(context, bob, 45).Times(1)
	payCarol := escrow.ExpectPay(context, carol, 45).Times(1).After(payBob)
	escrow.ExpectRefund(context, bob, 90).Times(1).After(payCarol)

	playAllMoves(t, msgServer, context, "1", game1Moves)
}

func TestCompleteGameAddPlayerInfo(t *testing.T) {
	msgServer, keeper, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)

	playAllMoves(t, msgServer, context, "1", game1Moves)

	bobInfo, found := keeper.GetPlayerInfo(ctx, bob)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          bob,
		WonCount:       1,
		LostCount:      0,
		ForfeitedCount: 0,
	}, bobInfo)
	carolInfo, found := keeper.GetPlayerInfo(ctx, carol)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          carol,
		WonCount:       0,
		LostCount:      1,
		ForfeitedCount: 0,
	}, carolInfo)
}

func TestCompleteGameUpdatePlayerInfo(t *testing.T) {
	msgServer, keeper, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)

	keeper.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          bob,
		WonCount:       1,
		LostCount:      2,
		ForfeitedCount: 3,
	})
	keeper.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          carol,
		WonCount:       4,
		LostCount:      5,
		ForfeitedCount: 6,
	})

	playAllMoves(t, msgServer, context, "1", game1Moves)

	bobInfo, found := keeper.GetPlayerInfo(ctx, bob)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          bob,
		WonCount:       2,
		LostCount:      2,
		ForfeitedCount: 3,
	}, bobInfo)
	carolInfo, found := keeper.GetPlayerInfo(ctx, carol)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          carol,
		WonCount:       4,
		LostCount:      6,
		ForfeitedCount: 6,
	}, carolInfo)
}

func TestCompleteGameLeaderboardAddWinner(t *testing.T) {
	msgServer, keeper, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)

	playAllMoves(t, msgServer, context, "1", game1Moves)

	leaderboard, found := keeper.GetLeaderboard(ctx)
	require.True(t, found)
	require.EqualValues(t, []types.WinningPlayer{
		{
			PlayerAddress: bob,
			WonCount:      1,
			DateAdded:     types.FormatDateAdded(types.GetDateAdded(ctx)),
		},
	}, leaderboard.Winners)
}

func TestCompleteGameLeaderboardUpdatedWinner(t *testing.T) {
	msgServer, keeper, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	keeper.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:    bob,
		WonCount: 2,
	})
	keeper.SetLeaderboard(ctx, types.Leaderboard{
		Winners: []types.WinningPlayer{
			{
				PlayerAddress: bob,
				WonCount:      2,
				DateAdded:     "2006-01-02 15:05:06.999999999 +0000 UTC",
			},
		},
	})

	playAllMoves(t, msgServer, context, "1", game1Moves)

	leaderboard, found := keeper.GetLeaderboard(ctx)
	require.True(t, found)
	require.EqualValues(t, []types.WinningPlayer{
		{
			PlayerAddress: bob,
			WonCount:      3,
			DateAdded:     types.FormatDateAdded(types.GetDateAdded(ctx)),
		},
	}, leaderboard.Winners)
}
