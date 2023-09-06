package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

const node = "tcp://3.218.170.198:26657"
const signer = ""
const chainID = "athens_7001-1"
const amount = "100000000000000000000"
const broadcastMode = "sync"

//const node = "tcp://localhost:26657"
//const signer = "zeta"
//const chain_id = "localnet_101-1"
//const amount = "100000000" // Amount in azeta
//const broadcast_mode = "block"

type TokenDistribution struct {
	Address           string   `json:"address"`
	BalanceBefore     sdk.Coin `json:"balance_before"`
	BalanceAfter      sdk.Coin `json:"balance_after"`
	TokensDistributed sdk.Coin `json:"tokens_distributed"`
}

func main() {
	GetCCTXHash()

	file, _ := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "address-list.json"))
	addresses, err := readLines(file)
	if err != nil {
		panic(err)
	}
	addresses = removeDuplicates(addresses)
	fileS, _ := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "successful_address.json"))
	fileF, _ := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "failed_address.json"))

	distributionList := make([]TokenDistribution, len(addresses))
	for i, address := range addresses {
		cmd := exec.Command("zetacored", "q", "bank", "balances", address, "--output", "json", "--denom", "azeta", "--node", node) // #nosec G204
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(cmd.String())
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			return
		}
		balance := sdk.Coin{}
		err = json.Unmarshal(output, &balance)
		if err != nil {
			panic(err)
		}
		distributionAmount, ok := sdkmath.NewIntFromString(amount)
		if !ok {
			panic("parse error for amount")
		}
		distributionList[i] = TokenDistribution{
			Address:           address,
			BalanceBefore:     balance,
			TokensDistributed: sdk.NewCoin(config.BaseDenom, distributionAmount),
		}
	}

	args := []string{"tx", "bank", "multi-send", signer}
	for _, address := range addresses {
		args = append(args, address)
	}

	args = append(args, []string{distributionList[0].TokensDistributed.String(), "--keyring-backend", "test", "--chain-id", chainID, "--yes",
		"--broadcast-mode", broadcastMode, "--gas=auto", "--gas-adjustment=2", "--gas-prices=0.001azeta", "--node", node}...)

	cmd := exec.Command("zetacored", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd.String())
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))

	time.Sleep(7 * time.Second)

	for i, address := range addresses {
		cmd := exec.Command("zetacored", "q", "bank", "balances", address, "--output", "json", "--denom", "azeta", "--node", node) // #nosec G204
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(cmd.String())
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			return
		}
		balance := sdk.Coin{}
		err = json.Unmarshal(output, &balance)
		if err != nil {
			panic(err)
		}
		distributionList[i].BalanceAfter = balance
	}
	var successfullDistributions []TokenDistribution
	var failedDistributions []TokenDistribution
	for _, distribution := range distributionList {
		if distribution.BalanceAfter.Sub(distribution.BalanceBefore).IsEqual(distribution.TokensDistributed) {
			successfullDistributions = append(successfullDistributions, distribution)
		} else {
			failedDistributions = append(failedDistributions, distribution)
		}
	}
	successFile, _ := json.MarshalIndent(successfullDistributions, "", " ")
	_ = os.WriteFile(fileS, successFile, 0600)
	failedFile, _ := json.MarshalIndent(failedDistributions, "", " ")
	_ = os.WriteFile(fileF, failedFile, 0600)

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path) // #nosec G304
	if err != nil {
		return nil, err
	}
	/* #nosec G307 */
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func removeDuplicates(s []string) []string {
	bucket := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}

func GetCCTXHash() {
	amount := "10000000000000"
	asset := ""
	creator := "zeta1mte0r3jzkf2rkd7ex4p3xsd3fxqg7q29q0wxl5"
	gas_limit := 90000
	in_block_height := 9649384
	in_tx_hash := "0x48ae453a1d4d65774320570646d4c9b2287e1f86ef22beec191a0f2295579d1f"
	message := "e66b7b71070747c43cbdbdf607f25da8f073239e48f80608b672dc30dc7e3dbbd0343c5f02c738eb2c0e5ec8794aeba837a894dff2c3a605a353e56e"
	receiver := "0x2C0E5EC8794aEba837a894dFf2C3a605a353E56e"
	receiver_chain := 7001
	sender := "0x2C0E5EC8794aEba837a894dFf2C3a605a353E56e"
	sender_chain_id := 5
	tx_origin := "0x2C0E5EC8794aEba837a894dFf2C3a605a353E56e"
	msg := crosschaintypes.NewMsgSendVoter(
		creator,
		sender,
		int64(sender_chain_id),
		tx_origin,
		receiver,
		int64(receiver_chain),
		sdkmath.NewUintFromString(amount),
		message,
		in_tx_hash,
		uint64(in_block_height),
		uint64(gas_limit),
		common.CoinType_Gas,
		asset)
	fmt.Println(msg.Digest())

}
