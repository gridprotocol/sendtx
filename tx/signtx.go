package tx

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// make and sign a eth tx for sending, tx data as the param
func MakeSignedTx(client *ethclient.Client,
	sk string, // sk of the sender
	to common.Address, // to address of this tx
	value *big.Int, // value in this tx
	gasLimit uint64, // gas limit of this tx
	data []byte, // data of this tx
) (*types.Transaction, error) {
	// sk
	privateKey, err := crypto.HexToECDSA(sk) // 你的以太坊账户私钥
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type *ecdsa.PublicKey")
	}

	// get the from addr with pk
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// get the nonce from client
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	//gasLimit := uint64(21000)

	// gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// make tx
	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	// sign tx
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
