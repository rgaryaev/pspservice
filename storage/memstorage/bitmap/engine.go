package bitmap

const (
	//PassportSeries can take values [0000; 9999]
	passportSeries = 10000
	//PassportNumber can take max value 999999
	passportMaxNumber = 1000000

	// BitSize is default bitmap window  - uint64
	bitSize = 64
	// FirstBit is first bit position
	firstBit = uint64(1 << (bitSize - 1))

	// BitmapLengthNumber is 15625
	bitmapLengthNumber = passportMaxNumber / bitSize
)

// StorageEngine represents interfaces of memory storage engine
type StorageEngine interface {
	Init() error
	AddPassport(series uint16, number uint32) (bool, error)
	CheckPassport(series uint16, number uint32) (bool, error)
}
