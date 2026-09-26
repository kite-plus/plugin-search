//go:build wasip1

// Command plugin-search is the search plugin's module: once a site is built,
// it writes the index the browser searches.
package main

import "github.com/extism/go-pdk"

func main() {}

//go:wasmexport build_complete
func buildComplete() int32 {
	var in Input
	if err := pdk.InputJSON(&in); err != nil {
		pdk.SetError(err)
		return 1
	}
	data, err := Index(in)
	if err != nil {
		pdk.SetError(err)
		return 1
	}
	out := map[string]any{"files": []map[string]string{{"path": "index.json", "content": string(data)}}}
	if err := pdk.OutputJSON(out); err != nil {
		pdk.SetError(err)
		return 1
	}
	return 0
}
