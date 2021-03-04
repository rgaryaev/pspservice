package bitmap

import (
	"errors"

	"github.com/RoaringBitmap/roaring"
)

// RoarBitmap implements storage as a roaring bitmap in memory
type RoarBitmap struct {
	passportData []*roaring.Bitmap
}

// Init initialize internal data
func (sb *RoarBitmap) Init() error {

	sb.passportData = make([]*roaring.Bitmap, passportSeries)

	return nil
}

// AddPassport adds passport to storage
func (sb *RoarBitmap) AddPassport(series uint16, number uint32) (bool, error) {

	// create bitmap if it was not created before
	if sb.passportData[series] == nil {
		sb.passportData[series] = roaring.NewBitmap()
	}

	sb.passportData[series].Add(number)
	return true, nil

}

// CheckPassport checks do passport exist in the storage
// if exists then true will be returned
func (sb *RoarBitmap) CheckPassport(series uint16, number uint32) (bool, error) {

	return sb.passportData[series].Contains(number), nil
}

// ImportData return error because there is no standart import
func (sb *RoarBitmap) ImportData(fileName string) error {
	return errors.New("import not exists")
}
