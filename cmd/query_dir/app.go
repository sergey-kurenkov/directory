package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/directory/internal"
)

func main() {
	file := flag.String("file", "./test/org.json", "file with organizaitonal structure")
	employee1 := flag.String("first", "", "first employee")
	employee2 := flag.String("second", "", "second employee")

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()

	if employee1 == nil || employee2 == nil {
		fmt.Fprintf(os.Stderr, "no employees\n")
		os.Exit(1)
	}

	jsonFile, err := os.Open(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open file error: %v\n", err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var orgUnit internal.OrgUnit
	json.Unmarshal(byteValue, &orgUnit)

	dir := internal.NewDirectory(&orgUnit)

	commonManagers := dir.FindClosestCommonManager(*employee1, *employee2)
	if len(commonManagers) == 0 {
		fmt.Fprintf(os.Stderr, "no common manager\n")
		os.Exit(1)
	}
	for _, commonManager := range commonManagers {
		fmt.Printf("employee #1: %s, employee #2: %s, common manager: %s\n",
			commonManager.Employee1, commonManager.Employee2, commonManager.Manager)
	}
}
