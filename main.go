package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

func e(err error) {
	if err != nil {
		panic(err)
	}
}

type MohawkCSV struct {
	Invoice            string    `csv:"Invoice"`
	ShipDate           time.Time `csv:"ShipDate"`
	Contract           string    `csv:"Contract"`
	CustNbr            string    `csv:"Cust#"`
	ProductCode        string    `csv:"ProductCode"`
	ProductDescription string    `csv:"ProductDescription"`
	Qty                int       `csv:"Qty"`
	Unit               string    `csv:"Unit"`
	ContCost           float32   `csv:"ContCost"`
	IntoStock          float32   `csv:"IntoStock"`
	RebateAmt          float32   `csv:"RebateAmt"`
	CustName           string    `csv:"CustName"`
	CustAddress1       string    `csv:"CustAddress1"`
	CustAddress2       string    `csv:"CustAddress2"`
	CustCity           string    `csv:"CustCity"`
	ST                 string    `csv:"ST"`
	ZipCd              string    `csv:"ZipCd"`
}

func main() {
	filePath := flag.String("fp", "", "file path")

	flag.Parse()

	if *filePath == "" {
		flag.Usage()
		return
	}

	// do something
	fs, err := os.Stat(*filePath)
	e(err)

	fmt.Printf("%+v\n", fs)

	src, err := os.OpenFile(*filePath, os.O_RDONLY, 0666)
	e(err)

	dst, err := os.OpenFile(*filePath+".new", os.O_WRONLY|os.O_CREATE, 0666)
	e(err)

	trim := regexp.MustCompile(`(^\s+|\s+$)`)
	moneyChars := regexp.MustCompile((`[\$\"\=]`))
	// regexp find dates 06/01/2011
	dateRegexp := regexp.MustCompile(`(\d{2}\/\d{2}\/\d{4})`)
	dateFormat := "01/02/2006"
	isoDateFormat := "2006-01-02T15:04:05Z07:00"

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, "                                   ", "", -1)
		line = trim.ReplaceAllString(line, "")
		line = moneyChars.ReplaceAllString(line, "")
		line = dateRegexp.ReplaceAllStringFunc(line, func(date string) string {
			t, err := time.Parse(dateFormat, date)
			e(err)
			return t.Format(isoDateFormat)
		})
		dst.WriteString(line + "\n")
	}

	dst.Close()

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	newSrc, err := os.OpenFile(*filePath+".new", os.O_RDONLY, 0666)
	e(err)
	defer newSrc.Close()

	rebates := []MohawkCSV{}

	if err := gocsv.UnmarshalFile(newSrc, &rebates); err != nil {
		panic(err)
	}

	for _, r := range rebates {
		fmt.Printf("%+v\n", r)
	}

}
