package middleware

import (
	"context"
	"fmt"
	"net/http"

	cloudv1alpha1 "github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/common"
	"github.com/gorilla/mux"
	"github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sClientManager struct{}

func (km *K8sClientManager) AdminClient(namespace string) (client.Client, error) {
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

	config.Impersonate.UserName = "system:serviceaccount:" + namespace + ":" + "admin"

	// Create a new scheme and register the API types
	scheme := NewScheme()

	// Create the controller-runtime client using the impersonated rest.Config
	c, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to create controller-runtime client: %v", err)
	}

	return c, nil
}

func NewScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	// add corev1 to the scheme
	_ = corev1.AddToScheme(scheme)
	_ = cloudv1alpha1.AddToScheme(scheme)
	_ = v1beta1.AddToScheme(scheme)

	return scheme
}

func K8sImpersonationMiddleware(km *K8sClientManager) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Initialize the Kubernetes admin client
			client, err := km.AdminClient("default")
			if err != nil {
				http.Error(w, "Failed to initialize Kubernetes client", http.StatusInternalServerError)
				return
			}

			// Attach the client to the context
			ctx := context.WithValue(r.Context(), common.AdminClientKey, client)

			// Proceed with the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
