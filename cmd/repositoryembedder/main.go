package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
	v1alpha1 "github.com/encoder-run/operator/api/cloud/v1alpha1"
	rediscache "github.com/encoder-run/operator/pkg/cache/redis"
	"github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage"
	redigoredis "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RepositorySpec represents the specification of the repository.
type RepositorySpec struct {
	Type string
	URL  string
}

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
	ChunkID     int       `json:"chunk_id"`
	FileHash    string    `json:"file_hash"`
	Code        string    `json:"code"`
	StartLine   int       `json:"start_line"`
	EndLine     int       `json:"end_line"`
	StartColumn int       `json:"start_column"`
	EndColumn   int       `json:"end_column"`
	Embedding   []float32 `json:"embedding"`
}

type CodeEmbeddings struct {
	Embeddings []CodeEmbeddingChunk `json:"embeddings"`
}

type CodeEmbeddingsResponse struct {
	Results map[string]CodeEmbeddings `json:"results"`
}

// NewEmbeddingClient creates a new client for fetching embeddings.
func NewEmbeddingClient(modelId, namespace string) *EmbeddingClient {
	return &EmbeddingClient{
		httpClient: &http.Client{},
		baseURL:    fmt.Sprintf("http://%s-predictor-default.%s.svc.cluster.local:80/v1/models/custom-model:predict", modelId, namespace),
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

func main() {
	// Define flags
	var storageId string
	var repositoryId string
	var modelId string

	flag.StringVar(&storageId, "storageId", "", "Storage ID")
	flag.StringVar(&repositoryId, "repositoryId", "", "Repository ID")
	flag.StringVar(&modelId, "modelId", "", "Model ID")

	// Parse flags
	flag.Parse()

	// Example usage of flags in the application logic
	fmt.Printf("Using storage ID: %s\n", storageId)
	fmt.Printf("Using repository ID: %s\n", repositoryId)
	fmt.Printf("Using model ID: %s\n", modelId)

	// Check if all required arguments are provided
	if storageId == "" || repositoryId == "" || modelId == "" {
		Warning("All arguments (storageId, repositoryId, modelId) are required.")
		os.Exit(1)
	}

	// Get the kubernetes client.
	c, err := defaultClient()
	if err != nil {
		CheckIfError(err)

	}

	// Get the namespace from the service account.
	ns, err := namespace()
	if err != nil {
		CheckIfError(err)
	}

	// Get the repository by name.
	repo := &v1alpha1.Repository{}
	if err := c.Get(context.TODO(), client.ObjectKey{Name: repositoryId, Namespace: ns}, repo); err != nil {
		CheckIfError(err)
	}

	// Get the storage by name.
	st := &v1alpha1.Storage{}
	if err := c.Get(context.TODO(), client.ObjectKey{Name: storageId, Namespace: ns}, st); err != nil {
		CheckIfError(err)
	}
	// Initialize client
	client := NewEmbeddingClient(modelId, ns)

	// Repository URL for remote git repository
	var url string
	switch repo.Spec.Type {
	case v1alpha1.RepositoryTypeGithub:
		url = repo.Spec.Github.URL
	default:
		CheckIfError(fmt.Errorf("unsupported repository type: %s", repo.Spec.Type))
	}

	var storer storage.Storer
	var redisearchClient *redisearch.Client
	var redisClient *redis.Client
	switch st.Spec.Type {
	case v1alpha1.StorageTypeRedis:
		// Get the go-git storage storer based on the storage type
		s, opts, err := redisStorageStorer(c, st, url)
		if err != nil {
			CheckIfError(err)
		}
		storer = s
		redisearchClient = getRedisearchClient(opts, "embedding")
		redisClient = redis.NewClient(opts)
		err = createIndex(redisearchClient)
		if err != nil {
			CheckIfError(err)
		}
	default:
		CheckIfError(fmt.Errorf("unsupported storage type: %s", st.Spec.Type))
	}
	// auth gets the go-git auth created.
	auth, err := gitAuth(c, repo)
	if err != nil {
		CheckIfError(err)
	}

	r, err := git.Open(storer, nil)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			r, err = git.Clone(storer, nil, &git.CloneOptions{
				URL:  url,
				Auth: auth,
			})
			CheckIfError(err)
		} else {
			CheckIfError(err)
		}
	}
	// Fetch changes from the remote repository
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			"+refs/heads/*:refs/remotes/origin/*",
		},
		Auth:  auth,
		Depth: 3,
	})
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}

	// err = r.Prune(git.PruneOptions{})
	// CheckIfError(err)

	// Manually update local branch reference to match the remote tracking branch
	// Typically in a bare repo, this might be done in response to a push or a hook
	remoteRef, err := r.Reference(plumbing.NewRemoteReferenceName("origin", "main"), false)
	CheckIfError(err)

	// Update local main directly to point to the same commit
	localRef := plumbing.NewHashReference(plumbing.NewBranchReferenceName("main"), remoteRef.Hash())
	err = r.Storer.SetReference(localRef)
	CheckIfError(err)

	commit, err := r.CommitObject(remoteRef.Hash())
	if err != nil {
		log.Fatal(err)
	}

	tree, err := commit.Tree()
	if err != nil {
		log.Fatal(err)
	}

	// Check for existing processed hashes
	existingHashes := make(map[string]bool)
	existingTreeHash, err := redisClient.Get(context.Background(), "embedding:tree").Result()
	if err != nil && err != redis.Nil {
		log.Fatal(err)
	}
	if err == nil {
		existingHashesStr, err := redisClient.Get(context.Background(), fmt.Sprintf("embedding:tree:%s", existingTreeHash)).Result()
		if err != nil && err != redis.Nil {
			log.Fatal(err)
		}
		for _, h := range strings.Split(existingHashesStr, ",") {
			existingHashes[h] = true
		}
	}

	hashes := make(map[string]bool)
	filesBatch := []CodeEmbeddingRequest{}
	treeIter := tree.Files()
	batchSize := 10
	count := 0

	for {
		file, err := treeIter.Next()
		if err != nil {
			if err == io.EOF {
				break // No more files
			}
			log.Fatal(err)
		}

		hashes[file.Hash.String()] = true

		// Process only if the file hash is not in the existing hashes
		if !existingHashes[file.Hash.String()] {
			content, err := file.Contents()
			if err != nil {
				log.Fatal(err)
			}

			filesBatch = append(filesBatch, CodeEmbeddingRequest{
				Path:    file.Name,
				Content: content,
				Hash:    file.Hash.String(),
			})
			count++

			// Process in batches of 10
			if count >= batchSize {
				fmt.Printf("Processing batch of files\n")
				embeddings, err := client.FetchEmbeddings(filesBatch)
				if err != nil {
					log.Fatal(err)
				}
				err = setCodeEmbeddings(redisearchClient, *embeddings)
				if err != nil {
					log.Fatal(err)
				}
				// Reset the batch
				filesBatch = []CodeEmbeddingRequest{}
				count = 0
			}
		} else {
			fmt.Printf("Skipping file '%s' since its hash is already processed\n", file.Name)
		}
	}

	// Process any remaining files
	if len(filesBatch) > 0 {
		fmt.Printf("Processing remaining files\n")
		embeddings, err := client.FetchEmbeddings(filesBatch)
		if err != nil {
			log.Fatal(err)
		}
		err = setCodeEmbeddings(redisearchClient, *embeddings)
		if err != nil {
			log.Fatal(err)
		}
	}

	// hashList is a list of all hashes in the tree.
	hashList := make([]string, 0, len(hashes))
	for k := range hashes {
		hashList = append(hashList, k)
	}
	// store the list of hashes as a comma-separated string for the tree (embedding:tree:<tree-hash>)
	err = redisClient.Set(context.Background(), fmt.Sprintf("embedding:tree:%s", tree.Hash), strings.Join(hashList, ","), 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	// set the embedding:tree equal to the latest tree hash that was processed
	err = redisClient.Set(context.Background(), "embedding:tree", tree.Hash.String(), 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func gitAuth(c client.Client, r *v1alpha1.Repository) (githttp.AuthMethod, error) {
	switch r.Spec.Type {
	case v1alpha1.RepositoryTypeGithub:
		// get the secret
		secret := &corev1.Secret{}
		if err := c.Get(context.TODO(), client.ObjectKey{Name: r.Name, Namespace: r.Namespace}, secret); err != nil {
			return nil, err
		}

		tokenBytes, ok := secret.Data["token"]
		if !ok {
			return nil, errors.New("token not found in secret")
		}
		token := string(tokenBytes)

		return &githttp.BasicAuth{Username: "encoder-run", Password: token}, nil
	default:
		return nil, fmt.Errorf("unsupported repository type: %s", r.Spec.Type)
	}
}

func redisStorageStorer(c client.Client, s *v1alpha1.Storage, nsPrefix string) (storage.Storer, *redis.Options, error) {
	// get the password from the redis secret
	secret := &corev1.Secret{}
	if err := c.Get(context.TODO(), client.ObjectKey{Name: s.Name, Namespace: s.Namespace}, secret); err != nil {
		return nil, nil, err
	}

	passwordBytes, ok := secret.Data["password"]
	if !ok {
		return nil, nil, errors.New("password not found in secret")
	}
	password := string(passwordBytes)

	url := fmt.Sprintf("%s.%s.svc.cluster.local:6379", s.Name, s.Namespace)
	// New redis storage
	opts := &redis.Options{
		Addr:     url,
		Password: password,
		DB:       0,
	}
	return rediscache.NewStorage(opts, nsPrefix), opts, nil
}

func namespace() (string, error) {
	ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", fmt.Errorf("failed to read namespace from service account: %w", err)
	}
	return string(ns), nil
}

func getRedisearchClient(opts *redis.Options, index string) *redisearch.Client {
	pool := &redigoredis.Pool{Dial: func() (redigoredis.Conn, error) {
		return redigoredis.Dial("tcp", opts.Addr, redigoredis.DialPassword(opts.Password))
	}}

	return redisearch.NewClientFromPool(pool, index)
}

func createIndex(r *redisearch.Client) error {
	// Create a schema for the index
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("file_hash")).
		AddField(redisearch.NewNumericField("chunk_id")).
		AddField(redisearch.NewNumericField("start_line")).
		AddField(redisearch.NewNumericField("end_line")).
		AddField(redisearch.NewNumericField("start_column")).
		AddField(redisearch.NewNumericField("end_column")).
		AddField(redisearch.NewVectorFieldOptions("embedding", redisearch.VectorFieldOptions{
			Algorithm: redisearch.Flat,
			Attributes: map[string]interface{}{
				"TYPE":            "FLOAT32",
				"DIM":             768, // Adjust this to the dimension of your embeddings
				"DISTANCE_METRIC": "COSINE",
			},
		}))

	indexDef := redisearch.NewIndexDefinition().AddPrefix("embedding")

	info, _ := r.Info()
	if info == nil {
		// Create the index with the schema
		fmt.Printf("Creating index\n")
		if err := r.CreateIndexWithIndexDefinition(sc, indexDef); err != nil {
			return err
		}
		return nil
	}
	fmt.Printf("Index already exists\n")
	return nil
}

func setCodeEmbeddings(r *redisearch.Client, embeddings CodeEmbeddingsResponse) error {
	// Use a redis pipeline to set all embeddings in one go
	docs := make([]redisearch.Document, 0)
	for _, embs := range embeddings.Results {
		for _, emb := range embs.Embeddings {
			// Create a unique key for each embedding
			key := fmt.Sprintf("embedding:%s:%s:%d", "code", emb.FileHash, emb.ChunkID)
			doc := redisearch.NewDocument(key, 1.0)
			doc.Set("file_hash", emb.FileHash)
			doc.Set("chunk_id", emb.ChunkID)
			doc.Set("start_line", emb.StartLine)
			doc.Set("end_line", emb.EndLine)
			doc.Set("start_column", emb.StartColumn)
			doc.Set("end_column", emb.EndColumn)
			// Convert embedding float slice to bytes
			buf := new(bytes.Buffer)
			for _, val := range emb.Embedding {
				if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
					return err // Handle error appropriately
				}
			}
			doc.Set("embedding", buf.Bytes())
			docs = append(docs, doc)
		}
	}

	if err := r.IndexOptions(redisearch.IndexingOptions{
		Replace: true,
	}, docs...); err != nil {
		return err
	}

	return nil
}

func newScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)
	return scheme
}

func defaultClient() (client.Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		// Load kubeconfig from default location or specified path
		kubeConfigPath := clientcmd.RecommendedHomeFile
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
			&clientcmd.ConfigOverrides{},
		).ClientConfig()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load in-cluster configuration")
		}
	}

	// Create a new scheme and register the API types
	scheme := newScheme()

	// Create the controller-runtime client using the impersonated rest.Config
	c, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to create controller-runtime client: %v", err)
	}

	return c, nil
}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
