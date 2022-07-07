package db

import (
	"nomadcoin/utils"
	"sync"

	"github.com/boltdb/bolt"
)

//initialize
var db *bolt.DB
var once sync.Once

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
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
	//bolt 에는 bucket 이라는게있음 table 이 아니라
	return db
}

func SaveBlock(hash string, data []byte) {

	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockChain(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

//bolt 는 정렬기능이 없음 key value 저장소이기 때문에

func Checkpoint() []byte {
	var data []byte
	err := DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		//Get =return  slice or nil
		return nil
	})

	utils.HandleErr(err)
	return data
}

func Block(hash string) []byte {
	var data []byte
	err := DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		//bucket.Get() =return  slice or nil
		return nil
	})

	utils.HandleErr(err)
	return data
}
