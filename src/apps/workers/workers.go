package workers

import (
	"github.com/socious-io/gomq"
)

func RegisterConsumers() {
	var consumers = []gomq.AddConsumerParams{
		{
			Channel:       "sociousid/event:user.delete",
			Consumer:      gomq.NewConsumer(DeleteUser),
			IsCategorized: false,
		},
		{
			Channel:       "sociousid/event:identities.sync",
			Consumer:      gomq.NewConsumer(SyncIdentities),
			IsCategorized: false,
		},
	}

	for _, consumer := range consumers {
		gomq.AddConsumer(consumer)
	}
}
