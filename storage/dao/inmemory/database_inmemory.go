package inmemory

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/allegro/bigcache"
	"github.com/denouche/plex-watcher/storage/dao"
	"github.com/denouche/plex-watcher/storage/model"
	"github.com/denouche/plex-watcher/utils"
)

type DatabaseInMemory struct {
	Cache *bigcache.BigCache
	file  *os.File
}

type DatabaseInMemoryExport struct {
	Libraries []*model.Library
}

func NewDatabaseInMemory(file string) dao.Database {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	if err != nil {
		utils.GetLogger().WithError(err).Fatal("Error while instantiate cache")
	}

	result := &DatabaseInMemory{
		Cache: cache,
	}

	if file != "" {
		f, err := os.OpenFile(file, os.O_RDWR, 0666)
		if err != nil {
			utils.GetLogger().WithError(err).Error("error while opening file for database in memory")
		} else {
			result.file = f

			content, err := ioutil.ReadAll(result.file)
			if err != nil {
				utils.GetLogger().WithError(err).Error("error while reading file for database in memory")
			}

			export := &DatabaseInMemoryExport{}
			err = json.Unmarshal(content, export)
			if err != nil {
				utils.GetLogger().WithError(err).Error("error while unmarshalling for database in memory")
			}

			result.fromExport(export)
		}
	}

	return result
}

func (db *DatabaseInMemory) save(key string, data []interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		utils.GetLogger().WithError(err).WithField("key", key).Error("Error while marshal fake")
		db.Cache.Set(key, []byte("[]"))
		return
	}
	err = db.Cache.Set(key, b)
	if err != nil {
		utils.GetLogger().WithError(err).WithField("key", key).Error("Error while saving fake")
	}

	if db.file != nil {
		export := db.toExport()
		exportMarshalled, _ := json.Marshal(export)
		_ = db.file.Truncate(0)
		_, _ = db.file.Seek(0, 0)
		_, err := db.file.Write(exportMarshalled)
		if err != nil {
			utils.GetLogger().WithError(err).Error("error while writing to file")
		}
		err = db.file.Sync()
		if err != nil {
			utils.GetLogger().WithError(err).Error("error while sync file")
		}
	}
}

func (db *DatabaseInMemory) toExport() *DatabaseInMemoryExport {
	return &DatabaseInMemoryExport{
		Libraries: db.loadLibraries(),
	}
}

func (db *DatabaseInMemory) fromExport(e *DatabaseInMemoryExport) {
	db.saveLibraries(e.Libraries)
}
