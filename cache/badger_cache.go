package cache

import (
	"time"

	"github.com/dgraph-io/badger"
)

type BadgerCache struct {
	Conn   *badger.DB
	Prefix string
}

func (b *BadgerCache) Has(str string) (bool, error) {
	_, err := b.Get(str)
	if err == badger.ErrKeyNotFound {
		// false is same as key not found
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (b *BadgerCache) Get(str string) (interface{}, error) {
	var fromCache []byte

	err := b.Conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(str))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			fromCache = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// decode cache key
	decoded, err := decode(string(fromCache))
	if err != nil {
		return nil, err
	}

	item := decoded[str]

	return item, nil
}

func (b *BadgerCache) Set(str string, value interface{}, ttl ...int) error {
	entry := Entry{}

	entry[str] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	// check expiry
	if len(ttl) > 0 {
		err = b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded).WithTTL(time.Second * time.Duration(ttl[0]))
			err = txn.SetEntry(e)
			return err
		})
		return err
	} else {
		err = b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded)
			err = txn.SetEntry(e)
			return err
		})
		return err
	}

	return nil
}

func (b *BadgerCache) Forget(str string) error {
	err := b.Conn.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(str))
		return err
	})

	return err
}

func (b *BadgerCache) EmptyByMatch(str string) error {
	return b.emptyByMatch(str)
}

func (b *BadgerCache) Empty() error {
	return b.emptyByMatch("")
}

// drops items with no prefix or
func (b *BadgerCache) emptyByMatch(str string) error {
	// search entire badger cache for keys to delete
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := b.Conn.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}

		return nil
	}

	// how much of the cache to handle at any given time
	collectSize := 100000

	err := b.Conn.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		// create iterator with the options
		it := txn.NewIterator(opts)

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0

		for it.Seek([]byte(str)); it.ValidForPrefix([]byte(str)); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			// dont delete more than 100k items at a time because it could take a while
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					return err
				}
			}
		}

		// delete the collected keys
		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
