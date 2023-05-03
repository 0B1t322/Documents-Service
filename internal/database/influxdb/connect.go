package influxdb

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"
)

func Connect(serverUrl, token string) (influxdb2.Client, error) {
	client := influxdb2.NewClientWithOptions(serverUrl, token, influxdb2.DefaultOptions())
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if ok, err := client.Ping(ctx); err != nil || !ok {
		return nil, err
	}

	return client, nil
}
