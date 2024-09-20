package tx

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/grid/contracts/go/registry"

	"github.com/grid/contracts/eth"
)

type Tx struct {
	ep string
	c  *ethclient.Client

	SignedTx *types.Transaction
	// tx in byte
	JsonTx []byte
}

// a nil tx with client
func NewTx(ep string) *Tx {
	// connect to an eth client
	log.Println("connecting client")
	c, err := ethclient.Dial(ep)
	if err != nil {
		log.Fatal(err)
	}

	return &Tx{ep, c, nil, nil}
}

// Make tx for register cp
func (tx *Tx) MakeRegisterTx() error {
	// tx data
	data := RegisterData()

	log.Println("making signed register tx")
	// Make a signed tx
	fmt.Println("cp sk: ", P_SK)
	SignedTx, err := MakeSignedTx(tx.c, P_SK, common.HexToAddress(Contracts.Registry), nil, 1000000, data)
	if err != nil {
		return err
	}

	// marshal tx into json
	js, err := SignedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	tx.SignedTx = SignedTx
	tx.JsonTx = js

	return nil
}

// add node tx
func (tx *Tx) MakeAddNodeTx(node *registry.RegistryNode) error {
	// tx data
	data := AddNodeData(node)

	log.Println("making signed add node tx")
	// Make a signed tx with data
	fmt.Println("cp sk: ", P_SK)
	SignedTx, err := MakeSignedTx(tx.c, P_SK, common.HexToAddress(Contracts.Registry), nil, 1000000, data)
	if err != nil {
		return err
	}

	// marshal tx into json
	js, err := SignedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	tx.SignedTx = SignedTx
	tx.JsonTx = js

	return nil
}

// Make tx for approving credit to market
func (tx *Tx) MakeApproveTx() error {
	// data for tx
	data := ApproveData()

	log.Println("making approve tx")
	// Make a signed tx for approve to credit
	SignedTx, err := MakeSignedTx(tx.c, U_SK, common.HexToAddress(Contracts.Credit), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := SignedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	tx.SignedTx = SignedTx
	tx.JsonTx = js

	return nil
}

// Make tx for create order
func (tx *Tx) MakeCreateOrderTx() error {
	// data for tx
	data := CreateOrderData()

	log.Println("making createorder tx")
	// Make a signed tx for createorder, sender must be user
	SignedTx, err := MakeSignedTx(tx.c, U_SK, common.HexToAddress(Contracts.Market), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := SignedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	tx.SignedTx = SignedTx
	tx.JsonTx = js

	return nil
}

// Make tx for calling registry.revise
func (tx *Tx) MakeReviseTx() error {
	// data for tx
	data := ReviseData()

	log.Println("making registry.revise tx")
	// Make a signed tx for revise, sender must be provider
	SignedTx, err := MakeSignedTx(tx.c, P_SK, common.HexToAddress(Contracts.Registry), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := SignedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	tx.SignedTx = SignedTx
	tx.JsonTx = js

	return nil
}

// Make tx for user confirm
func (tx *Tx) MakeUserConfirmTx() error {
	// data for tx
	data := UserConfirmData()

	log.Println("making user confirm tx")
	// Make a signed tx for createorder, sender must be user
	SignedTx, err := MakeSignedTx(tx.c, U_SK, common.HexToAddress(Contracts.Market), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := SignedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	tx.SignedTx = SignedTx
	tx.JsonTx = js

	return nil
}

func (tx *Tx) MakeUserCancelTx() error {
	// data for tx
	data := UserCancelData()

	log.Println("making user cancel tx")
	// Make a signed tx for createorder, sender must be user
	SignedTx, err := MakeSignedTx(tx.c, U_SK, common.HexToAddress(Contracts.Market), nil, 1000000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := SignedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	tx.SignedTx = SignedTx
	tx.JsonTx = js

	return nil
}

// send tx to chain
func (tx *Tx) Send() error {
	log.Printf("sending signed tx")

	// send the tx to client
	if err := tx.c.SendTransaction(context.Background(), tx.SignedTx); err != nil {
		fmt.Println("send tx failed:", err.Error())
		return err
	}

	// wait tx ok
	fmt.Println("waiting for tx to be ok")
	err := eth.CheckTx(tx.ep, tx.SignedTx.Hash(), "")
	if err != nil {
		fmt.Println("tx failed:", err.Error())
		return err
	}

	fmt.Println("tx ok")

	return nil
}
