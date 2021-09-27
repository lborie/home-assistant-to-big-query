package p

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/bigquery"
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type HomeAssistantMessage struct {
	State       string `json:"state"`
	LastChanged string `json:"last_changed"`
	EntityId    string `json:"entity_id"`
}

type HomeAssistantBigquery struct {
	State       float64   `bigquery:"state"`
	LastChanged time.Time `bigquery:"occured_at"`
	EntityId    string    `bigquery:"entity_id"`
}

// Main consumes a Pub/Sub message.
func Main(ctx context.Context, m PubSubMessage) error {
	logrus.SetFormatter(joonix.NewFormatter())
	logrus.Infof(string(m.Data))

	var homeAssistantMessage HomeAssistantMessage
	err := json.Unmarshal(m.Data, &homeAssistantMessage)
	if err != nil {
		return fmt.Errorf("error unmarshalling message : %w", err)
	}

	temperature, err := strconv.ParseFloat(homeAssistantMessage.State, 32)
	if err != nil {
		return fmt.Errorf("state is not a float : %w", err)
	}

	occuredAt, err := time.Parse(time.RFC3339Nano, homeAssistantMessage.LastChanged)
	if err != nil {
		return fmt.Errorf("LastChanged is not RFC3339Nano : %w", err)
	}

	projectId := os.Getenv("GCP_PROJECTID")
	if projectId == "" {
		return fmt.Errorf("environment variable GCP_PROJECTID is mandatory")
	}
	client, err := bigquery.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer func(client *bigquery.Client) {
		_ = client.Close()
	}(client)

	dataset := os.Getenv("BIGQUERY_DATASET")
	if dataset == "" {
		return fmt.Errorf("environment variable BIGQUERY_DATASET is mandatory")
	}
	table := os.Getenv("BIGQUERY_TABLE")
	if table == "" {
		return fmt.Errorf("environment variable BIGQUERY_TABLE is mandatory")
	}
	inserter := client.Dataset(dataset).Table(table).Inserter()
	items := []*HomeAssistantBigquery{
		{EntityId: homeAssistantMessage.EntityId, State: temperature, LastChanged: occuredAt},
	}
	if err := inserter.Put(ctx, items); err != nil {
		return fmt.Errorf("error when inserting into BQ : %w", err)
	}
	return nil
}
