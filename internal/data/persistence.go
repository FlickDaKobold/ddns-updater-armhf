package data

import (
	"fmt"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/records"
)

func (db *Database) Update(id uint, record records.Record) (err error) {
	db.Lock()
	defer db.Unlock()
	if id > uint(len(db.data))-1 {
		return fmt.Errorf("%w: for id %d", ErrRecordNotFound, id)
	}
	currentCount := len(db.data[id].History)
	newCount := len(record.History)
	db.data[id] = record
	// new IP address added
	if newCount > currentCount {
		if err := db.persistentDB.StoreNewIP(
			record.Provider.Domain(),
			record.Provider.Owner(),
			record.History.GetCurrentIP(),
			record.History.GetSuccessTime(),
		); err != nil {
			return err
		}
	}
	return nil
}
