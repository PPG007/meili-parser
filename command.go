package main

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/cobra"
)

var (
	SearchCommand *cobra.Command
	indexPage     string
	meiliHost     string
	uid           string
	apiKey        string
)

func init() {
	SearchCommand = &cobra.Command{
		Use:     "meiliParser [options] [values]",
		Example: "meiliParser --indexPage https://example.com --meiliHost https://meilihost.com:7700 --uid testName --apiKey testKey",
		RunE: func(cmd *cobra.Command, args []string) error {
			return start()
		},
	}
	SearchCommand.Flags().StringVar(&indexPage, "indexPage", "", "Your static website index page")
	SearchCommand.Flags().StringVar(&meiliHost, "meiliHost", "", "Your MeiliSearch instance host")
	SearchCommand.Flags().StringVar(&uid, "uid", "", "Target MeiliSearch uid")
	SearchCommand.Flags().StringVar(&apiKey, "apiKey", "", "Your static website index page")
	SearchCommand.MarkFlagRequired("indexPage")
	SearchCommand.MarkFlagRequired("meiliHost")
	SearchCommand.MarkFlagRequired("uid")
	SearchCommand.MarkFlagRequired("apiKey")
}

func start() error {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   meiliHost,
		APIKey: apiKey,
	})
	_, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        uid,
		PrimaryKey: "id",
	})
	if err != nil {
		return err
	}
	index, err := client.GetIndex(uid)
	if err != nil {
		return err
	}
	parseFromStartPage(indexPage, index)
	logger.Info("Parse end")
	return nil
}
