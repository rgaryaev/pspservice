package storage

// Engine represents interfaces of memory storage engine
type Engine interface {
	Init() error
	AddPassport(series uint16, number uint32) (bool, error)
	CheckPassport(series uint16, number uint32) (bool, error)
	// ImportData  calls the engins standard import procedures if exists
	ImportData(fileName string) error
}
