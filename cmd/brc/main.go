package main

import (
	"go-labs/internal/brc"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	input := filepath.Join(filepath.Dir(file), "weather_stations.csv")

	err := brc.SolutionOne(input, os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
}
