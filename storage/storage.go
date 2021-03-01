package storage

import (
	"github.com/rgaryaev/pspservice/config"
)

type storageData struct {
	StorageType string
	StorageDate string
}

// Storage - set of common interfaces of a passport data storage
type Storage interface {
	StartStorage(cfg *config.Configuration) error
	CheckStorage(cfg *config.Configuration) (bool, error)
	UpdateStorage(cfg *config.Configuration) error
	IsPassportInList(series string, number string) (bool, error)
}
