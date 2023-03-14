package main

import (
	"os"

	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	// определяем map подключаемых правил
	checks := map[string]bool{
		//"SA*": true,
		"S1039":  true,
		"ST1000": true,
		"QF1006": true,
	}
	var mychecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	//mychecks = append(mychecks, errcheckanalyzer.ErrCheckAnalyzer)
	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, shift.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)
	mychecks = append(mychecks, errcheck.Analyzer)

	multichecker.Main(
		mychecks...,
	)

	//var return_Url string
	//fmt.Println(return_Url)
	os.Exit(10)
	//f, _ := os.OpenFile("notes.txt", os.O_RDWR|os.O_CREATE, 0755)

}
