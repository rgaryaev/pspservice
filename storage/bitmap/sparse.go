package bitmap

import "errors"

// BitMap defines size of bitmap window
type bitMap uint64

// SparseBitmap implements storage as a sparse bitmap in memory
//  it require 9999 * (999999/64) * 8 ~ 1.25 Gb
type SparseBitmap struct {
	passportData [][]bitMap
}

// Init initialize internal data
// matrix BitmapLengthNumber(rows) x BitmapLengthSeries (cols)
func (sb *SparseBitmap) Init() error {
	// PassportSeries is rows in the matrix

	sb.passportData = (make([][]bitMap, passportSeries))

	return nil
}

// AddPassport adds passport to storage
func (sb *SparseBitmap) AddPassport(series uint16, number uint32) (bool, error) {

	// create bitmap if it was not created before
	if sb.passportData[series] == nil {
		sb.passportData[series] = make([]bitMap, bitmapLengthNumber)
	}
	col := (number - 1) / bitSize
	colBitPosition := bitMap(firstBit >> ((number - 1) % bitSize))

	// Set required bit to 1
	sb.passportData[series][col] = sb.passportData[series][col] | colBitPosition

	return true, nil

}

// CheckPassport checks do passport exist in the storage
// if exists then true will be returned
func (sb *SparseBitmap) CheckPassport(series uint16, number uint32) (bool, error) {

	column := (number - 1) / bitSize
	columnBitPosition := bitMap(firstBit >> ((number - 1) % bitSize))

	// Set required bit to 1
	return (sb.passportData[series][column] & columnBitPosition) != 0, nil
}

// ImportData return error because there is no standart import
func (sb *SparseBitmap) ImportData(fileName string) error {
	return errors.New("import not exists")
}
