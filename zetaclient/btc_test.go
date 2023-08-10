package zetaclient

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/zeta-chain/zetacore/common"
)

type BTCSignTestSuite struct {
	suite.Suite
	testSigner *TestSigner
	db         *gorm.DB
}

const (
	prevOut = "07a84f4bd45a633e93871be5c98d958afd13a37f3cf5010f40eec0840d19f5fa"
	// tb1q7r6lnqjhvdjuw9uf4ehx7fs0euc6cxnqz7jj50
	pk        = "cQkjdfeMU8vHvE6jErnFVqZYYZnGGYy64jH6zovbSXdfTjte6QgY"
	utxoCount = 5
)

func (suite *BTCSignTestSuite) SetupTest() {
	//skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	//privateKey, err := crypto.HexToECDSA(skHex)
	//pkBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	//suite.T().Logf("pubkey: %d", len(pkBytes))
	//suite.Require().NoError(err)

	wif, _ := btcutil.DecodeWIF(pk)
	privateKey := secp.PrivateKey(*wif.PrivKey)

	suite.testSigner = &TestSigner{ // fake TSS
		PrivKey: privateKey.ToECDSA(),
	}
	addr := suite.testSigner.BTCAddressWitnessPubkeyHash()
	suite.T().Logf("segwit addr: %s", addr)

	db, err := gorm.Open(sqlite.Open(TempSQLiteDbPath), &gorm.Config{})
	suite.NoError(err)

	suite.db = db
}

func (suite *BTCSignTestSuite) TearDownSuite() {
}

func (suite *BTCSignTestSuite) TestSign() {

	// build a tx used for both signatures
	tx, txSigHashes, idx, amt, subscript, privKey, compress, err := buildTX()
	suite.Require().NoError(err)

	// sign tx using wallet signature
	walletSignedTX, err := getWalletTX(tx, txSigHashes, idx, amt, subscript, txscript.SigHashAll, privKey, compress)
	suite.Require().NoError(err)
	suite.T().Logf("wallet signed tx : %v\n", walletSignedTX)

	// sign tx using tss signature

	tssSignedTX, err := getTSSTX(suite.testSigner, tx, txSigHashes, idx, amt, subscript, txscript.SigHashAll)
	suite.Require().NoError(err)
	suite.T().Logf("tss signed tx :    %v\n", tssSignedTX)
}

func (suite *BTCSignTestSuite) TestSubmittedTx() {

}

func TestBTCSign(t *testing.T) {
	suite.Run(t, new(BTCSignTestSuite))
}

func buildTX() (*wire.MsgTx, *txscript.TxSigHashes, int, int64, []byte, *btcec.PrivateKey, bool, error) {
	wif, err := btcutil.DecodeWIF(pk)
	if err != nil {
		return nil, nil, 0, 0, nil, nil, false, err
	}
	fmt.Printf("is wif.netid testnet3 %v\n", wif.IsForNet(&chaincfg.TestNet3Params))
	fmt.Printf("is wif pubkey compressed %v\n", wif.CompressPubKey)

	addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
	if err != nil {
		return nil, nil, 0, 0, nil, nil, false, err
	}
	fmt.Printf("addr %v\n", addr.EncodeAddress())

	hash, err := chainhash.NewHashFromStr(prevOut)
	if err != nil {
		return nil, nil, 0, 0, nil, nil, false, err
	}
	outpoint := wire.NewOutPoint(hash, 0)

	// build tx
	tx := wire.NewMsgTx(wire.TxVersion)
	txIn := wire.NewTxIn(outpoint, nil, nil)
	tx.AddTxIn(txIn)

	pkScript, err := payToWitnessPubKeyHashScript(addr.WitnessProgram())
	if err != nil {
		return nil, nil, 0, 0, nil, nil, false, err
	}
	txOut := wire.NewTxOut(47000, pkScript)
	tx.AddTxOut(txOut)

	prevFetcher := txscript.NewCannedPrevOutputFetcher(
		txOut.PkScript, txOut.Value,
	)

	txSigHashes := txscript.NewTxSigHashes(tx, prevFetcher)

	privKey := *wif.PrivKey

	return tx, txSigHashes, 0, int64(65236), pkScript, &privKey, wif.CompressPubKey, nil
}

func getWalletTX(
	tx *wire.MsgTx,
	sigHashes *txscript.TxSigHashes,
	idx int,
	amt int64,
	subscript []byte,
	hashType txscript.SigHashType,
	privKey *btcec.PrivateKey,
	compress bool,
) (string, error) {
	txWitness, err := txscript.WitnessSignature(tx, sigHashes, idx, amt, subscript, hashType, privKey, compress)
	if err != nil {
		return "", err
	}

	tx.TxIn[0].Witness = txWitness

	buf := new(bytes.Buffer)
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	walletTx := hex.EncodeToString(buf.Bytes())
	return walletTx, nil
}

func getTSSTX(tss *TestSigner, tx *wire.MsgTx, sigHashes *txscript.TxSigHashes, idx int, amt int64, subscript []byte, hashType txscript.SigHashType) (string, error) {
	witnessHash, err := txscript.CalcWitnessSigHash(subscript, sigHashes, txscript.SigHashAll, tx, idx, amt)
	if err != nil {
		return "", err
	}

	// create signature
	sig65B, err := tss.Sign(witnessHash, 10, &common.Chain{})
	if err != nil {
		return "", err
	}

	var r, s btcec.ModNScalar
	if overflow := r.SetByteSlice(sig65B[:32]); overflow {
		return "", errors.New("r: byte slice overflow")
	}
	if overflow := s.SetByteSlice(sig65B[32:64]); overflow {
		return "", errors.New("s: byte slice overflow")
	}

	sig := ecdsa.NewSignature(&r, &s)

	pkCompressed := tss.PubKeyCompressedBytes()
	txWitness := wire.TxWitness{append(sig.Serialize(), byte(hashType)), pkCompressed}
	tx.TxIn[0].Witness = txWitness

	buf := new(bytes.Buffer)
	err = tx.Serialize(buf)
	if err != nil {
		return "", err
	}

	tssTX := hex.EncodeToString(buf.Bytes())
	return tssTX, nil
}
