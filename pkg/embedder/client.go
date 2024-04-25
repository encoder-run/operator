package embedder

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/encoder-run/operator/pkg/common"
)

// Client for embedding source code files
type EmbeddingClient struct {
	httpClient *http.Client
	baseURL    string
}

type CodeEmbeddingRequest struct {
	Path    string
	Content string
	Hash    string
}

type CodeEmbeddingChunk struct {
	ChunkID    int       `json:"chunk_id"`
	FileHash   string    `json:"file_hash"`
	Code       string    `json:"code"`
	StartIndex int       `json:"start_index"`
	EndIndex   int       `json:"end_index"`
	Embedding  []float32 `json:"embedding"`
}

type CodeEmbeddings struct {
	Embeddings []CodeEmbeddingChunk `json:"embeddings"`
}

type CodeEmbeddingsResponse struct {
	Results map[string]CodeEmbeddings `json:"results"`
}

// NewClient creates a new client for fetching embeddings.
func NewClient(modelId, namespace string) *EmbeddingClient {
	return &EmbeddingClient{
		httpClient: &http.Client{},
		baseURL:    common.ModelServiceURL(modelId, namespace),
	}
}

// FetchEmbeddings sends a batch of file content to the inference API and retrieves embeddings.
func (ec *EmbeddingClient) FetchEmbeddings(requests []CodeEmbeddingRequest) (*CodeEmbeddingsResponse, error) {
	instances := make([]map[string]string, 0)
	for _, f := range requests {
		instances = append(instances, map[string]string{"file_path": f.Path, "code": f.Content, "file_hash": f.Hash})
	}
	payload, err := json.Marshal(map[string]interface{}{"instances": instances})
	if err != nil {
		return nil, err
	}

	resp, err := ec.httpClient.Post(ec.baseURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CodeEmbeddingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
