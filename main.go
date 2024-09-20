package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/rockiecn/sendtx/tx"

	"github.com/grid/contracts/eth/contracts"
)

func main() {
	var txType uint
	flag.UintVar(&txType, "tx", 1, "1=register 2=approve 3=createOrder 4=revise 5=userConfirm 6=userCancel")
	chain := flag.String("chain", "local", "local:local chain, sepo:sepolia test chain")
	auto := flag.Bool("auto", false, "auto send the tx to chain")

	flag.Parse()

	fmt.Println("chain: ", *chain)

	fmt.Println("type:", txType)

	var endpoint string
	switch *chain {
	case "local":
		endpoint = tx.Endpoint

		// load contracts
		local := contracts.Local{}
		local.LoadPath("../grid-contracts/eth/contracts/local.json")
		tx.Contracts = local.Contracts

		fmt.Println("contract addresses:", tx.Contracts)
	case "sepo":
		endpoint = tx.Endpoint2

		// load contracts
		sepo := contracts.Sepo{}
		sepo.LoadPath("../grid-contracts/eth/contracts/sepo.json")
		tx.Contracts = sepo.Contracts

		fmt.Println("contract addresses:", tx.Contracts)
	}

	txObj := tx.NewTx(endpoint)

	switch txType {
	case 1:
		// signed register tx for send to chain directly
		err := txObj.MakeRegisterTx()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("signedTx for [registcp]: \n%s\n", txObj.JsonTx)

		if *auto {
			err := txObj.Send()
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Printf("add 2 nodes for this cp")

		// add node1
		node1, err := tx.NewNode()
		if err != nil {
			panic(err)
		}
		err = txObj.MakeAddNodeTx(node1)
		if err != nil {
			log.Fatal(err)
		}

		// send tx
		if *auto {
			err := txObj.Send()
			if err != nil {
				log.Fatal(err)
			}
		}

		// add node2
		node2, err := tx.NewNode2()
		if err != nil {
			panic(err)
		}
		err = txObj.MakeAddNodeTx(node2)
		if err != nil {
			log.Fatal(err)
		}

		// send tx
		if *auto {
			err := txObj.Send()
			if err != nil {
				log.Fatal(err)
			}
		}

	case 2:
		// tx for send to chain directly
		err := txObj.MakeApproveTx()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [approve] credit to market: \n%s\n", txObj.JsonTx)

		if *auto {
			err := txObj.Send()
			if err != nil {
				log.Fatal(err)
			}
		}
	case 3:
		// signed market.createorder tx for send to chain directly
		err := txObj.MakeCreateOrderTx()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [createorder]: \n%s\n", txObj.JsonTx)

		if *auto {
			err := txObj.Send()
			if err != nil {
				log.Fatal(err)
			}
		}
	case 4:
		// signed registry.revise tx for send to chain directly
		err := txObj.MakeReviseTx()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [revise]: \n%s\n", txObj.JsonTx)

		if *auto {
			err := txObj.Send()
			if err != nil {
				log.Fatal(err)
			}
		}

	case 5:
		// signed market.userconfirm tx for send to chain directly
		err := txObj.MakeUserConfirmTx()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [userconfirm]: \n%s\n", txObj.JsonTx)

		if *auto {
			err := txObj.Send()
			if err != nil {
				log.Fatal(err)
			}
		}

	case 6:
		// signed market.userconfirm tx for send to chain directly
		err := txObj.MakeUserCancelTx()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("signedTx for [usercancel]: \n%s\n", txObj.JsonTx)

		if *auto {
			err := txObj.Send()
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}
