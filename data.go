package main

import (
	"context"
	"fmt"

	bq "cloud.google.com/go/bigquery"
	"github.com/mchmarny/stocker/pkg/object"
	"github.com/mchmarny/stocker/pkg/util"
	"google.golang.org/api/iterator"
)

var (
	bqDataSet = util.MustEnvVar("BQ_DATSET", "")
)

func getCompanies(ctx context.Context) (symbols []*object.Company, err error) {

	logger.Println("Getting companies...")

	client, err := bq.NewClient(ctx, projectID)
	if err != nil {
		logger.Printf("Error creating BQ client: %v", err)
		return nil, err
	}

	qSQL := fmt.Sprintf("SELECT symbol, aliases FROM %s.company ORDER BY 1", bqDataSet)
	q := client.Query(qSQL)
	it, err := q.Read(ctx)
	if err != nil {
		logger.Printf("Error quering BQ: %v", err)
		return nil, err
	}

	list := make([]*object.Company, 0)

	for {
		var c object.Company
		err := it.Next(&c)
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Printf("Error looping through BQ values: %v", err)
			return nil, err
		}
		list = append(list, &c)
	}

	return list, nil

}
