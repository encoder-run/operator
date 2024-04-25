package converters

import (
	"fmt"
	"strconv"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/graph/model"
)

func RedisEmbeddingDocToSearchResult(doc *redisearch.Document, repo *v1alpha1.Repository) (*model.SearchResult, error) {
	sr := &model.SearchResult{}
	sr.ID = doc.Id
	if repo.Spec.Type != v1alpha1.RepositoryTypeGithub {
		return nil, fmt.Errorf("unsupported repository type: %s", repo.Spec.Type)
	}

	sr.Owner = repo.Spec.Github.Owner
	sr.Repo = repo.Spec.Github.Name

	// Get the chunk id
	chunkIDString, ok := doc.Properties["chunkID"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to convert chunk to int")
	}
	chunkID, err := strconv.Atoi(chunkIDString)
	if err != nil {
		return nil, err
	}
	sr.ChunkID = chunkID

	// Get the hash
	hash, ok := doc.Properties["fileHash"].(string)
	if !ok {
		return nil, fmt.Errorf("hash property not found in the document")
	}
	sr.Hash = hash

	// Get the file path
	path, ok := doc.Properties["filePath"].(string)
	if !ok {
		return nil, fmt.Errorf("filePath property not found in the document")
	}
	sr.Path = path

	// Get the score
	scoreString, ok := doc.Properties["__vec_score"].(string)
	if !ok {
		return nil, fmt.Errorf("vec_score property not found in the document")
	}
	score, err := strconv.ParseFloat(scoreString, 64)
	if err != nil {
		return nil, err
	}
	sr.Score = score

	// Get the start index
	startIndexString, ok := doc.Properties["startIndex"].(string)
	if !ok {
		return nil, fmt.Errorf("startIndex property not found in the document")
	}
	startIndex, err := strconv.Atoi(startIndexString)
	if err != nil {
		return nil, err
	}
	sr.StartIndex = startIndex

	// Get the end index
	endIndexString, ok := doc.Properties["endIndex"].(string)
	if !ok {
		return nil, fmt.Errorf("endIndex property not found in the document")
	}
	endIndex, err := strconv.Atoi(endIndexString)
	if err != nil {
		return nil, err
	}
	sr.EndIndex = endIndex

	return sr, nil

}
