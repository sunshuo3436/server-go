package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	localL2        = "https://mainnet.infura.io/v3/c0fd902e8c32475f85909278c7352314"
	dbPath         = "./block_info.db"
	flagConfigPath = "config-path"
)

var (
	configPath string // 变量用于存储配置文件路径
)

func initFlags() {
	flag.StringVar(&configPath, flagConfigPath, "", "Path to config file")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func recoverPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(w, fmt.Sprintln("Recovered from panic:", r), http.StatusInternalServerError)
	}
}

func getBlocksInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to the Ethernet client
	client, err := ethclient.Dial(localL2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Connect to a SQLitedb
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Create a table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS block_info (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		block_number BIGINT,
		block_hash TEXT,
		block_difficulty BIGINT,
		transaction_count INT
	)`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the latest 100 block information and store it in the database
	for i := 0; i < 100; i++ {
		// Get the block header
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Println("HeaderByNumber: ", err)
		}

		//get the block
		block, err := client.BlockByNumber(context.Background(), header.Number)
		if err != nil {
			log.Println("BlockByNumber: ", err)
		}

		//insert data into db
		defer recoverPanic(w)
		_, err = db.Exec(`INSERT INTO block_info (block_number, block_hash, block_difficulty, transaction_count)
			VALUES (?, ?, ?, ?)`, header.Number.String(), header.Hash().String(), block.Difficulty().Uint64(), len(block.Transactions()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//print block information
		fmt.Printf("Block %d - Hash: %s, Difficulty: %d, Transactions: %d\n", header.Number, header.Hash().String(), block.Difficulty().Uint64(), len(block.Transactions()))
	}
	fmt.Fprintln(w, "Block info saved successfully")
}

func main() {
	initFlags()

	http.HandleFunc("/blocks", getBlocksInfoHandler)
	fmt.Println("Server is running at :8080")
	http.ListenAndServe(":8080", nil)
}
