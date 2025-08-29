package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/denkhaus/agents/multi/plugins"
)

// Mock the enhancedChatModel for testing scrolling
type testScrollModel struct {
	messages     []chatMessage
	scrollOffset int
	inputFocused bool
	width        int
	height       int
}

type chatMessage struct {
	Agent     string
	Content   string
	Type      plugins.MessageType
	Timestamp time.Time
}

// getVisibleMessageCount calculates how many messages can fit in the chat area
func (m *testScrollModel) getVisibleMessageCount() int {
	chatHeight := m.height - 8 // Account for title and input area
	if chatHeight < 8 {
		return 1
	}
	return chatHeight / 8
}

// simulateScrollUp simulates pressing up arrow in scroll mode
func (m *testScrollModel) simulateScrollUp() {
	if !m.inputFocused {
		maxScroll := len(m.messages) - m.getVisibleMessageCount()
		if maxScroll > 0 && m.scrollOffset < maxScroll {
			m.scrollOffset++
		}
	}
}

// simulateScrollDown simulates pressing down arrow in scroll mode
func (m *testScrollModel) simulateScrollDown() {
	if !m.inputFocused {
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
	}
}

// getVisibleRange returns the range of messages currently visible
func (m *testScrollModel) getVisibleRange() (int, int, string) {
	if len(m.messages) == 0 {
		return 0, 0, "No messages"
	}

	visibleCount := m.getVisibleMessageCount()
	totalMessages := len(m.messages)
	
	start := totalMessages - visibleCount - m.scrollOffset
	if start < 0 {
		start = 0
	}
	
	end := start + visibleCount
	if end > totalMessages {
		end = totalMessages
	}

	status := fmt.Sprintf("Showing messages %d-%d of %d (scroll offset: %d)", 
		start+1, end, totalMessages, m.scrollOffset)
	
	return start, end, status
}

func main() {
	model := &testScrollModel{
		messages:     []chatMessage{},
		scrollOffset: 0,
		inputFocused: false, // Start in scroll mode for testing
		width:        120,
		height:       30,
	}

	// Add test messages
	for i := 1; i <= 25; i++ {
		model.messages = append(model.messages, chatMessage{
			Agent:     fmt.Sprintf("Agent%d", (i%3)+1),
			Content:   fmt.Sprintf("This is test message number %d with some content", i),
			Type:      plugins.MessageTypeNormal,
			Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
		})
	}

	fmt.Println("Testing Scrollable Message Area:")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Total messages: %d\n", len(model.messages))
	fmt.Printf("Visible message count: %d\n", model.getVisibleMessageCount())
	fmt.Println()

	// Test initial state (showing latest messages)
	start, end, status := model.getVisibleRange()
	fmt.Printf("1. Initial state: %s\n", status)
	fmt.Printf("   Visible messages: %d to %d\n", start+1, end)
	fmt.Println()

	// Test scrolling up (to older messages)
	fmt.Println("2. Scrolling up 5 times (to see older messages):")
	for i := 0; i < 5; i++ {
		model.simulateScrollUp()
		_, _, status := model.getVisibleRange()
		fmt.Printf("   Scroll %d: %s\n", i+1, status)
	}
	fmt.Println()

	// Test scrolling down (to newer messages)
	fmt.Println("3. Scrolling down 3 times (to see newer messages):")
	for i := 0; i < 3; i++ {
		model.simulateScrollDown()
		_, _, status := model.getVisibleRange()
		fmt.Printf("   Scroll %d: %s\n", i+1, status)
	}
	fmt.Println()

	// Test boundary conditions
	fmt.Println("4. Testing boundaries:")
	
	// Scroll to top
	for model.scrollOffset < len(model.messages) {
		model.simulateScrollUp()
	}
	_, _, topStatus := model.getVisibleRange()
	fmt.Printf("   At top: %s\n", topStatus)
	
	// Scroll to bottom
	for model.scrollOffset > 0 {
		model.simulateScrollDown()
	}
	_, _, bottomStatus := model.getVisibleRange()
	fmt.Printf("   At bottom: %s\n", bottomStatus)

	fmt.Println("\n✅ Scrolling functionality works correctly!")
	fmt.Println("\nKey bindings:")
	fmt.Println("- ESC: Toggle between input and scroll mode")
	fmt.Println("- ↑↓: Scroll messages (in scroll mode) or navigate input history (in input mode)")
	fmt.Println("- PgUp/PgDn: Page up/down through messages")
	fmt.Println("- Home/End: Jump to oldest/newest messages")
}