package main

// will build here the retrieving of relavant documnent that wan be used for the LLM

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Document represents a text chunk with metadata
type Document struct {
	ID       string    `json:"id"`
	Content  string    `json:"content"`
	FilePath string    `json:"file_path"`
	ChunkIdx int       `json:"chunk_idx"`
	Created  time.Time `json:"created"`
}

// Embedding represents a vector embedding with metadata
type Embedding struct {
	ID       string    `json:"id"`
	Vector   []float32 `json:"vector"`
	Document Document  `json:"document"`
	Created  time.Time `json:"created"`
}

// VectorDB represents our simple in-memory vector database
type VectorDB struct {
	Embeddings []Embedding `json:"embeddings"`
}

// RAGPipeline orchestrates the entire RAG process
type RAGPipeline struct {
	geminiAPIKey string
	client       *genai.Client
	vectorDB     *VectorDB
	dbPath       string
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Document   Document
	Similarity float32
}

// 1. INDEXING PIPELINE

// NewRAGPipeline creates a new RAG pipeline instance
func NewRAGPipeline(geminiAPIKey, dbPath string) (*RAGPipeline, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	pipeline := &RAGPipeline{
		geminiAPIKey: geminiAPIKey,
		client:       client,
		vectorDB:     &VectorDB{Embeddings: []Embedding{}},
		dbPath:       dbPath,
	}

	// Try to load existing database
	pipeline.loadDB()

	return pipeline, nil
}

// isTextFile checks if a file is a text file based on extension
func isTextFile(filename string) bool {
	textExtensions := []string{".txt", ".md", ".go", ".py", ".js", ".java", ".html", ".css", ".json", ".yaml", ".yml", ".xml", ".csv"}
	ext := strings.ToLower(filepath.Ext(filename))

	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	return false
}

// ExtractDocuments extracts all text documents from the file system tree
func (r *RAGPipeline) ExtractDocuments(root *FileSystemNode) []Document {
	var documents []Document
	r.extractDocumentsRecursive(root, &documents)
	return documents
}

func (r *RAGPipeline) extractDocumentsRecursive(node *FileSystemNode, documents *[]Document) {
	if node == nil {
		return
	}

	// If it's a file with content, process it
	if !node.IsDir && len(node.Content) > 0 {
		chunks := r.splitText(string(node.Content))

		for i, chunk := range chunks {
			if strings.TrimSpace(chunk) == "" {
				continue
			}

			docID := r.generateDocumentID(node.Path, i)
			doc := Document{
				ID:       docID,
				Content:  strings.TrimSpace(chunk),
				FilePath: node.Path,
				ChunkIdx: i,
				Created:  time.Now(),
			}
			*documents = append(*documents, doc)
		}
	}

	// Process children
	for _, child := range node.Children {
		r.extractDocumentsRecursive(child, documents)
	}
}

// splitText splits text into chunks (simple paragraph-based splitting)
func (r *RAGPipeline) splitText(text string) []string {
	// Split on double newlines (paragraphs) or every 1000 characters for long lines
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(text, -1)

	var chunks []string
	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		// If paragraph is too long, split it further
		if len(paragraph) > 1000 {
			subChunks := r.splitLongText(paragraph, 1000)
			chunks = append(chunks, subChunks...)
		} else {
			chunks = append(chunks, paragraph)
		}
	}

	return chunks
}

func (r *RAGPipeline) splitLongText(text string, maxLength int) []string {
	var chunks []string
	words := strings.Fields(text)

	var currentChunk strings.Builder
	for _, word := range words {
		if currentChunk.Len()+len(word)+1 > maxLength {
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
			}
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString(" ")
		}
		currentChunk.WriteString(word)
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// generateDocumentID creates a unique ID for a document chunk
func (r *RAGPipeline) generateDocumentID(filePath string, chunkIdx int) string {
	data := fmt.Sprintf("%s_%d", filePath, chunkIdx)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)[:16]
}

// CreateEmbedding creates an embedding for a given text using Gemini
func (r *RAGPipeline) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	model := r.client.EmbeddingModel("embedding-001") // not sure if gemini has that, other use chromadb

	batch := model.NewBatch().AddContent(genai.Text(text))
	result, err := model.BatchEmbedContents(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return result.Embeddings[0].Values, nil
}

// IndexDocuments processes and indexes all documents
func (r *RAGPipeline) IndexDocuments(documents []Document) error {
	ctx := context.Background()

	log.Printf("Indexing %d documents...", len(documents))

	for i, doc := range documents {
		if i%10 == 0 {
			log.Printf("Processing document %d/%d", i+1, len(documents))
		}

		// Check if document is already indexed
		exists := false
		for _, existing := range r.vectorDB.Embeddings {
			if existing.Document.ID == doc.ID {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		// Create embedding
		vector, err := r.CreateEmbedding(ctx, doc.Content)
		if err != nil {
			log.Printf("Error creating embedding for document %s: %v", doc.ID, err)
			continue
		}

		// Store embedding
		embedding := Embedding{
			ID:       doc.ID,
			Vector:   vector,
			Document: doc,
			Created:  time.Now(),
		}

		r.vectorDB.Embeddings = append(r.vectorDB.Embeddings, embedding)

		// Add small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	// Save to disk
	return r.saveDB()
}

// 2. RETRIEVAL

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (normA * normB)
}

// SearchDocuments performs semantic search and returns top n results
func (r *RAGPipeline) SearchDocuments(ctx context.Context, query string, topN int) ([]SearchResult, error) {
	// Create embedding for query
	queryVector, err := r.CreateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create query embedding: %w", err)
	}

	// Calculate similarities
	var results []SearchResult

	for _, embedding := range r.vectorDB.Embeddings {
		similarity := cosineSimilarity(queryVector, embedding.Vector)

		results = append(results, SearchResult{
			Document:   embedding.Document,
			Similarity: similarity,
		})
	}

	// Sort by similarity (descending)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Similarity < results[j].Similarity {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Return top N results
	if topN > len(results) {
		topN = len(results)
	}

	return results[:topN], nil
}

// 3. GENERATION

// GenerateAnswer generates an answer using retrieved context and Gemini
func (r *RAGPipeline) GenerateAnswer(ctx context.Context, query string, context []string) (string, error) {
	// Create prompt with context
	prompt := r.createRAGPrompt(query, context)

	// might need to give apiKey here ?!?!?!?

	// Generate response using Gemini
	model := r.client.GenerativeModel("gemini-2.5-flash") // try pro later !!!

	response, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return fmt.Sprintf("%v", response.Candidates[0].Content.Parts[0]), nil
}

// createRAGPrompt creates a prompt with context for generation
func (r *RAGPipeline) createRAGPrompt(query string, context []string) string {
	contextText := strings.Join(context, "\n\n")

	prompt := fmt.Sprintf(`You are a helpful and knowledgeable assistant that answers questions based on the provided context from code repositories and documentation.

Use the context below to answer the user's question. Be comprehensive and include relevant details, but explain technical concepts in a clear, accessible way.

If the context doesn't contain enough information to fully answer the question, say so and provide what information you can. If the question is out of context say so.

CONTEXT:
%s

QUESTION: %s

ANSWER:`, contextText, query)

	return prompt
}

// 4. MAIN RAG QUERY FUNCTION

// Query performs the complete RAG pipeline: retrieve relevant documents and generate answer
func (r *RAGPipeline) Query(ctx context.Context, query string, topN int) (string, []SearchResult, error) {
	// Retrieve relevant documents
	results, err := r.SearchDocuments(ctx, query, topN)
	if err != nil {
		return "", nil, fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		return "No relevant documents found for your query.", nil, nil
	}

	// Extract context from search results
	var context []string
	for _, result := range results {
		context = append(context, result.Document.Content)
	}

	// Generate answer
	answer, err := r.GenerateAnswer(ctx, query, context)
	if err != nil {
		return "", results, fmt.Errorf("answer generation failed: %w", err)
	}

	return answer, results, nil
}

// 5. DATABASE PERSISTENCE

// saveDB saves the vector database to disk
func (r *RAGPipeline) saveDB() error {
	data, err := json.MarshalIndent(r.vectorDB, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	return os.WriteFile(r.dbPath, data, 0644)
}

// loadDB loads the vector database from disk
func (r *RAGPipeline) loadDB() error {
	data, err := os.ReadFile(r.dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Database doesn't exist yet, start with empty
			return nil
		}
		return fmt.Errorf("failed to read database: %w", err)
	}

	return json.Unmarshal(data, r.vectorDB)
}

// GetStats returns statistics about the indexed data
func (r *RAGPipeline) GetStats() map[string]interface{} {
	fileCount := make(map[string]int)

	for _, embedding := range r.vectorDB.Embeddings {
		ext := filepath.Ext(embedding.Document.FilePath)
		fileCount[ext]++
	}

	return map[string]interface{}{
		"total_documents":    len(r.vectorDB.Embeddings),
		"files_by_extension": fileCount,
		"database_size":      fmt.Sprintf("%.2f MB", float64(len(r.vectorDB.Embeddings)*len(r.vectorDB.Embeddings[0].Vector)*4)/1024/1024),
	}
}

/*
Steps

1) need to load and store data
2) which embedding models to use




prompt --> [  vector search --> llm + context] --> output
 data indexing (loading, splitting, embedding, and later retrieval)
 as part of my vectorized search so gemini can get more context before giving me a prompted answer

*/
