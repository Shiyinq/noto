package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestPaginationAggregate(t *testing.T) {
	testCases := []struct {
		name     string
		page     int
		limit    int
		expected bson.M
	}{
		{
			name:  "Test page 1 with limit 10",
			page:  1,
			limit: 10,
			expected: bson.M{
				"metadata": []bson.M{
					{"$count": "totalData"},
					{
						"$project": bson.M{
							"totalData": 1,
							"totalPage": bson.M{
								"$toInt": bson.M{
									"$ceil": bson.M{
										"$divide": []interface{}{"$totalData", 10},
									},
								},
							},
							"previousPage": bson.M{
								"$cond": bson.M{
									"if":   bson.M{"$lte": []interface{}{1, 1}},
									"then": nil,
									"else": bson.M{"$subtract": []interface{}{1, 1}},
								},
							},
							"currentPage": bson.M{
								"$cond": bson.M{
									"if":   bson.M{"$eq": []interface{}{1, 1}},
									"then": 1,
									"else": bson.M{"$toInt": bson.M{"$ceil": bson.M{"$divide": []interface{}{1, 1}}}},
								},
							},
							"nextPage": bson.M{
								"$cond": bson.M{
									"if": bson.M{
										"$lte": []interface{}{
											bson.M{"$add": []interface{}{1, 1}},
											bson.M{"$toInt": bson.M{"$ceil": bson.M{"$divide": []interface{}{"$totalData", 10}}}},
										},
									},
									"then": bson.M{"$add": []interface{}{1, 1}},
									"else": nil,
								},
							},
						},
					},
				},
				"data": []bson.M{
					{"$skip": 0},
					{"$limit": 10},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := PaginationAggregate(testCase.page, testCase.limit)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
