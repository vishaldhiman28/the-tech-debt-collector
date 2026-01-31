package rag

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// EmbeddingService generates embeddings
type EmbeddingService struct {
	client *openai.Client
	model  string
}

// NewEmbeddingService creates service
func NewEmbeddingService(apiKey string) *EmbeddingService {
	return &EmbeddingService{
		client: openai.NewClient(apiKey),
		model:  openai.SmallEmbedding3,
	}
}

// Embed generates single embedding
func (es *EmbeddingService) Embed(ctx context.Context, text string) ([]float32, error) {
	resp, err := es.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: es.model,
	})

	if err != nil {
		return nil, fmt.Errorf("embedding error: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}

// EmbedBatch generates multiple embeddings
func (es *EmbeddingService) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	resp, err := es.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: texts,
		Model: es.model,
	})

	if err != nil {
		return nil, fmt.Errorf("batch embedding error: %w", err)
	}

	embeddings := make([][]float32, len(resp.Data))
	for _, data := range resp.Data {
		embeddings[data.Index] = data.Embedding
	}

	return embeddings, nil
}
