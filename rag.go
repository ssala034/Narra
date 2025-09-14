package main

// Package implements a Retrieval-Augmented Generation (RAG) pipeline for document analysis and querying.

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
	"golang.org/x/time/rate"
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

// VectorDB represents an in-memory vector database for storing document embeddings.
type VectorDB struct {
	Embeddings []Embedding `json:"embeddings"`
}

// RAGPipeline orchestrates the entire RAG process including document indexing,
// vector embeddings, semantic search, and answer generation.
type RAGPipeline struct {
	geminiAPIKey string        `json:"gemini_api_key"`
	client       *genai.Client `json:"-"`
	vectorDB     *VectorDB     `json:"vector_db"`
	dbPath       string        `json:"db_path"`
	rateLimiter  *rate.Limiter `json:"-"`
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Document   Document `json:"document"`
	Similarity float32  `json:"similarity"`
}

// globalRagPipeline provides dependency injection for the RAG pipeline instance.
var globalRagPipeline *RAGPipeline

// StartEmbeddings initializes the RAG pipeline and processes documents from the specified root path.
// This method should be called first to build the file system tree and create document embeddings.
func (r *RAGPipeline) StartEmbeddings(rootPath string) []Document {

	geminiAPIKey := "Your-API-Key-Here"
	if geminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set")
	}

	// Initialize RAG pipeline
	pipeline, err := NewRAGPipeline(geminiAPIKey, "vector_db.json")
	if err != nil {
		log.Fatal("Failed to create RAG pipeline:", err)
	}

	globalRagPipeline = pipeline

	// Build file system tree
	root, err := BuildFileSystemTree(rootPath)
	if err != nil {
		log.Fatal("Failed to build file system tree:", err)
	}

	fmt.Printf("Built file system tree for: %s\n", root.Name)

	// Extract documents and create embeddings for indexing
	documents := getGlobalPipeline().ExtractDocuments(root)
	fmt.Printf("Extracted %d document chunks\n", len(documents))
	return documents

}

func getGlobalPipeline() *RAGPipeline {
	if globalRagPipeline == nil {
		panic("RAGPipeline not initialized. Call StartEmbeddings() first.")
	}
	return globalRagPipeline
}

// NewRAGPipeline creates a new RAG pipeline instance with the specified API key and database path.
func NewRAGPipeline(geminiAPIKey, dbPath string) (*RAGPipeline, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Create rate limiter to prevent API quota exhaustion: 1 request per second with burst capacity of 2
	rateLimiter := rate.NewLimiter(rate.Every(time.Second), 2)

	pipeline := &RAGPipeline{
		geminiAPIKey: geminiAPIKey,
		client:       client,
		vectorDB:     &VectorDB{Embeddings: []Embedding{}},
		dbPath:       dbPath,
		rateLimiter:  rateLimiter,
	}

	return pipeline, nil
}

func (r *RAGPipeline) SimpleQuery(query string) string {

	ctx := context.Background()

	fmt.Printf("\n=== Query: %s ===\n", query)

	answer, results, err := getGlobalPipeline().Query(ctx, query, 3)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return ""
	}

	var sb strings.Builder

	// Format the response with answer and supporting documents
	sb.WriteString(fmt.Sprintf("Response: %s\n\n", answer))

	// Include metadata about relevant source documents
	sb.WriteString("Relevant documents:\n")
	for i, result := range results {
		sb.WriteString(fmt.Sprintf("%d. [%.3f] %s (chunk %d)\n",
			i+1, result.Similarity, result.Document.FilePath, result.Document.ChunkIdx))
	}

	return sb.String()
}

// isTextFile determines if a file should be processed based on its extension.
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

	// If it's a file with content, process it into document chunks
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

	// Process child nodes recursively
	for _, child := range node.Children {
		r.extractDocumentsRecursive(child, documents)
	}
}

// splitText divides text into manageable chunks using paragraph boundaries and character limits.
func (r *RAGPipeline) splitText(text string) []string {
	// Split on double newlines (paragraphs) or every 1000 characters for long sections
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(text, -1)

	var chunks []string
	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		// Split overly long paragraphs into smaller chunks
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

// generateDocumentID creates a unique identifier for a document chunk using MD5 hashing.
func (r *RAGPipeline) generateDocumentID(filePath string, chunkIdx int) string {
	data := fmt.Sprintf("%s_%d", filePath, chunkIdx)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)[:16]
}

// CreateEmbedding generates vector embeddings for text using the Gemini API.
func (r *RAGPipeline) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Apply rate limiting to prevent API quota exhaustion
	err := r.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	model := r.client.EmbeddingModel("embedding-001")

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

// IndexDocuments processes documents and creates vector embeddings for the database.
func (r *RAGPipeline) IndexDocuments(documents []Document) error {
	ctx := context.Background()

	log.Printf("Indexing %d documents...", len(documents))

	for i, doc := range documents {
		if i%10 == 0 {
			log.Printf("Processing document %d/%d", i+1, len(documents))
		}

		// Skip documents that are already indexed
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

		// Generate vector embedding for document content
		vector, err := r.CreateEmbedding(ctx, doc.Content)
		if err != nil {
			log.Printf("Error creating embedding for document %s: %v", doc.ID, err)
			continue
		}

		// Store the embedding in the vector database
		embedding := Embedding{
			ID:       doc.ID,
			Vector:   vector,
			Document: doc,
			Created:  time.Now(),
		}

		r.vectorDB.Embeddings = append(r.vectorDB.Embeddings, embedding)

		// Add delay to prevent rate limiting issues
		time.Sleep(100 * time.Millisecond)
	}

	// Persist the updated database to disk
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

// SearchDocuments performs semantic search using vector similarity and returns the top N results.
func (r *RAGPipeline) SearchDocuments(ctx context.Context, query string, topN int) ([]SearchResult, error) {
	// Generate query embedding for similarity comparison
	queryVector, err := r.CreateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create query embedding: %w", err)
	}

	// Calculate cosine similarity between query and all document embeddings
	var results []SearchResult

	for _, embedding := range r.vectorDB.Embeddings {
		similarity := cosineSimilarity(queryVector, embedding.Vector)

		results = append(results, SearchResult{
			Document:   embedding.Document,
			Similarity: similarity,
		})
	}

	// Sort results by similarity score in descending order
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Similarity < results[j].Similarity {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Return the top N most relevant results
	if topN > len(results) {
		topN = len(results)
	}

	return results[:topN], nil
}

// GenerateAnswer creates responses using retrieved context and the Gemini language model.
func (r *RAGPipeline) GenerateAnswer(ctx context.Context, query string, context []string) (string, error) {
	// Apply rate limiting to prevent API quota exhaustion
	err := r.rateLimiter.Wait(ctx)
	if err != nil {
		return "", fmt.Errorf("rate limiter error: %w", err)
	}

	// TODO: Implement context-aware prompt generation
	// prompt := r.createRAGPrompt(query, context)

	// Generate response using Gemini language model
	model := r.client.GenerativeModel("gemini-1.5-flash")

	response, err := model.GenerateContent(ctx, genai.Text(query))
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return fmt.Sprintf("%v", response.Candidates[0].Content.Parts[0]), nil
}

// createRAGPrompt constructs a context-aware prompt for answer generation.
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

// Query executes the complete RAG pipeline: retrieval of relevant documents followed by answer generation.
func (r *RAGPipeline) Query(ctx context.Context, query string, topN int) (string, []SearchResult, error) {
	// Retrieve the most relevant documents for the query
	results, err := r.SearchDocuments(ctx, query, topN)
	if err != nil {
		return "", nil, fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		return "No relevant documents found for your query.", nil, nil
	}

	// Extract textual content from search results for context
	var context []string
	for _, result := range results {
		context = append(context, result.Document.Content)
	}

	// Generate contextual answer using retrieved documents
	answer, err := r.GenerateAnswer(ctx, query, context)
	if err != nil {
		return "", results, fmt.Errorf("answer generation failed: %w", err)
	}

	return answer, results, nil
}

// saveDB persists the vector database to disk in JSON format.
func (r *RAGPipeline) saveDB() error {
	data, err := json.MarshalIndent(r.vectorDB, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	return os.WriteFile(r.dbPath, data, 0644)
}

// loadDB restores the vector database from disk storage.
func (r *RAGPipeline) loadDB() error {
	data, err := os.ReadFile(r.dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Initialize with empty database if file doesn't exist
			return nil
		}
		return fmt.Errorf("failed to read database: %w", err)
	}

	return json.Unmarshal(data, r.vectorDB)
}

// GetStats returns comprehensive statistics about the indexed document collection.
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
