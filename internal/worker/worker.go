package worker

import (
	"avito/internal/storage/sqliteDB"
	"fmt"
)

func GetDbKey(storagePath string, h, m, l []string) *sqliteDB.Storage {
	storage, install, err := sqliteDB.CreateDB(storagePath)
	if err != nil {
		fmt.Println(err)
	}
	if install {
		for _, v := range h {
			_, err := storage.SaveHighKey(v)
			if err != nil {
				fmt.Println(err)
			}
		}
		for _, v := range m {
			storage.SaveMiddleKey(v)
		}
		for _, v := range l {
			storage.SaveLowKey(v)
		}
	}
	storage.GetHightKey()
	storage.GetMiddleKey()
	return storage
}
