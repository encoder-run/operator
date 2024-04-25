package common

import "fmt"

const (
	// AdminClientKey is the key used to store the admin client in the context
	AdminClientKey = "adminClient"
)

func ModelServiceURL(modelId string, namespace string) string {
	return fmt.Sprintf("http://%s-predictor-default.%s.svc.cluster.local:80/v1/models/custom-model:predict", modelId, namespace)
}

func RedisServiceURL(redisId string, namespace string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:6379", redisId, namespace)
}
