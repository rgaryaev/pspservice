package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/rgaryaev/pspservice/config"
	"github.com/rgaryaev/pspservice/lstnr"
	"github.com/rgaryaev/pspservice/storage"
)

func main() {

	var cfg *config.Configuration

	printMemStats()
	//configData
	cfg, err := config.LoadConfiguration()
	if err != nil {
		panic("Config file is not loaded! " + err.Error())
	}

	var storage *storage.Storage = new(storage.Storage)

	log.Println("Initialize and start passport data storage...")
	err = (*storage).StartStorage(cfg)
	if err != nil {
		panic("Passport storage is not initialized : " + err.Error())
	}
	runtime.GC()
	printMemStats()

	log.Println("Starting storage update scheduler ...")
	go func() {
		err = (*storage).UpdateStorage(cfg)
		if err != nil {
			panic("Storage update scheduler error: " + err.Error())
		}
	}()
	log.Println("storage updater has been started")
	//
	log.Println("Starting http listener...")
	go func() {
		err := lstnr.StartListener(cfg, storage)
		if err != nil {
			panic("Start http listener: " + err.Error())
		}
	}()
	log.Println("Http listener has been started")

	//
	select {}

}

func printMemStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Println("mem.Sys(Mb) : ", bToMb(mem.Sys))
	fmt.Println("mem.Alloc(Mb): ", bToMb(mem.Alloc))
	fmt.Println("mem.TotalAlloc(Mb): ", bToMb(mem.TotalAlloc))
	fmt.Println("mem.HeapAlloc(Mb):  ", bToMb(mem.HeapAlloc))

}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
