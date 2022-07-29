package blockchain

import (
	"nomadcoin/utils"
	"reflect"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = FakeDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an instace of a block")
	}

}

func TestFindBlock(t *testing.T) {
	t.Run("Block not found", func(t *testing.T) {
		dbStorage = FakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 1,
				}
				return utils.ToBytes(b)
			},
		}
		block, _ := FindBlock("Xx")
		if reflect.TypeOf(block) != reflect.TypeOf(&Block{}) {
			t.Error("Block should be found.")
		}

	})

}
