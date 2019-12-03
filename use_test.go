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

	filename := "data/iris/iris/Iris.csv"

	useHeader := false
	trees := 100
	labelColomn := 5
	useColumn := []int{1, 2, 3, 4}
	if err := Use(useHeader, trees, labelColomn, useColumn, filename); err != nil {
		t.Fatal(err)
	}
	t.Log("success")
}
