package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/directory/internal"
	"io/ioutil"
	"os"
)

func main() {
	file := flag.String("file", "./test/org.json", "file with organizaitonal structure")
	employee1 := flag.String("first ", "", "first employee")
	employee2 := flag.String("second ", "", "second employee")

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
	m, err := dir.FindClosestCommonManager(*employee1, *employee2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", m.Name)
}
