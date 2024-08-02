package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rockiecn/sendtx/tx"
)

func main() {
	var txType uint
	flag.UintVar(&txType, "tx", 1, "1=register 2=approve 3=createOrder 4=revise 5=userConfirm")
	chain := flag.String("chain", "local", "local:local chain, sepo:sepolia test chain")

	flag.Parse()

	fmt.Println("chain: ", *chain)

	var endpoint string
	switch *chain {
	case "local":
		endpoint = tx.Endpoint
	case "sepo":
		endpoint = tx.Endpoint2
	}

	// connect to an eth client
	log.Println("connecting client")
	c, err := ethclient.Dial(endpoint)
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

		//log.Printf("signedTx for [registcp]: \n%s\n", regTx)
		log.Printf("sending signed register tx")
		ctx := context.Background()
		c.SendTransaction(ctx, regTx)

	case 2:
		// tx for send to chain directly
		apprTx, err := makeApproveTx(c)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("signedTx for [approve] credit to market: \n%s\n", apprTx)
		log.Printf("sending signed approve tx")
		ctx := context.Background()
		c.SendTransaction(ctx, apprTx)
	case 3:
		// signed market.createorder tx for send to chain directly
		coTx, err := makeCreateOrderTx(c)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("signedTx for [createorder]: \n%s\n", coTx)
		log.Printf("sending signed create order tx")
		ctx := context.Background()
		c.SendTransaction(ctx, coTx)
	case 4:
		// signed registry.revise tx for send to chain directly
		reviseTx, err := makeReviseTx(c)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("signedTx for [revise]: \n%s\n", reviseTx)
		log.Printf("sending signed revise tx")
		ctx := context.Background()
		c.SendTransaction(ctx, reviseTx)
	case 5:
		// signed market.userconfirm tx for send to chain directly
		confirmTx, err := makeUserConfirmTx(c)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("signedTx for [userconfirm]: \n%s\n", confirmTx)
		log.Printf("sending signed user confirm tx")
		ctx := context.Background()
		c.SendTransaction(ctx, confirmTx)
	case 6:
		// signed market.userconfirm tx for send to chain directly
		cancelTx, err := makeUserCancelTx(c)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("signedTx for [usercancel]: \n%s\n", cancelTx)
		log.Printf("sending signed user cancel tx")
		ctx := context.Background()
		c.SendTransaction(ctx, cancelTx)
	}
}

// make tx for register cp
func makeRegisterTx(c *ethclient.Client) (*types.Transaction, error) {
	// tx data
	data := tx.RegisterData()

	log.Println("making signed register tx")
	// make a signed tx
	fmt.Println("cp sk: ", tx.P_SK)
	signedTx, err := tx.MakeSignedTx(c, tx.P_SK, common.HexToAddress(tx.Contracts.Registry), nil, 400000, data)
	if err != nil {
		return nil, err
	}

	// // marshal tx into json
	// js, err := signedTx.MarshalJSON()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return signedTx, nil
}

// make tx for approving credit to market
func makeApproveTx(c *ethclient.Client) (*types.Transaction, error) {
	// data for tx
	data := tx.ApproveData()

	log.Println("making approve tx")
	// make a signed tx for approve to credit
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Credit), nil, 400000, data)
	if err != nil {
		log.Fatal(err)
	}

	// // marshal tx into json
	// js, err := signedTx.MarshalJSON()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return signedTx, nil
}

// make tx for create order
func makeCreateOrderTx(c *ethclient.Client) (*types.Transaction, error) {
	// data for tx
	data := tx.CreateOrderData()

	log.Println("making createorder tx")
	// make a signed tx for createorder, sender must be user
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Market), nil, 400000, data)
	if err != nil {
		log.Fatal(err)
	}

	// // marshal tx into json
	// js, err := signedTx.MarshalJSON()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return signedTx, nil
}

// make tx for calling registry.revise
func makeReviseTx(c *ethclient.Client) (*types.Transaction, error) {
	// data for tx
	data := tx.ReviseData()

	log.Println("making registry.revise tx")
	// make a signed tx for revise, sender must be provider
	signedTx, err := tx.MakeSignedTx(c, tx.P_SK, common.HexToAddress(tx.Contracts.Registry), nil, 400000, data)
	if err != nil {
		log.Fatal(err)
	}

	// // marshal tx into json
	// js, err := signedTx.MarshalJSON()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return signedTx, nil
}

// make tx for user confirm
func makeUserConfirmTx(c *ethclient.Client) (*types.Transaction, error) {
	// data for tx
	data := tx.UserConfirmData()

	log.Println("making user confirm tx")
	// make a signed tx for createorder, sender must be user
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Market), nil, 400000, data)
	if err != nil {
		log.Fatal(err)
	}

	// // marshal tx into json
	// js, err := signedTx.MarshalJSON()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return signedTx, nil
}

func makeUserCancelTx(c *ethclient.Client) (*types.Transaction, error) {
	// data for tx
	data := tx.UserCancelData()

	log.Println("making user cancel tx")
	// make a signed tx for createorder, sender must be user
	signedTx, err := tx.MakeSignedTx(c, tx.U_SK, common.HexToAddress(tx.Contracts.Market), nil, 400000, data)
	if err != nil {
		log.Fatal(err)
	}

	// // marshal tx into json
	// js, err := signedTx.MarshalJSON()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return signedTx, nil
}
