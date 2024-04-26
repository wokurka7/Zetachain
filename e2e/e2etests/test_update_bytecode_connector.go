package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/e2e/contracts/testconnectorzevm"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// TestUpdateBytecodeConnector tests updating the bytecode of a connector and interact with it
func TestUpdateBytecodeConnector(r *runner.E2ERunner, _ []string) {
	// Can withdraw 10ZETA
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(10))
	r.DepositAndApproveWZeta(amount)
	tx := r.WithdrawZeta(amount, true)
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Errorf(
			"expected cctx status to be %s; got %s, message %s",
			crosschaintypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		))
	}

	// Deploy the test contract
	newTestConnectorAddr, tx, _, err := testconnectorzevm.DeployTestZetaConnectorZEVM(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.WZetaAddr,
	)
	if err != nil {
		panic(err)
	}

	// Wait for the contract to be deployed
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("contract deployment failed")
	}

	// Get the code hash of the new contract
	codeHashRes, err := r.FungibleClient.CodeHash(r.Ctx, &fungibletypes.QueryCodeHashRequest{
		Address: newTestConnectorAddr.String(),
	})
	if err != nil {
		panic(err)
	}
	r.Logger.Info("New contract code hash: %s", codeHashRes.CodeHash)

	r.Logger.Info("Updating the bytecode of the Connector")
	msg := fungibletypes.NewMsgUpdateContractBytecode(
		r.ZetaTxServer.GetAccountAddress(0),
		r.ConnectorZEVMAddr.Hex(),
		codeHashRes.CodeHash,
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Update connector bytecode tx hash: %s", res.TxHash)

	r.Logger.Info("Can interact with the new code of the contract")
	testConnectorContract, err := testconnectorzevm.NewTestZetaConnectorZEVM(r.ConnectorZEVMAddr, r.ZEVMClient)
	if err != nil {
		panic(err)
	}

	response, err := testConnectorContract.Foo(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	if response != "foo" {
		panic("unexpected response")
	}

	// Can continue to interact with the connector: withdraw 10ZETA
	r.DepositAndApproveWZeta(amount)
	tx = r.WithdrawZeta(amount, true)
	cctx = utils.WaitCctxMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Errorf(
			"expected cctx status to be %s; got %s, message %s",
			crosschaintypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		))
	}
}
