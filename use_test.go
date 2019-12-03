package randomforest

import (
	"fmt"
	"testing"
	"time"
)

func TestUse(t *testing.T) {
	start := time.Now()
	defer func() {
		end := time.Now()
		fmt.Println("exec time: ", end.Sub(start))
	}()

	filename := "/Volumes/DailySD/SD_Desktop/data/iris/iris/Iris.csv"

	useHeader := false
	trees := 100
	labelColomn := 5
	useColumn := []int{1, 2, 3, 4}
	f, err := Use(useHeader, trees, labelColomn, useColumn, filename)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("success %f", f)
}
