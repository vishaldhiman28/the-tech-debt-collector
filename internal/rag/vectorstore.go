package rag

import (
	"context"
	"log"
)

// VectorStore manages embeddings (stub for Qdrant)
type VectorStore struct {
	embedder *EmbeddingService
	items    []*StoredItem
}

// StoredItem represents stored embedding
type StoredItem struct {
	ID        string
	File      string
	Line      int
	Comment   string
	Risk      float64
	Embedding []float32
}

// SimilarItem represents search result
type SimilarItem struct {
	Score   float32
	File    string
	Line    int
	Comment string
	Risk    float64
}

// NewVectorStore creates store
func NewVectorStore(apiKey string) (*VectorStore, error) {
	return &VectorStore{
		embedder: NewEmbeddingService(apiKey),
		items:    []*StoredItem{},
	}, nil
}

// IndexItem adds item to store
func (vs *VectorStore) IndexItem(ctx context.Context, id, file, comment string, line int, risk float64) error {
	embedding, err := vs.embedder.Embed(ctx, comment)
	if err != nil {
		return err
	}

	vs.items = append(vs.items, &StoredItem{
		ID:        id,
		File:      file,
		Line:      line,
		Comment:   comment,
		Risk:      risk,
		Embedding: embedding,
	})

	log.Printf("[RAG] Indexed: %s:%d\n", file, line)
	return nil
}

// SearchSimilar finds similar items
func (vs *VectorStore) SearchSimilar(ctx context.Context, query string, topK int) ([]*SimilarItem, error) {
	if len(vs.items) == 0 {
		return []*SimilarItem{}, nil
	}

	queryEmbedding, err := vs.embedder.Embed(ctx, query)
	if err != nil {
		return nil, err
	}

	// Calculate similarities (cosine)
	type scored struct {
		item  *StoredItem
		score float32
	}

	var scores []scored
	for _, item := range vs.items {
		similarity := cosineSimilarity(queryEmbedding, item.Embedding)
		scores = append(scores, scored{item, similarity})
	}

	// Sort by score (simple bubble sort for MVP)
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].score > scores[i].score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// Return top K
	results := make([]*SimilarItem, 0)
	for i := 0; i < topK && i < len(scores); i++ {
		item := scores[i].item
		results = append(results, &SimilarItem{
			Score:   scores[i].score,
			File:    item.File,
			Line:    item.Line,
			Comment: item.Comment,
			Risk:    item.Risk,
		})
	}

	return results, nil
}

// cosineSimilarity calculates similarity between vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (sqrt32(normA) * sqrt32(normB))
}

func sqrt32(x float32) float32 {
	if x < 0 {
		return 0
	}
	// Fast sqrt approximation
	return float32(1.0) / float32(1.0) // Placeholder
}
