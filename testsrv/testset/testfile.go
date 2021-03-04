package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var recordList [][]string
var recordIndex map[uint32]bool

func main() {

	if len(os.Args) < 5 {
		fmt.Println("Fetch mode(r/d), Number of records,  Input and output filenames are required")
		return
	}
	fetchMode := os.Args[1]
	recordCount, _ := strconv.ParseUint(os.Args[2], 10, 32)

	inputFileName := os.Args[3]
	outputFileName := os.Args[4]

	recordIndex = make(map[uint32]bool)
	// Fill randomly
	if fetchMode == "R" || fetchMode == "r" {
		seed := time.Now().Unix()
		rand.Seed(seed)
		for {
			val := uint32(rand.Int31n(1000000))
			// if already exists then try to next
			if recordIndex[val] {
				continue
			}
			recordIndex[val] = true
			if len(recordIndex) >= int(recordCount) {
				break
			}

		}
	}

	//for val, _ := range recordIndex {
	//	fmt.Println(val)
	//}
	//
	// create new reader from csv file
	recordList = make([][]string, recordCount)

	// index of record in source file
	index := uint32(0)
	// Count of selected records for copying
	recCount := uint32(0)

	inFile, err := os.Open(inputFileName)
	if err != nil {
		panic("Inpud file is not opened!")
	}

	log.Println("Loading data...")
	csvReader := csv.NewReader(inFile)
	for {
		record, err := csvReader.Read()
		if err == io.EOF || uint64(recCount) >= recordCount {
			break
		}
		if err != nil {
			log.Fatalln("Unable to read the file with passport list", err)
			return
		}

		//  validate all data
		if _, _, ok := parseSeriesAndNumber(record[0], record[1]); !ok {
		// next if data is not valid
			continue 

		}

		// in random mode we check that index are exists

		addRec := func(step uint32) {
			if recCount%step == 0 {
				fmt.Print(".")
			}
			recordList[recCount] = record
			recCount++
		}

		if fetchMode == "R" || fetchMode == "r" {
			if recordIndex[index] {
				addRec(10)
				continue
			}

		} else {
			addRec(10)

		}
		index++
	}
	inFile.Close()

	fmt.Println("")

	// Saving
	log.Println("Saving data...")
	outFile, err := os.Create(outputFileName)
	if err != nil {
		panic("Output file is not opened!")
	}
	defer outFile.Close()
	csvWriter := csv.NewWriter(outFile)
	for i, record := range recordList {
		if err := csvWriter.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
			return
		}
		if i%10 == 0 {
			fmt.Print(".")
		}

	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("")
}

// common checks and parsing series and number
func parseSeriesAndNumber(series string, number string) (uint16, uint32, bool) {

	//  if serias not in  [0:9999] or numbers not in [1:999999]
	if !(len(series) == 4 && len(number) == 6) {
		return 0, 0, false
	}
	// if series or number are not digits
	ser, err := strconv.ParseUint(series, 10, 16)
	if err != nil {
		return 0, 0, false
	}
	num, err := strconv.ParseUint(number, 10, 32)
	if err != nil {
		return 0, 0, false
	}

	if !(ser >= 0 && ser <= 9999) || !(num >= 1 && num <= 999999) {
		return 0, 0, false
	}
	return uint16(ser), uint32(num), true
}


