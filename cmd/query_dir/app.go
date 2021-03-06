package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/directory/internal/directory"
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
	var orgUnit directory.OrgUnit

	err = json.Unmarshal(byteValue, &orgUnit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "json unmarshalling error: %v\n", err)
		os.Exit(1)
	}

	dir := directory.NewDirectory(&orgUnit)

	items, err := dir.FindClosestCommonManager(*employee1, *employee2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	for _, item := range items {
		fmt.Printf("employee: \"%s\", employee: \"%s\", manager: \"%s\"\n",
			item.Employee1, item.Employee2, item.Manager)
	}
}
