package pilosadb

import (
	"bytes"
	"io/ioutil"
	"log"
	"sort"

	"github.com/pilosa/go-pilosa"
	"github.com/pilosa/go-pilosa/csv"
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
		log.Panic("pilosa schema was not created: " + err.Error())
	}
	// Create an Index object
	pdb.myIndex = schema.Index(indexName)

	pdb.seriesField = pdb.myIndex.Field(fieldName)

	err = pdb.client.SyncSchema(schema)
	if err != nil {
		log.Panic("Pilisa schema was not synchronized: " + err.Error())
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
		log.Println("error: inset into pilosa db:  " + err.Error())
	}

	return true, nil

}

// CheckPassport is pasport in the Pilosa
func (pdb *PilosaDB) CheckPassport(series uint16, number uint32) (bool, error) {

	response, err := pdb.client.Query(pdb.seriesField.RowsColumn(number))
	if err != nil {
		return false, err
	}
	seriesRows := response.Result().RowIdentifiers().IDs

	// we will rely on the fact that pilosa returns sorted slices
	// so we can use go's sort.Search function
	res := searchUints64(seriesRows, uint64(series))
	if res < len(seriesRows) {
		return seriesRows[res] == uint64(series), nil
	}
	/*
		// make search via service row with number 10000
		var serviceRow int = 10000
		var err error
		//	_, err := pdb.client.Query(pdb.seriesField.Set(serviceRow, int(number)))
		_, err = pdb.client.Query(pdb.seriesField.ClearRow(serviceRow))
		_, err = pdb.client.Query(pdb.seriesField.Set(serviceRow, int(number)))
		response, err := pdb.client.Query(pdb.myIndex.Intersect(pdb.seriesField.Row(int(series)), pdb.seriesField.Row(serviceRow)))
		if err != nil {
			return false, err
		}
		//log.Println(response.Result().Row().Columns)

		if len(response.Result().Row().Columns) > 0 && response.Result().Row().Columns[0] == uint64(number) {
			return true, nil
		}
	*/
	return false, nil
}

// ImportData import data using standard import
func (pdb *PilosaDB) ImportData(fileName string) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("pilosa import file is not opened")
		return err
	}
	log.Println("pilosa started import...")
	iterator := csv.NewColumnIterator(csv.RowIDColumnID, bytes.NewReader(file))
	err = pdb.client.ImportField(pdb.seriesField, iterator)
	if err != nil {
		log.Println("error during pilosa import: " + err.Error())
		return err
	}
	log.Println("pilosa finished import...")
	return nil
}
