package main

import (
	"flag"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rockiecn/sendtx/tx"
)

func main() {
	var txType uint
	flag.UintVar(&txType, "tx", 1, "1=register 2=approve 3=createOrder 4=revise")
	flag.Parse()

	// connect to an eth client
	log.Println("connecting client")
	c, err := ethclient.Dial(tx.Endpoint)
	if err != nil {
		log.Fatal(err)
	}

	switch txType {
	case 1:
		// signed register tx for send to chain directly
		regTx, err := makeRegisterTx(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [registcp]: \n%s\n", regTx)
	case 2:
		// tx for send to chain directly
		apprTx, err := makeApproveTx(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [approve] credit to market: \n%s\n", apprTx)
	case 3:
		// signed market.createorder tx for send to chain directly
		coTx, err := makeCreateOrderTx(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [createorder]: \n%s\n", coTx)
	case 4:
		// signed registry.update tx for send to chain directly
		reviseTx, err := makeReviseTx(c)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [revise]: \n%s\n", reviseTx)
	}
}

// make tx for register cp
func makeRegisterTx(c *ethclient.Client) ([]byte, error) {
	// tx data
	data := tx.RegisterData()

	log.Println("making signed register tx")
	// make a signed tx
	signedTx, err := tx.MakeSignedTx(c, tx.P_SK, common.HexToAddress(tx.Contracts.Registry), nil, 400000, data)
	if err != nil {
		return nil, err
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return js, nil
}

// make tx for approving credit to market
func makeApproveTx(c *ethclient.Client) ([]byte, error) {
	// data for tx
	data := tx.ApproveData()

	log.Println("making approve tx")
	// make a signed tx for approve to credit
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Credit), nil, 400000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return js, nil
}

// make tx for create order
func makeCreateOrderTx(c *ethclient.Client) ([]byte, error) {
	// data for tx
	data := tx.CreateOrderData()

	log.Println("making createorder tx")
	// make a signed tx for createorder, sender must be user
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Market), nil, 400000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return js, nil
}

// make tx for calling registry.revise
func makeReviseTx(c *ethclient.Client) ([]byte, error) {
	// data for tx
	data := tx.ReviseData()

	log.Println("making registry.revise tx")
	// make a signed tx for revise, sender must be provider
	signedTx, err := tx.MakeSignedTx(c, tx.P_SK, common.HexToAddress(tx.Contracts.Registry), nil, 400000, data)
	if err != nil {
		log.Fatal(err)
	}

	// marshal tx into json
	js, err := signedTx.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	return js, nil
}
