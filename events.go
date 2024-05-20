package events

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type EventsEmitter struct {
	*message.Router
	pubSub *gochannel.GoChannel
}

func NewEventsEmitter() *EventsEmitter {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	router, err := message.NewRouter(message.RouterConfig{}, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatal(err)
	}

	return &EventsEmitter{
		Router: router,
		pubSub: pubSub,
	}
}

func (emitter *EventsEmitter) On(topic string, handler message.HandlerFunc) {
	emitter.AddHandler(topic, topic, emitter.pubSub, topic, emitter.pubSub, handler)
}

func (emitter *EventsEmitter) Emit(topic string, payload []byte) {
	msg := message.NewMessage(watermill.NewUUID(), payload)
	emitter.pubSub.Publish(topic, msg)
}

func (emitter *EventsEmitter) Run() {
	err := emitter.Router.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func ListenForEvents(emitter *EventsEmitter) {
	go emitter.Run()
}
