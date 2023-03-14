package pkg

import (
	"errors"
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
)

func GetBadgerDBInstance(DBPath string) (*badger.DB, error) {

	ops := badger.DefaultOptions(DBPath)
	ops.Logger = nil

	db, err := badger.Open(ops)

	return db, err
}

type BadgerDBWrap struct {
	db *badger.DB
}

func InitBadgerDBWrapInstance(filepath string) (*BadgerDBWrap, error) {
	db, err := GetBadgerDBInstance(filepath)
	if err != nil {
		panic(err)
	}

	instance := &BadgerDBWrap{db: db}

	return instance, nil
}

func (instance *BadgerDBWrap) Close() {
	instance.db.Close()
}

func (instance *BadgerDBWrap) Update(key string, value string) error {
	err := instance.db.Update(func(txn *badger.Txn) error {
		// e := badger.NewEntry([]byte(key), []byte(value))
		// err := txn.SetEntry(e)
		// return err
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
	return err
}

func (instance *BadgerDBWrap) UpdateByte(key string, value []byte) error {
	err := instance.db.Update(func(txn *badger.Txn) error {
		// e := badger.NewEntry([]byte(key), value)
		// err := txn.SetEntry(e)
		err := txn.Set([]byte(key), value)
		return err
	})
	return err
}

func (instance *BadgerDBWrap) BatchUpdateByte(keys []string, values [][]byte) error {
	err := instance.db.Update(func(txn *badger.Txn) error {
		for index := 0; index < len(keys); index++ {
			key := []byte(keys[index])
			err := txn.Set([]byte(key), values[index])
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (instance *BadgerDBWrap) RemoveKey(key string) error {
	err := instance.db.Update(
		func(tx *badger.Txn) error {
			_, err := tx.Get([]byte(key))
			if err != nil {
				if errors.Is(err, badger.ErrKeyNotFound) {
					return nil
				}
			}
			return tx.Delete([]byte(key))
		})

	return err
}

func (instance *BadgerDBWrap) GetByte(key string) ([]byte, error) {
	var result []byte
	err := instance.db.View(
		func(tx *badger.Txn) error {
			record, err := tx.Get([]byte(key))
			if err != nil {
				if errors.Is(err, badger.ErrKeyNotFound) {
					return nil
				}
				return err
			}

			var valCopy []byte

			err = record.Value(func(val []byte) error {
				valCopy = val
				return nil
			})

			if err != nil {
				return err
			}

			result = valCopy

			return nil
		})

	return result, err
}

func (instance *BadgerDBWrap) Get(key string) (string, error) {
	var result string
	err := instance.db.View(
		func(tx *badger.Txn) error {
			record, err := tx.Get([]byte(key))
			if err != nil {
				if errors.Is(err, badger.ErrKeyNotFound) {
					return nil
				}
				return err
			}

			var valCopy []byte

			err = record.Value(func(val []byte) error {
				valCopy = val
				return nil
			})

			if err != nil {
				return err
			}

			result = string(valCopy)

			return nil
		})

	return result, err
}

func (instance *BadgerDBWrap) BatchGet(keys []string) ([]map[string][]byte, error) {

	dataList := []map[string][]byte{}

	err := instance.db.View(
		func(tx *badger.Txn) error {
			for index := 0; index < len(keys); index++ {
				key := []byte(keys[index])
				record, err := tx.Get(key)

				if err != nil {
					if errors.Is(err, badger.ErrKeyNotFound) {
						continue
					}
					return err
				}

				var valCopy []byte

				err = record.Value(func(val []byte) error {
					valCopy = append([]byte{}, val...)
					return nil
				})

				if err != nil {
					return err
				}

				dataList = append(dataList, map[string][]byte{"key": key, "value": valCopy})
			}
			return nil
		})

	return dataList, err
}

func (instance *BadgerDBWrap) BatchSet(key string, value []byte) error {
	wb := instance.db.NewWriteBatch()
	defer wb.Cancel()
	err := wb.SetEntry(badger.NewEntry([]byte(key), value).WithMeta(0))
	if err != nil {
		fmt.Println("Failed to write data to cache.", "key", string(key), "value", string(value), "err", err)
		return err
	}
	err = wb.Flush()
	if err != nil {
		fmt.Println("Failed to flush data to cache.", "key", string(key), "value", string(value), "err", err)
	}

	return err
}
