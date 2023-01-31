package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	"github.com/zeta-chain/zetacore/zetaclient"
	"math/big"
	"time"
)

var (
	BTCDeployerAddress *btcutil.AddressWitnessPubKeyHash
)

func (sm *SmokeTest) TestBitcoinSetup() {
	LoudPrintf("Setup Bitcoin\n")

	btc := sm.btcRpcClient
	_, err := btc.CreateWallet("smoketest", rpcclient.WithCreateWalletBlank())
	if err != nil {
		panic(err)
	}
	skBytes, err := hex.DecodeString(DeployerPrivateKey)
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, &chaincfg.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	err = btc.ImportPrivKeyRescan(privkeyWIF, "deployer", true)
	if err != nil {
		panic(err)
	}
	BTCDeployerAddress, err = btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(privkeyWIF.PrivKey.PubKey().SerializeCompressed()), &chaincfg.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("BTCDeployerAddress: %s\n", BTCDeployerAddress.EncodeAddress())

	err = btc.ImportAddress(BTCTSSAddress.EncodeAddress())
	if err != nil {
		panic(err)
	}
	_, err = btc.GenerateToAddress(101, BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	bal, err := btc.GetBalance("*")
	if err != nil {
		panic(err)
	}
	_, err = btc.GenerateToAddress(1, BTCTSSAddress, nil)
	if err != nil {
		panic(err)
	}
	bal, err = btc.GetBalance("*")
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance: %f\n", bal.ToBTC())

	bals, err := btc.GetBalances()
	if err != nil {
		panic(err)
	}
	fmt.Printf("balances: \n")
	fmt.Printf("  mine: %+v\n", bals.Mine)
	if bals.WatchOnly != nil {
		fmt.Printf("  watchonly: %+v\n", bals.WatchOnly)
	}
	fmt.Printf("TSS Address: %s\n", BTCTSSAddress.EncodeAddress())
	utxos, err := btc.ListUnspent()
	if err != nil {
		panic(err)
	}
	for _, utxo := range utxos {
		fmt.Printf("utxo: %+v\n", utxo)
	}

	// send 1 BTC to TSS address
	input0 := btcjson.TransactionInput{utxos[0].TxID, utxos[0].Vout}
	input1 := btcjson.TransactionInput{utxos[1].TxID, utxos[1].Vout}
	inputs := []btcjson.TransactionInput{input0, input1}
	fee := btcutil.Amount(0.0001 * btcutil.SatoshiPerBitcoin)
	change := btcutil.Amount((utxos[0].Amount+utxos[1].Amount)*(btcutil.SatoshiPerBitcoin)) - fee - btcutil.Amount(1*btcutil.SatoshiPerBitcoin)
	amounts := map[btcutil.Address]btcutil.Amount{
		BTCTSSAddress:      btcutil.Amount(1 * btcutil.SatoshiPerBitcoin),
		BTCDeployerAddress: change,
	}
	tx, err := btc.CreateRawTransaction(inputs, amounts, nil)
	if err != nil {
		panic(err)
	}
	// contruct memo just to deposit BTC into deployer address
	// the bytes in the memo (following OP_RETURN) is of format:
	// [ OP_RETURN(6a) <length of memo> <memo> ]
	// where <memo> is ASCII encoding of the base64 bytes (!we do this because popular bitcoin wallet
	// only input ASCII characters, and we need to encode binary data. We pick base64 StdEncoding).
	addrB64Str := base64.StdEncoding.EncodeToString(DeployerAddress.Bytes())

	addrB64StrLen := len(addrB64Str)
	fmt.Printf("addrB64StrLen: %d\naddrB64Str: %s\naddrB64StrASCII: %x\n", addrB64StrLen, addrB64Str, []byte(addrB64Str))
	//data := make([]byte, addrB64StrLen)
	//data[0] = byte(addrB64StrLen)
	//copy(data[1:], addrB64Str)
	nulldata, err := txscript.NullDataScript([]byte(addrB64Str)) // this adds a OP_RETURN + single BYTE len prefix to the data
	if err != nil {
		panic(err)
	}
	fmt.Printf("nulldata (len %d): %x\n", len(nulldata), nulldata)
	if err != nil {
		panic(err)
	}
	memoOutput := wire.TxOut{Value: 0, PkScript: nulldata}
	tx.TxOut = append(tx.TxOut, &memoOutput)
	tx.TxOut[1], tx.TxOut[2] = tx.TxOut[2], tx.TxOut[1]

	fmt.Printf("raw transaction: \n")
	for idx, txout := range tx.TxOut {
		fmt.Printf("txout %d\n", idx)
		fmt.Printf("  value: %d\n", txout.Value)
		fmt.Printf("  PkScript: %x\n", txout.PkScript)
	}
	stx, signed, err := btc.SignRawTransactionWithWallet(tx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("signed tx: all inputs signed?: %+v\n", signed)
	txid, err := btc.SendRawTransaction(stx, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("txid: %+v\n", txid)
	_, err = btc.GenerateToAddress(6, BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	gtx, err := btc.GetTransaction(txid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("rawtx confirmation: %d\n", gtx.BlockIndex)
	rawtx, err := btc.GetRawTransactionVerbose(txid)
	if err != nil {
		panic(err)
	}

	events := zetaclient.FilterAndParseIncomingTx([]btcjson.TxRawResult{*rawtx}, 0, BTCTSSAddress.EncodeAddress(), &log.Logger)
	fmt.Printf("bitcoin intx events:\n")
	for _, event := range events {
		fmt.Printf("  TxHash: %s\n", event.TxHash)
		fmt.Printf("  From: %s\n", event.FromAddress)
		fmt.Printf("  To: %s\n", event.ToAddress)
		fmt.Printf("  Amount: %d\n", event.Value)
		fmt.Printf("  Memo: %x\n", event.MemoBytes)
	}

	fmt.Printf("testing if the deposit into BTC ZRC20 is successful...")

	SystemContract, err := contracts.NewSystemContract(HexToAddress(SystemContractAddr), sm.zevmClient)
	if err != nil {
		panic(err)
	}
	sm.SystemContract = SystemContract
	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := contracts.NewZRC20(BTCZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(5 * time.Second)
		balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
		if err != nil {
			panic(err)
		}
		if balance.Cmp(big.NewInt(1*btcutil.SatoshiPerBitcoin)) != 0 {
			fmt.Printf("waiting for BTC balance to show up in ZRC contract... current bal %d\n", balance)
		} else {
			fmt.Printf("BTC balance is in ZRC contract! Success\n")
			break
		}
	}
}

func (sm *SmokeTest) TestBitcoinWithdraw() {
	LoudPrintf("Testing Bitcoin ZRC20 Withdraw...\n")
	// withdraw 0.1 BTC from ZRC20 to BTC address
	// first, approve the ZRC20 contract to spend 1 BTC from the deployer address
	amount := big.NewInt(0.1 * btcutil.SatoshiPerBitcoin)
	SystemContract, err := contracts.NewSystemContract(HexToAddress(SystemContractAddr), sm.zevmClient)
	if err != nil {
		panic(err)
	}
	sm.SystemContract = SystemContract
	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := contracts.NewZRC20(BTCZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	if balance.Cmp(amount) < 0 {
		panic(fmt.Errorf("not enough balance in ZRC20 contract"))
	}
	// approve the ZRC20 contract to spend 1 BTC from the deployer address
	{

		tx, err := BTCZRC20.Approve(sm.zevmAuth, BTCZRC20Addr, big.NewInt(amount.Int64()*2)) // approve more to cover withdraw fee
		if err != nil {
			panic(err)
		}
		receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
		fmt.Printf("approve receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("approve receipt status is not 1"))
		}
	}
	// withdraw 0.1 BTC from ZRC20 to BTC address
	{
		_, gasFee, err := BTCZRC20.WithdrawGasFee(&bind.CallOpts{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("withdraw gas fee: %d\n", gasFee)
		tx, err := BTCZRC20.Withdraw(sm.zevmAuth, []byte(BTCDeployerAddress.EncodeAddress()), amount)
		if err != nil {
			panic(err)
		}
		receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
		fmt.Printf("withdraw receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("withdraw receipt status is not 1"))
		}
		WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), sm.cctxClient)
	}
}
