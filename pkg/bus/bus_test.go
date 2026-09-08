package bus

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestInboundRoundTrip(t *testing.T) {
	b := NewMessageBus()
	defer b.Close()
	want := InboundMessage{Channel: "chat", SenderID: "sender", ChatID: "room", Content: "hello", Media: []string{"image"}, SessionKey: "session", Metadata: map[string]string{"key": "value"}}
	b.PublishInbound(want)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, ok := b.ConsumeInbound(ctx)
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, %v; want %#v", got, ok, want)
	}
}

func TestOutboundRoundTrip(t *testing.T) {
	b := NewMessageBus()
	defer b.Close()
	want := OutboundMessage{Channel: "chat", ChatID: "room", Content: "hello"}
	b.PublishOutbound(want)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, ok := b.SubscribeOutbound(ctx)
	if !ok || got != want {
		t.Fatalf("got %#v, %v; want %#v", got, ok, want)
	}
}

func TestCanceledEmptyQueues(t *testing.T) {
	b := NewMessageBus()
	defer b.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if msg, ok := b.ConsumeInbound(ctx); ok || !reflect.DeepEqual(msg, InboundMessage{}) {
		t.Fatalf("inbound: %#v, %v", msg, ok)
	}
	if msg, ok := b.SubscribeOutbound(ctx); ok || msg != (OutboundMessage{}) {
		t.Fatalf("outbound: %#v, %v", msg, ok)
	}
}

func TestConcurrentOutboundConsumers(t *testing.T) {
	b := NewMessageBus()
	defer b.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	const total = 100
	var wg sync.WaitGroup
	results := make(chan OutboundMessage, total)
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < total/4; i++ {
				msg, ok := b.SubscribeOutbound(ctx)
				if !ok {
					return
				}
				results <- msg
			}
		}()
	}
	for i := 0; i < total; i++ {
		b.PublishOutbound(OutboundMessage{Content: fmt.Sprint(i)})
	}
	wg.Wait()
	close(results)
	seen := map[string]int{}
	for msg := range results {
		seen[msg.Content]++
	}
	for i := 0; i < total; i++ {
		if seen[fmt.Sprint(i)] != 1 {
			t.Fatalf("message %d delivered %d times", i, seen[fmt.Sprint(i)])
		}
	}
}

func TestHandlersAreSelectedByChannel(t *testing.T) {
	b := NewMessageBus()
	defer b.Close()
	if _, ok := b.GetHandler("missing"); ok {
		t.Fatal("unexpected handler")
	}
	calls := map[string]int{}
	b.RegisterHandler("a", func(InboundMessage) error { calls["old-a"]++; return nil })
	b.RegisterHandler("b", func(InboundMessage) error { calls["b"]++; return nil })
	b.RegisterHandler("a", func(InboundMessage) error { calls["new-a"]++; return nil })
	for _, channel := range []string{"a", "b"} {
		h, ok := b.GetHandler(channel)
		if !ok {
			t.Fatalf("missing handler %s", channel)
		}
		if err := h(InboundMessage{Channel: channel}); err != nil {
			t.Fatal(err)
		}
	}
	if calls["old-a"] != 0 || calls["new-a"] != 1 || calls["b"] != 1 {
		t.Fatalf("wrong handlers called: %v", calls)
	}
}

func TestClosedQueuesDrainThenReportExhaustion(t *testing.T) {
	b := NewMessageBus()
	in := InboundMessage{Content: "last inbound"}
	out := OutboundMessage{Content: "last outbound"}
	b.PublishInbound(in)
	b.PublishOutbound(out)
	b.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if got, ok := b.ConsumeInbound(ctx); !ok || !reflect.DeepEqual(got, in) {
		t.Fatalf("buffered inbound lost: %#v, %v", got, ok)
	}
	if got, ok := b.SubscribeOutbound(ctx); !ok || got != out {
		t.Fatalf("buffered outbound lost: %#v, %v", got, ok)
	}
	if got, ok := b.ConsumeInbound(ctx); ok || !reflect.DeepEqual(got, InboundMessage{}) {
		t.Errorf("closed inbound returned %#v, %v", got, ok)
	}
	if got, ok := b.SubscribeOutbound(ctx); ok || got != (OutboundMessage{}) {
		t.Errorf("closed outbound returned %#v, %v", got, ok)
	}
}
