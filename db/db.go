package db

import (
	"fmt"
	"nomadcoin/utils"
	"os"
	"sync"

	bolt "go.etcd.io/bbolt"
)

//initialize
var db *bolt.DB
var once sync.Once

type DB struct{}

func (DB) FindBlock(hash string) []byte {
	return findBlock(hash)
}
func (DB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
}
func (DB) SaveChain(data []byte) {
	saveChain(data)
}
func (DB) LoadChain() []byte {
	return loadChain()
}

func (DB) DeleteAllBlocks() {
	emptyBlocks()
}

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

func getDBname() string {
	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func InitDB() {
	if db == nil {

		dbPointer, err := bolt.Open(getDBname(), 0600, nil)
		utils.HandleErr(err)

		db = dbPointer
		err = db.Update(func(tx *bolt.Tx) error {
			// tx.CreateBucket()
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			//CreateBucketIfNotExists = bucket 이 존재하지 않는 경우에만 버킷 생성
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}

}

func Close() {
	db.Close()
}

//
func saveBlock(hash string, data []byte) {

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func saveChain(data []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

//bolt 는 정렬기능이 없음 key value 저장소이기 때문에

func loadChain() []byte {
	var data []byte
	//DB.View()는 read-only 트랜잭션
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		//Get =return  slice or nil
		return nil
	})

	utils.HandleErr(err)
	return data
}

func findBlock(hash string) []byte {
	var data []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		//bucket.Get() =return  slice or nil
		return nil
	})

	utils.HandleErr(err)
	return data
}

func emptyBlocks() {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		_, err = tx.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)

		return nil
	})
	utils.HandleErr(err)
}
