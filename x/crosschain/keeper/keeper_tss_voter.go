package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// MESSAGES

func (k msgServer) CreateTSSVoter(goCtx context.Context, msg *types.MsgCreateTSSVoter) (*types.MsgCreateTSSVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAuthorizedNodeAccount(ctx, msg.Creator) {
		return nil, errorsmod.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s does not have a node account set", msg.Creator))
	}
	// No need to create a ballot if keygen does nto exist
	keygen, found := k.GetKeygen(ctx)
	if !found {
		return &types.MsgCreateTSSVoterResponse{}, types.ErrKeygenNotFound
	}
	index := msg.Digest()
	// Add votes and Set Ballot
	// GetBallot checks against the supported chains list before querying for Ballot
	ballot := zetaObserverTypes.Ballot{}
	ballot, found = k.zetaObserverKeeper.GetBallot(ctx, index)
	if !found {
		var voterList []string

		for _, nodeAccount := range k.GetAllNodeAccount(ctx) {
			voterList = append(voterList, nodeAccount.Creator)
		}
		ballot = zetaObserverTypes.Ballot{
			Index:            "",
			BallotIdentifier: index,
			VoterList:        voterList,
			Votes:            zetaObserverTypes.CreateVotes(len(msg.Creator)),
			ObservationType:  zetaObserverTypes.ObservationType_TSSKeyGen,
			BallotThreshold:  sdk.MustNewDecFromStr("1.00"),
			BallotStatus:     zetaObserverTypes.BallotStatus_BallotInProgress,
		}
	}
	err := error(nil)
	if msg.Status == common.ReceiveStatus_Success {
		ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, zetaObserverTypes.VoteType_SuccessObservation)
		if err != nil {
			return &types.MsgCreateTSSVoterResponse{}, err
		}
	} else if msg.Status == common.ReceiveStatus_Failed {
		ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, zetaObserverTypes.VoteType_FailureObservation)
		if err != nil {
			return &types.MsgCreateTSSVoterResponse{}, err
		}
	}

	ballot, isFinalized := k.CheckIfBallotIsFinalized(ctx, ballot)
	if !isFinalized {
		return &types.MsgCreateTSSVoterResponse{}, nil
	}

	if ballot.BallotStatus != zetaObserverTypes.BallotStatus_BallotFinalized_FailureObservation {
		k.SetTSS(ctx, types.TSS{
			TssPubkey:           msg.TssPubkey,
			SignerList:          ballot.VoterList,
			FinalizedZetaHeight: ctx.BlockHeight(),
			KeyGenZetaHeight:    msg.KeyGenZetaHeight,
		})
		keygen.Status = types.KeygenStatus_KeyGenSuccess
	} else if ballot.BallotStatus == zetaObserverTypes.BallotStatus_BallotFinalized_FailureObservation {
		keygen.Status = types.KeygenStatus_KeyGenFailed
	}
	k.SetKeygen(ctx, keygen)
	// Remove ballot
	return &types.MsgCreateTSSVoterResponse{}, nil
}
