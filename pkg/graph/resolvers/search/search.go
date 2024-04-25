package search

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/common"
	"github.com/encoder-run/operator/pkg/embedder"
	"github.com/encoder-run/operator/pkg/graph/converters"
	"github.com/encoder-run/operator/pkg/graph/model"
	redigoredis "github.com/gomodule/redigo/redis"
	"github.com/redis/go-redis/v9"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Semantic(ctx context.Context, query model.QueryInput) ([]*model.SearchResult, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// List all the pipelineList.
	pipelineList := &v1alpha1.PipelineList{}
	if err := ctrlClient.List(ctx, pipelineList, &client.ListOptions{Namespace: "default"}); err != nil {
		return nil, err
	}

	results := make([]*model.SearchResult, 0, len(pipelineList.Items))

	// search the pipelines by producing a vector representation of the query
	// with the model used in that pipeline.
	for _, pipeline := range pipelineList.Items {

		if pipeline.Spec.RepositoryEmbeddings == nil {
			continue
		}

		// Get the storage.
		storageCRD := &v1alpha1.Storage{}
		if err := ctrlClient.Get(ctx, client.ObjectKey{Name: pipeline.Spec.RepositoryEmbeddings.Storage.Name, Namespace: pipeline.Namespace}, storageCRD); err != nil {
			return nil, err
		}

		switch storageCRD.Spec.Type {
		case v1alpha1.StorageTypeRedis:
			rs, err := semanticSearchRedis(ctrlClient, &pipeline, storageCRD, &query)
			if err != nil {
				return nil, err
			}
			results = append(results, rs...)
		default:
			return nil, fmt.Errorf("unsupported storage type: %s", storageCRD.Spec.Type)
		}

	}
	// Sort based on score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score // Assuming the field is named 'Score' and we are sorting ascending
	})
	return results, nil
}

func semanticSearchRedis(ctrlClient client.Client, pipeline *v1alpha1.Pipeline, storage *v1alpha1.Storage, query *model.QueryInput) ([]*model.SearchResult, error) {
	// get repository object
	repository := v1alpha1.Repository{}
	if err := ctrlClient.Get(context.Background(), types.NamespacedName{Name: pipeline.Spec.RepositoryEmbeddings.Repository.Name, Namespace: pipeline.Namespace}, &repository); err != nil {
		return nil, err
	}

	if repository.Spec.Type != v1alpha1.RepositoryTypeGithub {
		return nil, fmt.Errorf("unsupported repository type: %s", repository.Spec.Type)
	}

	redisearchClient, err := getSearchClient(ctrlClient, storage, repository.Spec.Github.URL)
	if err != nil {
		return nil, err
	}
	redisClient, err := getRedisClient(ctrlClient, storage)
	if err != nil {
		return nil, err
	}
	modelId := pipeline.Spec.RepositoryEmbeddings.Model.Name
	// search the model with the query
	c := embedder.NewClient(modelId, pipeline.Namespace)
	response, err := c.FetchEmbeddings([]embedder.CodeEmbeddingRequest{
		{
			Path:    "/",
			Content: query.Query,
			Hash:    "query",
		},
	})
	if err != nil {
		return nil, err
	}
	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no results returned from the model")
	}

	codeEmb, ok := response.Results["/"]
	if !ok {
		return nil, fmt.Errorf("no embeddings returned for the query")
	}

	if len(codeEmb.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned for the query")
	}

	emb := codeEmb.Embeddings[0].Embedding

	// Query vector represented as blob
	queryBlob := convertToBlob(emb)

	// Set up KNN search
	knnQuery := fmt.Sprintf("*=>[KNN %d @embedding $B AS __vec_score]", 25)

	redisQuery := redisearch.NewQuery(knnQuery).
		SetParams(map[string]interface{}{"B": queryBlob}).
		SetSortBy("__vec_score", true). // Sort by the vector score
		AddReturnFields("__vec_score", "chunkID", "fileHash", "filePath", "startIndex", "endIndex").
		SetDialect(2).
		Limit(0, 25)

	docs, _, err := redisearchClient.Search(redisQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	results := make([]*model.SearchResult, 0, len(docs))
	for _, doc := range docs {
		sr, err := converters.RedisEmbeddingDocToSearchResult(&doc, &repository)
		if err != nil {
			return nil, err
		}
		results = append(results, sr)
	}

	// Get the file content for the search results
	for _, sr := range results {
		// Get the file content
		key := fmt.Sprintf("%s:%s:%s:%s", repository.Spec.Github.URL, "object", "blob", sr.Hash)
		content, err := redisClient.Get(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}
		adjustedContent, startLine, err := extractContentWindowIndex(content, sr.StartIndex, sr.EndIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to extract content window from file key %s: %w", key, err)
		}
		sr.Content = adjustedContent
		sr.StartLine = startLine
	}

	return results, nil
}

func extractContentWindowIndex(content string, startIndex int, endIndex int) (string, int, error) {
	if startIndex < 0 || endIndex < 0 || startIndex > endIndex {
		return "", 0, fmt.Errorf("Invalid index range")
	}
	if endIndex > len(content) {
		endIndex = len(content)
	}
	
	// Calculate the starting line number
	startLine := 1 // Line numbering starts at 1
	for _, ch := range content[:startIndex] {
		if ch == '\n' {
			startLine++
		}
	}
	
	return content[startIndex:endIndex], startLine, nil
}
func getRedisClient(k8sClient client.Client, storage *v1alpha1.Storage) (*redis.Client, error) {
	// Retrieve the secret containing the Redis password
	secret := &corev1.Secret{}
	err := k8sClient.Get(context.Background(), types.NamespacedName{Name: storage.Name, Namespace: storage.Namespace}, secret)
	if err != nil {
		return nil, err
	}

	passwordBytes, ok := secret.Data["password"]
	if !ok {
		return nil, errors.New("password not found in secret")
	}

	password := string(passwordBytes)
	host := common.RedisServiceURL(storage.Name, storage.Namespace)

	c := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})

	return c, nil
}

func getSearchClient(k8sClient client.Client, storage *v1alpha1.Storage, ns string) (*redisearch.Client, error) {

	// Retrieve the secret containing the Redis password
	secret := &corev1.Secret{}
	err := k8sClient.Get(context.Background(), types.NamespacedName{Name: storage.Name, Namespace: storage.Namespace}, secret)
	if err != nil {
		return nil, err
	}

	passwordBytes, ok := secret.Data["password"]
	if !ok {
		return nil, errors.New("password not found in secret")
	}

	password := string(passwordBytes)
	host := common.RedisServiceURL(storage.Name, storage.Namespace)

	pool := &redigoredis.Pool{Dial: func() (redigoredis.Conn, error) {
		return redigoredis.Dial("tcp", host, redigoredis.DialPassword(password))
	}}

	c := redisearch.NewClientFromPool(pool, fmt.Sprintf("%s:%s", ns, "embedding"))

	return c, nil
}

func convertToBlob(vector []float32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, vector)
	if err != nil {
		log.Fatal("Failed to convert vector to blob: ", err)
	}
	return buf.Bytes()
}
