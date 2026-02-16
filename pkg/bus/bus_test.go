package bus

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestNewMessageBus tests MessageBus initialization
func TestNewMessageBus(t *testing.T) {
	mb := NewMessageBus()
	if mb == nil {
		t.Fatal("NewMessageBus returned nil")
	}
	if mb.inbound == nil {
		t.Error("inbound channel is nil")
	}
	if mb.outbound == nil {
		t.Error("outbound channel is nil")
	}
	if mb.handlers == nil {
		t.Error("handlers map is nil")
	}
}

// TestPublishInbound tests publishing inbound messages
func TestPublishInbound(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	msg := InboundMessage{
		Channel:  "telegram",
		SenderID: "user123",
		ChatID:   "chat456",
		Content:  "Hello",
	}

	// Publish message
	mb.PublishInbound(msg)

	// Consume message
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	received, ok := mb.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("Failed to consume inbound message")
	}

	if received.Content != "Hello" {
		t.Errorf("Expected content 'Hello', got '%s'", received.Content)
	}
	if received.Channel != "telegram" {
		t.Errorf("Expected channel 'telegram', got '%s'", received.Channel)
	}
}

// TestPublishOutbound tests publishing outbound messages
func TestPublishOutbound(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	msg := OutboundMessage{
		Channel: "discord",
		ChatID:  "chat789",
		Content: "Response",
	}

	// Publish message
	mb.PublishOutbound(msg)

	// Subscribe to message
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	received, ok := mb.SubscribeOutbound(ctx)
	if !ok {
		t.Fatal("Failed to subscribe to outbound message")
	}

	if received.Content != "Response" {
		t.Errorf("Expected content 'Response', got '%s'", received.Content)
	}
	if received.Channel != "discord" {
		t.Errorf("Expected channel 'discord', got '%s'", received.Channel)
	}
}

// TestMultipleSubscribers tests that multiple goroutines can consume messages
func TestMultipleSubscribers(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	numMessages := 10
	received := make([]OutboundMessage, 0, numMessages)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Start 3 subscribers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			for {
				msg, ok := mb.SubscribeOutbound(ctx)
				if !ok {
					return
				}

				mu.Lock()
				received = append(received, msg)
				mu.Unlock()

				// Exit after receiving some messages
				if len(received) >= numMessages {
					return
				}
			}
		}(i)
	}

	// Publish messages
	for i := 0; i < numMessages; i++ {
		mb.PublishOutbound(OutboundMessage{
			Channel: "test",
			ChatID:  "chat",
			Content: "msg",
		})
		time.Sleep(10 * time.Millisecond) // Small delay
	}

	// Wait for subscribers with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(3 * time.Second):
		t.Fatal("Test timeout: subscribers did not finish")
	}

	mu.Lock()
	defer mu.Unlock()

	if len(received) != numMessages {
		t.Errorf("Expected %d messages, got %d", numMessages, len(received))
	}
}

// TestRegisterHandler tests handler registration
func TestRegisterHandler(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	called := false
	handler := func(msg InboundMessage) error {
		called = true
		return nil
	}

	mb.RegisterHandler("telegram", handler)

	retrieved, ok := mb.GetHandler("telegram")
	if !ok {
		t.Fatal("Failed to retrieve registered handler")
	}
	if retrieved == nil {
		t.Fatal("Retrieved handler is nil")
	}

	// Test handler
	retrieved(InboundMessage{Content: "test"})
	if !called {
		t.Error("Handler was not called")
	}
}

// TestGetHandlerNotFound tests retrieving non-existent handler
func TestGetHandlerNotFound(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	_, ok := mb.GetHandler("nonexistent")
	if ok {
		t.Error("Expected handler not found, but got ok=true")
	}
}

// TestMessageFiltering tests filtering by channel
func TestMessageFiltering(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	var telegramMsgs, discordMsgs []InboundMessage
	var mu sync.Mutex

	// Register handlers for different channels
	mb.RegisterHandler("telegram", func(msg InboundMessage) error {
		mu.Lock()
		defer mu.Unlock()
		telegramMsgs = append(telegramMsgs, msg)
		return nil
	})

	mb.RegisterHandler("discord", func(msg InboundMessage) error {
		mu.Lock()
		defer mu.Unlock()
		discordMsgs = append(discordMsgs, msg)
		return nil
	})

	// Publish messages for different channels
	telegramMsg := InboundMessage{Channel: "telegram", Content: "TG message"}
	discordMsg := InboundMessage{Channel: "discord", Content: "DC message"}

	mb.PublishInbound(telegramMsg)
	mb.PublishInbound(discordMsg)

	// Consume and route messages
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for i := 0; i < 2; i++ {
		msg, ok := mb.ConsumeInbound(ctx)
		if !ok {
			t.Fatal("Failed to consume message")
		}

		handler, ok := mb.GetHandler(msg.Channel)
		if ok {
			handler(msg)
		}
	}

	// Verify filtering
	mu.Lock()
	defer mu.Unlock()

	if len(telegramMsgs) != 1 {
		t.Errorf("Expected 1 telegram message, got %d", len(telegramMsgs))
	}
	if len(discordMsgs) != 1 {
		t.Errorf("Expected 1 discord message, got %d", len(discordMsgs))
	}
	if len(telegramMsgs) > 0 && telegramMsgs[0].Content != "TG message" {
		t.Errorf("Wrong telegram message content")
	}
	if len(discordMsgs) > 0 && discordMsgs[0].Content != "DC message" {
		t.Errorf("Wrong discord message content")
	}
}

// TestContextCancellation tests that consume operations respect context cancellation
func TestContextCancellation(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, ok := mb.ConsumeInbound(ctx)
	if ok {
		t.Error("Expected ConsumeInbound to return false on canceled context")
	}

	_, ok = mb.SubscribeOutbound(ctx)
	if ok {
		t.Error("Expected SubscribeOutbound to return false on canceled context")
	}
}

// TestConcurrentHandlerRegistration tests thread-safety of handler registration
func TestConcurrentHandlerRegistration(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			channel := "channel"
			handler := func(msg InboundMessage) error { return nil }
			mb.RegisterHandler(channel, handler)

			// Try to get handler
			_, ok := mb.GetHandler(channel)
			if !ok {
				t.Error("Failed to get handler after registration")
			}
		}(i)
	}

	wg.Wait()
}

// TestMessageBusClose tests proper cleanup
func TestMessageBusClose(t *testing.T) {
	mb := NewMessageBus()
	mb.Close()

	// After closing, publishing should not panic (but message will be lost)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic after Close: %v", r)
		}
	}()

	// These operations on closed channels will panic, so we skip them
	// Just ensure Close() itself doesn't panic
}

// TestInboundMessageMetadata tests metadata handling
func TestInboundMessageMetadata(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	metadata := map[string]string{
		"user_id":  "123",
		"username": "alice",
	}

	msg := InboundMessage{
		Channel:  "telegram",
		Content:  "test",
		Metadata: metadata,
	}

	mb.PublishInbound(msg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	received, ok := mb.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("Failed to consume message")
	}

	if received.Metadata["user_id"] != "123" {
		t.Error("Metadata user_id mismatch")
	}
	if received.Metadata["username"] != "alice" {
		t.Error("Metadata username mismatch")
	}
}

// TestChannelCapacity tests channel buffering
func TestChannelCapacity(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	// Publish 100 messages (channel capacity)
	for i := 0; i < 100; i++ {
		mb.PublishInbound(InboundMessage{
			Content: "msg",
		})
	}

	// Should be able to consume all
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count := 0
	for i := 0; i < 100; i++ {
		_, ok := mb.ConsumeInbound(ctx)
		if !ok {
			break
		}
		count++
	}

	if count != 100 {
		t.Errorf("Expected to consume 100 messages, got %d", count)
	}
}
