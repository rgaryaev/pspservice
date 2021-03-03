package pilosadb

import (
	"log"
	"sort"

	"github.com/pilosa/go-pilosa"
)

// PilosaDB as a storage engine
type PilosaDB struct {
	client      *pilosa.Client
	schema      *pilosa.Schema
	myIndex     *pilosa.Index
	seriesField *pilosa.Field
}

const indexName = "passport"
const fieldName = "series"

// Init initialize Pilosa DB
// default connecion on port 10101
func (pdb *PilosaDB) Init() error {
	pdb.client = pilosa.DefaultClient()

	// Retrieve the schema
	schema, err := pdb.client.Schema()
	if err != nil {
		log.Println("pilosa schema was not created: " + err.Error())
		return err
	}
	// Create an Index object
	pdb.myIndex = schema.Index(indexName)

	pdb.seriesField = pdb.myIndex.Field(fieldName)

	err = pdb.client.SyncSchema(schema)
	if err != nil {
		log.Println("Pilisa schema was not synchronized: " + err.Error())
		return err
	}
	return nil
}

func searchUints64(a []uint64, x uint64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// AddPassport adds passport to Pilosa
func (pdb *PilosaDB) AddPassport(series uint16, number uint32) (bool, error) {

	_, err := pdb.client.Query(pdb.seriesField.Set(int(series), int(number)))

	if err != nil {
		log.Fatal(err)
	}

	return true, nil

}

// CheckPassport is pasport in the Pilosa
func (pdb *PilosaDB) CheckPassport(series uint16, number uint32) (bool, error) {

	response, err := pdb.client.Query(pdb.seriesField.RowsColumn(997710))
	if err != nil {
		return false, err
	}
	seriesRows := response.Result().RowIdentifiers().IDs
	// we are relying  on the fact that Pilosa returns sorted slices
	// so we can use go's sort.Search function
	res := searchUints64(seriesRows, uint64(series))

	if res < len(seriesRows) {
		return seriesRows[res] == uint64(series), nil
	}
	return false, nil
}
