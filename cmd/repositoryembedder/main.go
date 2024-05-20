package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
	v1alpha1 "github.com/encoder-run/operator/api/cloud/v1alpha1"
	postgrescache "github.com/encoder-run/operator/pkg/cache/postgres"
	rediscache "github.com/encoder-run/operator/pkg/cache/redis"
	"github.com/encoder-run/operator/pkg/database"
	"github.com/encoder-run/operator/pkg/embedder"
	"github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage"
	redigoredis "github.com/gomodule/redigo/redis"
	"github.com/pgvector/pgvector-go"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
	// Initialize embClient
	embClient := embedder.NewClient(modelId, ns)

	// Repository URL for remote git repository
	var url string
	var branch string
	switch repo.Spec.Type {
	case v1alpha1.RepositoryTypeGithub:
		url = repo.Spec.Github.URL
		branch = repo.Spec.Github.Branch
	default:
		CheckIfError(fmt.Errorf("unsupported repository type: %s", repo.Spec.Type))
	}

	var storer storage.Storer
	var redisearchClient *redisearch.Client
	var redisClient *redis.Client
	var db *gorm.DB
	switch st.Spec.Type {
	case v1alpha1.StorageTypeRedis:
		// Get the go-git storage storer based on the storage type
		s, opts, err := redisStorageStorer(c, st, url)
		if err != nil {
			CheckIfError(err)
		}
		storer = s
		redisearchClient = getRedisearchClient(opts, fmt.Sprintf("%s:embedding", url))
		redisClient = redis.NewClient(opts)
		err = createIndex(redisearchClient, url)
		if err != nil {
			CheckIfError(err)
		}
	case v1alpha1.StorageTypePostgres:
		// Get the go-git storage storer based on the storage type
		s, dbClient, err := postgresStorageStorer(c, st, url)
		if err != nil {
			CheckIfError(err)
		}
		db = dbClient
		storer = s
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
				URL:           fmt.Sprintf("https://%s", url),
				Auth:          auth,
				ReferenceName: plumbing.NewBranchReferenceName(branch),
				Depth:         1,
				SingleBranch:  true,
			})
			CheckIfError(err)
		} else {
			CheckIfError(err)
		}
	}
	// Fetch changes from the remote repository
	refSpec := fmt.Sprintf("+refs/heads/%s:refs/remotes/origin/%s", branch, branch)
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(refSpec),
		},
		Auth:  auth,
		Depth: 1,
	})
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}

	// err = r.Prune(git.PruneOptions{})
	// CheckIfError(err)

	// Manually update local branch reference to match the remote tracking branch
	// Typically in a bare repo, this might be done in response to a push or a hook
	remoteRef, err := r.Reference(plumbing.NewRemoteReferenceName("origin", branch), true)
	CheckIfError(err)
	fmt.Printf("Remote Ref: %v\n", remoteRef)

	// // Update local main directly to point to the same commit
	// localRef := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branch), remoteRef.Hash())
	// err = r.Storer.SetReference(localRef)
	// CheckIfError(err)

	commit, err := r.CommitObject(remoteRef.Hash())
	if err != nil {
		log.Fatal(err)
	}

	tree, err := commit.Tree()
	if err != nil {
		log.Fatal(err)
	}

	// Switch based on the db type
	switch st.Spec.Type {
	case v1alpha1.StorageTypeRedis:
		processRedisEmbeddings(embClient, tree, redisClient, redisearchClient, url)
	case v1alpha1.StorageTypePostgres:
		processPostgresEmbeddings(embClient, tree, db, url)
	default:
		CheckIfError(fmt.Errorf("unsupported storage type: %s", st.Spec.Type))
	}
}

func processPostgresEmbeddings(embClient *embedder.EmbeddingClient, tree *object.Tree, db *gorm.DB, url string) {
	// Check for existing processed hashes
	existingHashes := make(map[string]bool)
	// List all rows in the codeembedding table given the url
	var embeddings []database.CodeEmbedding
	if err := db.Where("url = ?", url).Find(&embeddings).Error; err != nil {
		log.Fatalf("failed to query existing embeddings: %v", err)
	}

	// Create a set of unique hashes
	for _, embedding := range embeddings {
		existingHashes[fmt.Sprintf("%s.%s", embedding.FileHash, embedding.FilePath)] = true
	}

	hashes := make(map[string]bool)
	filesBatch := []embedder.CodeEmbeddingRequest{}
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

		hashes[fmt.Sprintf("%s.%s", file.Hash.String(), file.Name)] = true

		if !existingHashes[fmt.Sprintf("%s.%s", file.Hash.String(), file.Name)] {
			content, err := file.Contents()
			if err != nil {
				log.Fatal(err)
			}
			filesBatch = append(filesBatch, embedder.CodeEmbeddingRequest{
				Path:    file.Name,
				Content: content,
				Hash:    file.Hash.String(),
			})
			count++

			if count >= batchSize {
				processAndSaveEmbeddings(embClient, db, &filesBatch, url) // process embeddings
				filesBatch = []embedder.CodeEmbeddingRequest{}            // Reset the batch
				count = 0
			}
		} else {
			fmt.Printf("Skipping file '%s' since its hash is already processed\n", file.Name)
		}
	}

	if len(filesBatch) > 0 {
		processAndSaveEmbeddings(embClient, db, &filesBatch, url) // Process any remaining files
	}
}

func processAndSaveEmbeddings(embClient *embedder.EmbeddingClient, db *gorm.DB, filesBatch *[]embedder.CodeEmbeddingRequest, url string) {
	fmt.Printf("Processing batch of files\n")
	embeddings, err := embClient.FetchEmbeddings(*filesBatch)
	if err != nil {
		log.Fatal(err)
	}
	for filePath, embs := range embeddings.Results {
		for _, emb := range embs.Embeddings {
			newEmb := database.CodeEmbedding{
				URL:        url,
				FileHash:   emb.FileHash,
				FilePath:   filePath,
				ChunkID:    emb.ChunkID,
				StartIndex: emb.StartIndex,
				EndIndex:   emb.EndIndex,
				Embedding:  pgvector.NewVector(emb.Embedding), // Assuming emb.Vector is the appropriate data structure
			}
			// Upsert operation using Clauses with ON CONFLICT
			if err := db.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "chunk_id"}, {Name: "file_hash"}, {Name: "url"}, {Name: "file_path"}}, // Columns part of the unique constraint
				DoUpdates: clause.Assignments(map[string]interface{}{ // Update these fields if there is a conflict
					"start_index": newEmb.StartIndex,
					"end_index":   newEmb.EndIndex,
					"embedding":   newEmb.Embedding,
				}),
			}).Create(&newEmb).Error; err != nil {
				log.Fatalf("failed to save or update embedding: %v", err)
			}
		}
	}
}

func processRedisEmbeddings(embClient *embedder.EmbeddingClient, tree *object.Tree, redisClient *redis.Client, redisearchClient *redisearch.Client, url string) {
	// Check for existing processed hashes
	existingHashes := make(map[string]bool)
	existingTreeHash, err := redisClient.Get(context.Background(), fmt.Sprintf("%s:embedding:tree", url)).Result()
	if err != nil && err != redis.Nil {
		log.Fatal(err)
	}
	if err == nil {
		existingHashesStr, err := redisClient.Get(context.Background(), fmt.Sprintf("%s:embedding:tree:%s", url, existingTreeHash)).Result()
		if err != nil && err != redis.Nil {
			log.Fatal(err)
		}
		for _, h := range strings.Split(existingHashesStr, ",") {
			existingHashes[h] = true
		}
	}

	hashes := make(map[string]bool)
	filesBatch := []embedder.CodeEmbeddingRequest{}
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

			filesBatch = append(filesBatch, embedder.CodeEmbeddingRequest{
				Path:    file.Name,
				Content: content,
				Hash:    file.Hash.String(),
			})
			count++

			// Process in batches of 10
			if count >= batchSize {
				fmt.Printf("Processing batch of files\n")
				embeddings, err := embClient.FetchEmbeddings(filesBatch)
				if err != nil {
					log.Fatal(err)
				}
				err = setCodeEmbeddings(redisearchClient, *embeddings, url)
				if err != nil {
					log.Fatal(err)
				}
				// Reset the batch
				filesBatch = []embedder.CodeEmbeddingRequest{}
				count = 0
			}
		} else {
			fmt.Printf("Skipping file '%s' since its hash is already processed\n", file.Name)
		}
	}

	// Process any remaining files
	if len(filesBatch) > 0 {
		fmt.Printf("Processing remaining files\n")
		embeddings, err := embClient.FetchEmbeddings(filesBatch)
		if err != nil {
			log.Fatal(err)
		}
		err = setCodeEmbeddings(redisearchClient, *embeddings, url)
		if err != nil {
			log.Fatal(err)
		}
	}

	// hashList is a list of all hashes in the tree.
	hashList := make([]string, 0, len(hashes))
	for k := range hashes {
		hashList = append(hashList, k)
	}
	// store the list of hashes as a comma-separated string for the tree (embedding:ns:tree:<tree-hash>)
	err = redisClient.Set(context.Background(), fmt.Sprintf("%s:embedding:tree:%s", url, tree.Hash.String()), strings.Join(hashList, ","), 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	// set the embedding:ns:tree equal to the latest tree hash that was processed
	err = redisClient.Set(context.Background(), fmt.Sprintf("%s:embedding:tree", url), tree.Hash.String(), 0).Err()
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

func postgresStorageStorer(c client.Client, s *v1alpha1.Storage, nsPrefix string) (storage.Storer, *gorm.DB, error) {
	// get the password from the postgres secret
	secret := &corev1.Secret{}
	if err := c.Get(context.TODO(), client.ObjectKey{Name: s.Name, Namespace: s.Namespace}, secret); err != nil {
		return nil, nil, err
	}

	hostBytes, ok := secret.Data["host"]
	if !ok {
		return nil, nil, errors.New("host not found in secret")
	}
	host := string(hostBytes)

	usernameBytes, ok := secret.Data["username"]
	if !ok {
		return nil, nil, errors.New("username not found in secret")
	}
	username := string(usernameBytes)

	passwordBytes, ok := secret.Data["password"]
	if !ok {
		return nil, nil, errors.New("password not found in secret")
	}
	password := string(passwordBytes)

	databaseBytes, ok := secret.Data["database"]
	if !ok {
		return nil, nil, errors.New("database not found in secret")
	}
	db := string(databaseBytes)

	portBytes, ok := secret.Data["port"]
	if !ok {
		return nil, nil, errors.New("port not found in secret")
	}
	port := string(portBytes)

	sslModeBytes, ok := secret.Data["ssl_mode"]
	if !ok {
		return nil, nil, errors.New("ssl_mode not found in secret")
	}
	sslMode := string(sslModeBytes)

	timezoneBytes, ok := secret.Data["timezone"]
	if !ok {
		return nil, nil, errors.New("timezone not found in secret")
	}
	timezone := string(timezoneBytes)

	// Construct the DSN

	dsn := database.ConstructPostgresDSN(host, username, password, db, port, sslMode, timezone)

	dbClient, err := database.GetPostgresClient(dsn)
	if err != nil {
		return nil, nil, err
	}

	// New postgres storage
	return postgrescache.NewStorage(dbClient, nsPrefix), dbClient, nil
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

func createIndex(r *redisearch.Client, ns string) error {
	// Create a schema for the index
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("file_hash")).
		AddField(redisearch.NewTextField("file_path")).
		AddField(redisearch.NewNumericField("chunk_id")).
		AddField(redisearch.NewNumericField("start_index")).
		AddField(redisearch.NewNumericField("end_index")).
		AddField(redisearch.NewVectorFieldOptions("embedding", redisearch.VectorFieldOptions{
			Algorithm: redisearch.Flat,
			Attributes: map[string]interface{}{
				"TYPE":            "FLOAT32",
				"DIM":             768, // Adjust this to the dimension of your embeddings
				"DISTANCE_METRIC": "COSINE",
			},
		}))

	indexDef := redisearch.NewIndexDefinition().AddPrefix(fmt.Sprintf("%s:embedding:code:", ns))

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

func setCodeEmbeddings(r *redisearch.Client, embeddings embedder.CodeEmbeddingsResponse, namespace string) error {
	// Use a redis pipeline to set all embeddings in one go
	docs := make([]redisearch.Document, 0)
	for filePath, embs := range embeddings.Results {
		for _, emb := range embs.Embeddings {
			// Create a unique key for each embedding
			key := fmt.Sprintf("%s:embedding:%s:%s:%d", namespace, "code", emb.FileHash, emb.ChunkID)
			doc := redisearch.NewDocument(key, 1.0)
			doc.Set("fileHash", emb.FileHash)
			doc.Set("filePath", filePath)
			doc.Set("chunkID", emb.ChunkID)
			doc.Set("startIndex", emb.StartIndex)
			doc.Set("endIndex", emb.EndIndex)
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
