package bitmap

import (
	"github.com/RoaringBitmap/roaring"
)

// RoarBitmap implements storage as a roaring bitmap in memory
type RoarBitmap struct {
	passportData []*roaring.Bitmap
}

// Init initialize internal data
func (sb *RoarBitmap) Init() error {

	sb.passportData = make([]*roaring.Bitmap, passportSeries)

	//  Pasport bumber will be stored in the roaring bitmap
	for i := range sb.passportData {
		sb.passportData[i] = roaring.NewBitmap()
	}
	return nil
}

// AddPassport adds passport to storage
func (sb *RoarBitmap) AddPassport(series uint16, number uint32) (bool, error) {

	sb.passportData[series].Add(number)
	return true, nil

}

// CheckPassport checks do passport exist in the storage
// if exists then true will be returned
func (sb *RoarBitmap) CheckPassport(series uint16, number uint32) (bool, error) {

	return sb.passportData[series].Contains(number), nil
}
