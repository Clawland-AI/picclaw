package bus

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMessageBusPublishesAndConsumesInboundMessages(t *testing.T) {
	msgBus := NewMessageBus()
	defer msgBus.Close()

	expected := InboundMessage{
		Channel:    "telegram",
		SenderID:   "user-1",
		ChatID:     "chat-1",
		Content:    "temperature status",
		SessionKey: "telegram:chat-1",
		Metadata:   map[string]string{"site": "rack-a"},
	}
	msgBus.PublishInbound(expected)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, ok := msgBus.ConsumeInbound(ctx)

	if !ok {
		t.Fatal("expected inbound message to be consumed")
	}
	if got.Channel != expected.Channel || got.Content != expected.Content {
		t.Fatalf("expected inbound message %#v, got %#v", expected, got)
	}
	if got.Metadata["site"] != "rack-a" {
		t.Fatalf("expected metadata to be preserved, got %#v", got.Metadata)
	}
}

func TestMessageBusDeliversOutboundMessagesToMultipleSubscribers(t *testing.T) {
	msgBus := NewMessageBus()
	defer msgBus.Close()

	msgBus.PublishOutbound(OutboundMessage{Channel: "telegram", ChatID: "chat-1", Content: "first"})
	msgBus.PublishOutbound(OutboundMessage{Channel: "discord", ChatID: "chat-2", Content: "second"})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	first, ok := msgBus.SubscribeOutbound(ctx)
	if !ok {
		t.Fatal("expected first outbound message")
	}
	second, ok := msgBus.SubscribeOutbound(ctx)
	if !ok {
		t.Fatal("expected second outbound message")
	}

	got := map[string]string{
		first.Channel:  first.Content,
		second.Channel: second.Content,
	}
	if got["telegram"] != "first" || got["discord"] != "second" {
		t.Fatalf("expected both outbound messages to be delivered, got %#v", got)
	}
}

func TestMessageBusReturnsFalseWhenContextIsCanceled(t *testing.T) {
	msgBus := NewMessageBus()
	defer msgBus.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, ok := msgBus.ConsumeInbound(ctx); ok {
		t.Fatal("expected canceled inbound consume to return false")
	}
	if _, ok := msgBus.SubscribeOutbound(ctx); ok {
		t.Fatal("expected canceled outbound subscribe to return false")
	}
}

func TestMessageBusReturnsFalseAfterClose(t *testing.T) {
	msgBus := NewMessageBus()
	msgBus.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, ok := msgBus.ConsumeInbound(ctx); ok {
		t.Fatal("expected closed inbound channel to return false")
	}
	if _, ok := msgBus.SubscribeOutbound(ctx); ok {
		t.Fatal("expected closed outbound channel to return false")
	}
}

func TestMessageBusRegistersHandlersByChannel(t *testing.T) {
	msgBus := NewMessageBus()
	defer msgBus.Close()

	expectedErr := errors.New("handler called")
	msgBus.RegisterHandler("telegram", func(InboundMessage) error {
		return expectedErr
	})

	handler, ok := msgBus.GetHandler("telegram")
	if !ok {
		t.Fatal("expected telegram handler")
	}
	if err := handler(InboundMessage{Channel: "telegram"}); !errors.Is(err, expectedErr) {
		t.Fatalf("expected registered handler error, got %v", err)
	}
	if _, ok := msgBus.GetHandler("discord"); ok {
		t.Fatal("expected missing channel handler to be absent")
	}
}
