package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)


// ChatBot handles user interactions and uses LLMClient to fetch responses.
type ChatBot struct {
	llmClient *LLMClient
	systemMsg string // System-level instruction for guiding the LLM.
	tools     []openai.FunctionDefinition // Tools available for chatbot responses.
	conversationHistory []openai.ChatCompletionMessage
}

// NewChatBot initializes a ChatBot with an LLM client, system message, and tools.
func NewChatBot(llmClient *LLMClient, systemMsg string) *ChatBot {
	if systemMsg == "" {
		systemMsg = "You are an expert assistant focused on providing the best learning resources and web links. " +
			"Whenever a user asks for recommendations or resources, include multiple reputable sources with a brief explanation of their value. " +
			"Avoid overly generic responses."
	}
	return &ChatBot{
		llmClient: llmClient,
		systemMsg: systemMsg,
		tools:     []openai.FunctionDefinition{WebSearchTool()}, // Add WebSearchTool.
		conversationHistory: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemMsg},
		},
	}
}

// InitializeChatBot sets up the chatbot instance with the OpenAI API key.
func InitializeChatBot() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OpenAI API key is missing. Please set the OPENAI_API_KEY environment variable.")
	}

	llmClient := NewLLMClient(apiKey)
	chatbot = NewChatBot(llmClient, "You are a helpful assistant for general queries.")
}

// ChatbotHandler handles chatbot queries via HTTP requests.
func ChatbotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body for the user's question.
	var req struct {
		Question string `json:"question"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Question == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get the chatbot's response to the question.
	response, err := chatbot.AnswerQuestion(req.Question)
	if err != nil {
		http.Error(w, "Error processing question: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the chatbot's answer.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

// AnswerQuestion processes a user question and fetches a response from the LLM.
func (bot *ChatBot) AnswerQuestion(question string) (string, error) {
	// Add the user message to the conversation history
	bot.conversationHistory = append(bot.conversationHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: question,
	})

	// Generate LLM response
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: bot.conversationHistory,
	}
	resp, err := bot.llmClient.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	var llmResponse string
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		llmResponse = resp.Choices[0].Message.Content
	} else {
		llmResponse = "I'm sorry, I couldn't find an answer to your question."
	}

	// Perform web search for additional resources
	webSearchResults, err := bot.performWebSearch(question)
	if err != nil {
		webSearchResults = "No additional links found due to an error."
	}

	// Combine the LLM response with web search results
	finalResponse := fmt.Sprintf("%s\n\n%s", llmResponse, webSearchResults)

	// Add the assistant's response to the conversation history
	bot.conversationHistory = append(bot.conversationHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: finalResponse,
	})

	return finalResponse, nil
}


func (bot *ChatBot) ResetConversation() {
	bot.conversationHistory = []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: bot.systemMsg},
	}
}

func (bot *ChatBot) performWebSearch(query string) (string, error) {
	// Prepare the query arguments
	args := map[string]string{"query": query}
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("failed to encode web search arguments: %w", err)
	}

	// Build the request to OpenAI
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: bot.conversationHistory,
		Functions: []openai.FunctionDefinition{
			WebSearchTool(),
		},
		FunctionCall: &openai.FunctionCall{
			Name:      "web_search",
			Arguments: string(argsJSON),
		},
	}

	// Send the request and handle the response
	resp, err := bot.llmClient.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		log.Printf("Web search failed: %v", err)
		
	}

	// Ensure response contains valid data
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		formattedResults := bot.formatSearchResults(resp.Choices[0].Message.Content)
		if formattedResults != "" {
			return formattedResults, nil
		}
	}
	return "",err
}

// generateLLMResponse fetches a default response for non-tool questions.
func (bot *ChatBot) generateLLMResponse(question string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: bot.systemMsg},
			{Role: openai.ChatMessageRoleUser, Content: question},
		},
	}

	resp, err := bot.llmClient.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	// Ensure the response contains choices and extract the content.
	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		if choice.Message.Content != "" {
			return choice.Message.Content, nil
		}
	}

	return "I couldn't find an answer to your question.", nil
}

func (bot *ChatBot) formatSearchResults(rawContent string) string {
	lines := strings.Split(rawContent, "\n")
	var formatted strings.Builder

	formatted.WriteString("Below are some useful examples and resources:\n\n")

	for _, line := range lines {
		if strings.Contains(line, "http") {
			// Assume "Description - Link" format
			parts := strings.SplitN(line, " - ", 2)
			if len(parts) == 2 {
				description := strings.TrimSpace(parts[0])
				link := strings.TrimSpace(parts[1])
				formatted.WriteString(fmt.Sprintf("- %s: %s\n", description, link))
			} else {
				// Handle standalone links
				link := strings.TrimSpace(line)
				formatted.WriteString(fmt.Sprintf("- %s\n", link))
			}
		}
	}

	return formatted.String()
}
