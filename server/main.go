package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/mattn/go-sqlite3"
)

const (
	localL2 = "https://mainnet.infura.io/v3/c0fd902e8c32475f85909278c7352314"
	dbPath  = "./block_info.db"
)

func recoverPanic() {
	if r := recover(); r != nil {
		log.Println("Recovered from panic:", r)
	}
}

func main() {
	// 连接到以太坊客户端
	client, err := ethclient.Dial(localL2)
	if err != nil {
		log.Fatal(err)
	}

	// 连接到 SQLite 数据库
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建表格
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS block_info (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		block_number BIGINT,
		block_hash TEXT,
		block_difficulty BIGINT,
		transaction_count INT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// 获取最新的100个区块信息并存储到数据库中
	for i := int64(0); i < 100; i++ {
		// 获取区块头
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Println("HeaderByNumber: ", err)
			continue
		}

		// 获取区块
		block, err := client.BlockByNumber(context.Background(), header.Number)
		if err != nil {
			log.Println("BlockByNumber: ", err)
			continue
		}

		// 插入数据到数据库
		defer recoverPanic() // 在每次迭代前调用recoverPanic函数
		_, err = db.Exec(`INSERT INTO block_info (block_number, block_hash, block_difficulty, transaction_count)
			VALUES (?, ?, ?, ?)`, header.Number.String(), header.Hash().String(), block.Difficulty().Uint64(), len(block.Transactions()))
		if err != nil {
			log.Println("Insert into block_info: ", err)
			continue
		}

		// 打印区块信息
		fmt.Printf("Block %d - Hash: %s, Difficulty: %d, Transactions: %d\n", header.Number, header.Hash().String(), block.Difficulty().Uint64(), len(block.Transactions()))
	}
}
