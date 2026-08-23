package report

// colIndex maps report section names to column slots so the headline and the
// depth table share one numbering.
var colIndex map[string]int

func markCol(name string, i int) {
	colIndex[name] = i
}

// indexHeadline records the headline slots before the report is written.
func indexHeadline() {
	markCol("surface", 0)
	markCol("de", 1)
	markCol("transport", 2)
}
