package randomforest

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"

	"gonum.org/v1/gonum/stat"

	"github.com/sajari/random-forest/RF"

	"github.com/labstack/gommon/log"
)

func Use(useHeader bool, trees, labelN int, useParam []int, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	var (
		useCol = useParam
		header []string
		data   [][]interface{}
		label  []string
	)

	r := csv.NewReader(f)
	for i := 0; ; i++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if useHeader {
			if i == 0 { // ヘッダー
				for j := range row {
					header = append(header, row[j])
				}
				continue
			}
		}

		d := make([]interface{}, len(row))
		for j := range row {
			d[j] = toF(row[j])
		}

		dd := make([]interface{}, len(useCol))
		for i, n := range useCol {
			dd[i] = d[n]
		}

		data = append(
			data,
			dd)

		// labeling
		label = append(label, row[labelN])
	}

	fmt.Printf("%+v\n", data[:2])
	fmt.Printf("%+v\n", label[:2])

	// 決定木を増やしつつ、正答率を見て最適なtrees numberを取得する
	count := trees

	var (
		checkCount = 1000
		hikakuE    = make(map[int]float64)
	)

	if len(data) < checkCount {
		checkCount = len(data)
	}

	for i := 0; i < count; i++ {
		rt := RF.DefaultForest(data, label, i)

		// 既存データで正答率
		errorCount := 0
		for j := 0; j < checkCount; j++ {
			predict := rt.Predict(data[j])
			if label[j] != predict {
				errorCount++
				continue
			}
		}

		e := float64(errorCount) / float64(checkCount)
		hikakuE[i] = e
	}

	// 出力層
	// Error
	var (
		errorPrint string
		min        = math.Inf(1)
		eMean      []float64
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { // エラーが少ないtreesNを取得
		for key, val := range hikakuE {
			if val < min {
				min = val
				errorPrint = fmt.Sprintf("tree: %d, min-error: %f", key, min)
			}
			eMean = append(eMean, val)
		}
		wg.Done()
	}()

	wg.Wait()

	mean, std := stat.MeanStdDev(eMean, nil)
	fmt.Printf("施行回数: %d, %s ----- success: %f (mean: %f, stdv: %f) \n", checkCount, errorPrint, 1-min, mean, std)

	return nil
}

func toF(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.NaN()
	}
	return f
}
