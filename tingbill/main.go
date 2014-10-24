package main

import (
	"gopkg.in/alecthomas/kingpin.v1"
	"log"
)

var (
	username    = kingpin.Flag("username", "Username for login.").Short('u').String()
	password    = kingpin.Flag("password", "Password for login.").Short('p').String()
	billCount   = kingpin.Flag("count", "Number of bills to fetch/process.").Short('c').Default("3").Uint64()
	btc         = kingpin.Flag("btc", "Show prices in mBTC").Default("false").Short('b').Bool()
	btcDiscount = kingpin.Flag("btc-discount", "Calculate mBTC with a %% discount").Default("0").Float()
)

func main() {
	kingpin.Parse()

	if *btc {
		err := GetRate()
		if err != nil {
			log.Fatalln(err)
		}
	}

	err := StartSession(Username(), Password())
	if err != nil {
		log.Fatalln(err)
	}
	defer EndSession()
	bills, err := GetBills()
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range bills {
		v.Print()
	}
}
