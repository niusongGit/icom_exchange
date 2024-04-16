package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var db *LevelDB

type LevelDB struct {
	db   *leveldb.DB
	path string
}

func InitDB(path string) {
	if db != nil {
		return
	}

	if path == "" {
		panic("leveldb path is empty!!!")
	}

	var err error
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}

	db = &LevelDB{
		db:   ldb,
		path: path,
	}

}

func GetDB() *LevelDB {
	return db
}

func (this *LevelDB) GetLevelDB() *leveldb.DB {
	return this.db
}

/*
保存
*/
func (this *LevelDB) Save(id []byte, bs *[]byte) error {

	//levedb保存相同的key，原来的key保存的数据不会删除，因此保存之前先删除原来的数据
	err := this.db.Delete(id, nil)
	if err != nil {
		return err
	}
	if bs == nil {
		err = this.db.Put(id, nil, nil)
	} else {
		err = this.db.Put(id, *bs, nil)
	}
	return err
}

/*
检查key是否存在
*/
func (this *LevelDB) CheckKeyExist(id []byte) bool {
	var err error
	_, err = this.Find(id)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return false
		}
		return true
	}
	return true
}

/*
查找
*/
func (this *LevelDB) Find(txId []byte) (*[]byte, error) {
	value, err := this.db.Get(txId, nil)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

/*
删除
*/
func (this *LevelDB) Remove(id []byte) error {
	return this.db.Delete(id, nil)
}

/*
查询指定前缀的key
*/
func (this *LevelDB) FindPrefixKeyAll(tag []byte) ([][]byte, [][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)
	iter := this.db.NewIterator(util.BytesPrefix(tag), nil)
	for iter.Next() {
		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())
		keys = append(keys, key)
		val := make([]byte, len(iter.Value()))
		copy(val, iter.Value())
		values = append(values, val)
	}
	iter.Release()
	err := iter.Error()
	return keys, values, err
}
