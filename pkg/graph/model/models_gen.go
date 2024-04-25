// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AddModelDeploymentInput struct {
	ID     string `json:"id"`
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type AddModelInput struct {
	Type        ModelType         `json:"type"`
	HuggingFace *HuggingFaceInput `json:"huggingFace,omitempty"`
}

type AddPipelineDeploymentInput struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

type AddPipelineInput struct {
	Type                 PipelineType                  `json:"type"`
	Name                 string                        `json:"name"`
	RepositoryEmbeddings *AddRepositoryEmbeddingsInput `json:"repositoryEmbeddings,omitempty"`
}

type AddRepositoryEmbeddingsInput struct {
	RepositoryID string `json:"repositoryID"`
	ModelID      string `json:"modelID"`
	StorageID    string `json:"storageID"`
}

type AddRepositoryInput struct {
	URL   *string         `json:"url,omitempty"`
	Token *string         `json:"token,omitempty"`
	Type  *RepositoryType `json:"type,omitempty"`
	Owner *string         `json:"owner,omitempty"`
	Name  *string         `json:"name,omitempty"`
}

type AddStorageDeploymentInput struct {
	ID     string `json:"id"`
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type AddStorageInput struct {
	Type StorageType `json:"type"`
	Name string      `json:"name"`
}

type HuggingFace struct {
	Organization      string `json:"organization"`
	Name              string `json:"name"`
	MaxSequenceLength int    `json:"maxSequenceLength"`
}

type HuggingFaceInput struct {
	Organization      string `json:"organization"`
	Name              string `json:"name"`
	MaxSequenceLength int    `json:"maxSequenceLength"`
}

type Model struct {
	ID          string           `json:"id"`
	Type        ModelType        `json:"type"`
	DisplayName string           `json:"displayName"`
	Status      ModelStatus      `json:"status"`
	HuggingFace *HuggingFace     `json:"huggingFace,omitempty"`
	Deployment  *ModelDeployment `json:"deployment,omitempty"`
}

type ModelDeployment struct {
	Enabled bool   `json:"enabled"`
	CPU     string `json:"cpu"`
	Memory  string `json:"memory"`
}

type Mutation struct {
}

type Pipeline struct {
	ID                   string                `json:"id"`
	Name                 string                `json:"name"`
	Type                 PipelineType          `json:"type"`
	Enabled              bool                  `json:"enabled"`
	Status               PipelineStatus        `json:"status"`
	RepositoryEmbeddings *RepositoryEmbeddings `json:"repositoryEmbeddings,omitempty"`
}

type PipelineExecution struct {
	ID     string                  `json:"id"`
	Status PipelineExecutionStatus `json:"status"`
}

type Query struct {
}

type QueryInput struct {
	Query string `json:"query"`
	Page  *int   `json:"page,omitempty"`
	Limit *int   `json:"limit,omitempty"`
}

type Repository struct {
	ID          string         `json:"id"`
	Type        RepositoryType `json:"type"`
	DisplayName string         `json:"displayName"`
	Owner       string         `json:"owner"`
	Name        string         `json:"name"`
	URL         string         `json:"url"`
}

type RepositoryEmbeddings struct {
	RepositoryID string `json:"repositoryID"`
	ModelID      string `json:"modelID"`
	StorageID    string `json:"storageID"`
}

type SearchResult struct {
	ID         string  `json:"id"`
	ChunkID    int     `json:"chunkID"`
	Content    string  `json:"content"`
	Hash       string  `json:"hash"`
	Path       string  `json:"path"`
	Owner      string  `json:"owner"`
	Repo       string  `json:"repo"`
	StartIndex int     `json:"startIndex"`
	EndIndex   int     `json:"endIndex"`
	StartLine  int     `json:"startLine"`
	Score      float64 `json:"score"`
}

type Storage struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Type       StorageType        `json:"type"`
	Status     StorageStatus      `json:"status"`
	Deployment *StorageDeployment `json:"deployment,omitempty"`
}

type StorageDeployment struct {
	Enabled bool   `json:"enabled"`
	CPU     string `json:"cpu"`
	Memory  string `json:"memory"`
}

type ModelStatus string

const (
	ModelStatusNotDeployed ModelStatus = "NOT_DEPLOYED"
	ModelStatusDeploying   ModelStatus = "DEPLOYING"
	ModelStatusReady       ModelStatus = "READY"
	ModelStatusError       ModelStatus = "ERROR"
)

var AllModelStatus = []ModelStatus{
	ModelStatusNotDeployed,
	ModelStatusDeploying,
	ModelStatusReady,
	ModelStatusError,
}

func (e ModelStatus) IsValid() bool {
	switch e {
	case ModelStatusNotDeployed, ModelStatusDeploying, ModelStatusReady, ModelStatusError:
		return true
	}
	return false
}

func (e ModelStatus) String() string {
	return string(e)
}

func (e *ModelStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ModelStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ModelStatus", str)
	}
	return nil
}

func (e ModelStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ModelType string

const (
	ModelTypeHuggingface ModelType = "HUGGINGFACE"
	ModelTypeOpenai      ModelType = "OPENAI"
	ModelTypeExternal    ModelType = "EXTERNAL"
)

var AllModelType = []ModelType{
	ModelTypeHuggingface,
	ModelTypeOpenai,
	ModelTypeExternal,
}

func (e ModelType) IsValid() bool {
	switch e {
	case ModelTypeHuggingface, ModelTypeOpenai, ModelTypeExternal:
		return true
	}
	return false
}

func (e ModelType) String() string {
	return string(e)
}

func (e *ModelType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ModelType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ModelType", str)
	}
	return nil
}

func (e ModelType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PipelineExecutionStatus string

const (
	PipelineExecutionStatusActive    PipelineExecutionStatus = "ACTIVE"
	PipelineExecutionStatusSucceeded PipelineExecutionStatus = "SUCCEEDED"
	PipelineExecutionStatusFailed    PipelineExecutionStatus = "FAILED"
	PipelineExecutionStatusPending   PipelineExecutionStatus = "PENDING"
)

var AllPipelineExecutionStatus = []PipelineExecutionStatus{
	PipelineExecutionStatusActive,
	PipelineExecutionStatusSucceeded,
	PipelineExecutionStatusFailed,
	PipelineExecutionStatusPending,
}

func (e PipelineExecutionStatus) IsValid() bool {
	switch e {
	case PipelineExecutionStatusActive, PipelineExecutionStatusSucceeded, PipelineExecutionStatusFailed, PipelineExecutionStatusPending:
		return true
	}
	return false
}

func (e PipelineExecutionStatus) String() string {
	return string(e)
}

func (e *PipelineExecutionStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PipelineExecutionStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PipelineExecutionStatus", str)
	}
	return nil
}

func (e PipelineExecutionStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PipelineStatus string

const (
	PipelineStatusNotDeployed PipelineStatus = "NOT_DEPLOYED"
	PipelineStatusDeploying   PipelineStatus = "DEPLOYING"
	PipelineStatusReady       PipelineStatus = "READY"
	PipelineStatusError       PipelineStatus = "ERROR"
)

var AllPipelineStatus = []PipelineStatus{
	PipelineStatusNotDeployed,
	PipelineStatusDeploying,
	PipelineStatusReady,
	PipelineStatusError,
}

func (e PipelineStatus) IsValid() bool {
	switch e {
	case PipelineStatusNotDeployed, PipelineStatusDeploying, PipelineStatusReady, PipelineStatusError:
		return true
	}
	return false
}

func (e PipelineStatus) String() string {
	return string(e)
}

func (e *PipelineStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PipelineStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PipelineStatus", str)
	}
	return nil
}

func (e PipelineStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PipelineType string

const (
	PipelineTypeRepositoryEmbeddings PipelineType = "REPOSITORY_EMBEDDINGS"
)

var AllPipelineType = []PipelineType{
	PipelineTypeRepositoryEmbeddings,
}

func (e PipelineType) IsValid() bool {
	switch e {
	case PipelineTypeRepositoryEmbeddings:
		return true
	}
	return false
}

func (e PipelineType) String() string {
	return string(e)
}

func (e *PipelineType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PipelineType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PipelineType", str)
	}
	return nil
}

func (e PipelineType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RepositoryType string

const (
	RepositoryTypeGithub    RepositoryType = "GITHUB"
	RepositoryTypeGitlab    RepositoryType = "GITLAB"
	RepositoryTypeBitbucket RepositoryType = "BITBUCKET"
)

var AllRepositoryType = []RepositoryType{
	RepositoryTypeGithub,
	RepositoryTypeGitlab,
	RepositoryTypeBitbucket,
}

func (e RepositoryType) IsValid() bool {
	switch e {
	case RepositoryTypeGithub, RepositoryTypeGitlab, RepositoryTypeBitbucket:
		return true
	}
	return false
}

func (e RepositoryType) String() string {
	return string(e)
}

func (e *RepositoryType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RepositoryType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RepositoryType", str)
	}
	return nil
}

func (e RepositoryType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type StorageStatus string

const (
	StorageStatusNotDeployed StorageStatus = "NOT_DEPLOYED"
	StorageStatusDeploying   StorageStatus = "DEPLOYING"
	StorageStatusReady       StorageStatus = "READY"
	StorageStatusError       StorageStatus = "ERROR"
)

var AllStorageStatus = []StorageStatus{
	StorageStatusNotDeployed,
	StorageStatusDeploying,
	StorageStatusReady,
	StorageStatusError,
}

func (e StorageStatus) IsValid() bool {
	switch e {
	case StorageStatusNotDeployed, StorageStatusDeploying, StorageStatusReady, StorageStatusError:
		return true
	}
	return false
}

func (e StorageStatus) String() string {
	return string(e)
}

func (e *StorageStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StorageStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StorageStatus", str)
	}
	return nil
}

func (e StorageStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type StorageType string

const (
	StorageTypeRedis         StorageType = "REDIS"
	StorageTypePostgres      StorageType = "POSTGRES"
	StorageTypeElasticsearch StorageType = "ELASTICSEARCH"
)

var AllStorageType = []StorageType{
	StorageTypeRedis,
	StorageTypePostgres,
	StorageTypeElasticsearch,
}

func (e StorageType) IsValid() bool {
	switch e {
	case StorageTypeRedis, StorageTypePostgres, StorageTypeElasticsearch:
		return true
	}
	return false
}

func (e StorageType) String() string {
	return string(e)
}

func (e *StorageType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StorageType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StorageType", str)
	}
	return nil
}

func (e StorageType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
