package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
)

type row struct {
	kota  string
	count int
}

func main() {
	f, err := os.Open("data/pendaftaran.csv")
	if err != nil {
		fmt.Println("gagal membuka file:", err)
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, _ := r.ReadAll()

	tally := map[string]int{}
	for _, rec := range records[1:] { // skip header
		if len(rec) >= 4 {
			tally[rec[3]]++
		}
	}

	rows := []row{}
	for k, v := range tally {
		rows = append(rows, row{k, v})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].count > rows[j].count
	})

	fmt.Printf("%-15s %s\n", "Kota", "Pendaftar")
	fmt.Println("-------------------------")
	for _, r := range rows {
		fmt.Printf("%-15s %d\n", r.kota, r.count)
	}
}
