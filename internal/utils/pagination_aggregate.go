package utils

import "go.mongodb.org/mongo-driver/bson"

func PaginationAggregate(page, limit int) bson.M {
	skip := limit * (page - 1)

	return bson.M{
		"metadata": []bson.M{{
			"$count": "totalData",
		}, {
			"$project": bson.M{
				"totalData": 1,
				"totalPage": bson.M{
					"$toInt": bson.M{
						"$ceil": bson.M{
							"$divide": []interface{}{"$totalData", limit},
						},
					},
				},
				"previousPage": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$lte": []interface{}{page, 1}},
						"then": nil,
						"else": bson.M{"$subtract": []interface{}{page, 1}},
					},
				},
				"currentPage": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": []interface{}{page, 1}},
						"then": 1,
						"else": bson.M{"$toInt": bson.M{"$ceil": bson.M{"$divide": []interface{}{page, 1}}}},
					},
				},
				"nextPage": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$lte": []interface{}{
								bson.M{"$add": []interface{}{page, 1}},
								bson.M{"$toInt": bson.M{"$ceil": bson.M{"$divide": []interface{}{"$totalData", limit}}}},
							},
						},
						"then": bson.M{"$add": []interface{}{page, 1}},
						"else": nil,
					},
				},
			},
		}},
		"data": []bson.M{
			{"$skip": skip},
			{"$limit": limit},
		},
	}
}
