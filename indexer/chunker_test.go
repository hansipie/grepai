package indexer

import (
	"strings"
	"testing"
)

func TestChunker_Chunk(t *testing.T) {
	chunker := NewChunker(100, 10) // Small chunks for testing

	content := strings.Repeat("line of code\n", 50)
	chunks := chunker.Chunk("test.go", content)

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	// Verify chunk properties
	for i, chunk := range chunks {
		if chunk.ID == "" {
			t.Errorf("chunk %d has empty ID", i)
		}
		if chunk.FilePath != "test.go" {
			t.Errorf("chunk %d has wrong file path: %s", i, chunk.FilePath)
		}
		if chunk.StartLine < 1 {
			t.Errorf("chunk %d has invalid start line: %d", i, chunk.StartLine)
		}
		if chunk.EndLine < chunk.StartLine {
			t.Errorf("chunk %d has end line before start line", i)
		}
		if chunk.Content == "" {
			t.Errorf("chunk %d has empty content", i)
		}
	}
}

func TestChunker_ChunkEmptyContent(t *testing.T) {
	chunker := NewChunker(512, 50)
	chunks := chunker.Chunk("empty.go", "")

	if len(chunks) != 0 {
		t.Errorf("expected 0 chunks for empty content, got %d", len(chunks))
	}
}

func TestChunker_ChunkWhitespaceOnly(t *testing.T) {
	chunker := NewChunker(512, 50)
	chunks := chunker.Chunk("whitespace.go", "   \n\n\t\t\n   ")

	if len(chunks) != 0 {
		t.Errorf("expected 0 chunks for whitespace-only content, got %d", len(chunks))
	}
}

func TestChunker_ChunkWithContext(t *testing.T) {
	chunker := NewChunker(512, 50)
	chunks := chunker.ChunkWithContext("myfile.go", "package main\n\nfunc main() {}")

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	// Verify context is added
	if !strings.Contains(chunks[0].Content, "File: myfile.go") {
		t.Error("chunk should contain file path context")
	}
}

func TestChunker_DefaultValues(t *testing.T) {
	// Test with invalid values
	chunker := NewChunker(0, -1)

	if chunker.chunkSize != DefaultChunkSize {
		t.Errorf("expected default chunk size %d, got %d", DefaultChunkSize, chunker.chunkSize)
	}

	if chunker.overlap != DefaultChunkOverlap {
		t.Errorf("expected default overlap %d, got %d", DefaultChunkOverlap, chunker.overlap)
	}
}

func TestChunker_OverlapTooLarge(t *testing.T) {
	// Overlap >= chunk size should be reduced
	chunker := NewChunker(100, 150)

	if chunker.overlap >= chunker.chunkSize {
		t.Error("overlap should be less than chunk size")
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		text      string
		minTokens int
		maxTokens int
	}{
		{"", 0, 0},
		{"hello", 1, 2},
		{"hello world", 2, 4},
		{"func main() { fmt.Println(\"hello\") }", 4, 10},
	}

	for _, tt := range tests {
		result := EstimateTokens(tt.text)
		if result < tt.minTokens || result > tt.maxTokens {
			t.Errorf("EstimateTokens(%q) = %d, expected between %d and %d",
				tt.text, result, tt.minTokens, tt.maxTokens)
		}
	}
}
