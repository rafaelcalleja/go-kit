package transaction

import (
	"context"
)

type BrokerInitializer struct {
	broker      *Broker
	initializer Initializer
	tx          Transaction
}

type BrokerConnClient struct {
	Id int
	Tx Transaction
}

func NewBrokerInitializer(initializer Initializer, broker *Broker) *BrokerInitializer {
	return &BrokerInitializer{
		broker:      broker,
		initializer: initializer,
	}
}

func (c *BrokerInitializer) Begin(ctx context.Context) (Transaction, error) {
	tx, err := c.initializer.Begin(ctx)
	c.broker.Publish(&BrokerConnClient{Tx: tx})
	c.tx = tx

	return c, err
}

func (c *BrokerInitializer) Commit() error {
	defer c.broker.Publish(nil)
	return c.tx.Commit()
}

func (c *BrokerInitializer) Rollback() error {
	defer c.broker.Publish(nil)
	return c.tx.Rollback()
}
