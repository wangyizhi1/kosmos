package utils

import (
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	jsonpatch1 "github.com/mattbaird/jsonpatch"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	kosmosversioned "github.com/kosmos.io/kosmos/pkg/generated/clientset/versioned"
)

type ClustersNodeSelection struct {
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Affinity     *corev1.Affinity    `json:"affinity,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
}

func CreateMergePatch(original, new interface{}) ([]byte, error) {
	pvByte, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}
	cloneByte, err := json.Marshal(new)
	if err != nil {
		return nil, err
	}
	patch, err := jsonpatch.CreateMergePatch(pvByte, cloneByte)
	if err != nil {
		return nil, err
	}
	return patch, nil
}

func CreateJSONPatch(original, new interface{}) ([]byte, error) {
	pvByte, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}
	cloneByte, err := json.Marshal(new)
	if err != nil {
		return nil, err
	}
	patchs, err := jsonpatch1.CreatePatch(pvByte, cloneByte)
	if err != nil {
		return nil, err
	}
	patchBytes, err := json.Marshal(patchs)
	if err != nil {
		return nil, err
	}
	return patchBytes, nil
}

type Opts func(*rest.Config)

func NewClient(configPath string, opts ...Opts) (kubernetes.Interface, error) {
	var (
		config *rest.Config
		err    error
	)
	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("could not read config file for cluster: %v", err)
		}
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create client for root cluster: %v", err)
	}
	return client, nil
}

func NewClientFromByte(kubeConfig []byte, opts ...Opts) (kubernetes.Interface, error) {
	var (
		config *rest.Config
		err    error
	)

	clientconfig, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}
	config, err = clientconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create client for root cluster: %v", err)
	}
	return client, nil
}

func NewConfigFromBytes(kubeConfig []byte, opts ...Opts) (*rest.Config, error) {
	var (
		config *rest.Config
		err    error
	)

	c, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}
	config, err = c.ClientConfig()
	if err != nil {
		return nil, err
	}

	for _, h := range opts {
		if h == nil {
			continue
		}
		h(config)
	}

	return config, nil
}

func NewClientFromConfigPath(configPath string, opts ...Opts) (kubernetes.Interface, error) {
	var (
		config *rest.Config
		err    error
	)
	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from configpath: %v", err)
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create clientset: %v", err)
	}
	return client, nil
}

func NewKosmosClientFromConfigPath(configPath string, opts ...Opts) (kosmosversioned.Interface, error) {
	var (
		config *rest.Config
		err    error
	)
	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from configpath: %v", err)
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	client, err := kosmosversioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create clientset: %v", err)
	}
	return client, nil
}

func NewClientFromBytes(kubeConfig []byte, opts ...Opts) (kubernetes.Interface, error) {
	var (
		config *rest.Config
		err    error
	)

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}
	config, err = clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create client failed: %v", err)
	}
	return client, nil
}

func NewKosmosClientFromBytes(kubeConfig []byte, opts ...Opts) (kosmosversioned.Interface, error) {
	var (
		config *rest.Config
		err    error
	)

	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}
	config, err = clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	client, err := kosmosversioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create client failed: %v", err)
	}
	return client, nil
}

func NewMetricClient(configPath string, opts ...Opts) (versioned.Interface, error) {
	var (
		config *rest.Config
		err    error
	)
	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("could not read config file for cluster: %v", err)
		}
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	metricClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create client for root cluster: %v", err)
	}
	return metricClient, nil
}

func NewMetricClientFromByte(kubeConfig []byte, opts ...Opts) (versioned.Interface, error) {
	var (
		config *rest.Config
		err    error
	)

	clientconfig, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}
	config, err = clientconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	metricClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create client for root cluster: %v", err)
	}
	return metricClient, nil
}

func IsVirtualNode(node *corev1.Node) bool {
	if node == nil {
		return false
	}
	valStr, exist := node.ObjectMeta.Labels[KosmosNodeLabel]
	if !exist {
		return false
	}
	return valStr == KosmosNodeValue
}

// IsVirtualPod
// TODO change func name
func IsVirtualPod(pod *corev1.Pod) bool {
	if pod.Labels != nil && pod.Labels[KosmosPodLabel] == "true" {
		return true
	}
	return false
}

func UpdateConfigMap(old, new *corev1.ConfigMap) {
	old.Labels = new.Labels
	old.Data = new.Data
	old.BinaryData = new.BinaryData
}

func UpdateSecret(old, new *corev1.Secret) {
	old.Labels = new.Labels
	old.Data = new.Data
	old.StringData = new.StringData
	old.Type = new.Type
}

func UpdateUnstructured(old, new *unstructured.Unstructured, g func() (interface{}, error), cb func(old, new interface{}) error) (*unstructured.Unstructured, error) {
	oldObj, err := g()
	if err != nil {
		return nil, err
	}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(old.UnstructuredContent(), &oldObj); err != nil {
		return nil, err
	}

	newObj, err := g()
	if err != nil {
		return nil, err
	}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(new.UnstructuredContent(), &newObj); err != nil {
		return nil, err
	}

	if err := cb(oldObj, newObj); err != nil {
		return nil, err
	}

	if retObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldObj); err == nil {
		return &unstructured.Unstructured{
			Object: retObj,
		}, nil
	} else {
		return nil, err
	}
}

func IsObjectGlobal(obj *metav1.ObjectMeta) bool {
	if obj.Annotations == nil {
		return false
	}

	if obj.Annotations[KosmosGlobalLabel] == "true" {
		return true
	}

	return false
}

func IsObjectUnstructuredGlobal(obj map[string]string) bool {
	if obj == nil {
		return false
	}

	if obj[KosmosGlobalLabel] == "true" {
		return true
	}

	return false
}
