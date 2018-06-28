package main

import "testing"

// Также можно замокать testsources.go и протестировать весь handler

func Test_SortAndGetResult(t *testing.T) {
	cases := []struct {
		rows   []row
		result row
	}{
		{
			rows: []row{
				{
					Price:  3,
					Source: "test1",
				},
				{
					Price:  4,
					Source: "test2",
				},
				{
					Price:  5,
					Source: "test3",
				},
				{
					Price:  6,
					Source: "test3",
				},
				{
					Price:  2,
					Source: "test2",
				},
				{
					Price:  1,
					Source: "test1",
				},
			},
			result: row{
				Price:  5,
				Source: "test3",
			},
		},
	}

	for _, c := range cases {
		testResult := sortAndGetResult(c.rows)

		if testResult.Price != c.result.Price {
			t.Fatalf("Unexpected price. Expected %d, got %d", c.result.Price, testResult.Price)
		}

		if testResult.Source != c.result.Source {
			t.Fatalf("Unexpected Source. Expected %s, got %s", c.result.Source, testResult.Source)
		}
	}
}
