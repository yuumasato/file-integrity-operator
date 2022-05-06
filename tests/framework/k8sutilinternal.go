package framework

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// GetKubeconfigAndNamespace returns the *rest.Config and default namespace defined in the
// kubeconfig at the specified path. If no path is provided, returns the default *rest.Config
// and namespace
func GetKubeconfigAndNamespace(configPath string) (*rest.Config, string, error) {
	var clientConfig clientcmd.ClientConfig
	var apiConfig *clientcmdapi.Config
	var err error
	if configPath != "" {
		apiConfig, err = clientcmd.LoadFromFile(configPath)
		if err != nil {
			return nil, "", fmt.Errorf("failed to load user provided kubeconfig: %v", err)
		}
	} else {
		apiConfig, err = clientcmd.NewDefaultClientConfigLoadingRules().Load()
		if err != nil {
			return nil, "", fmt.Errorf("failed to get kubeconfig: %v", err)
		}
	}
	clientConfig = clientcmd.NewDefaultClientConfig(*apiConfig, &clientcmd.ConfigOverrides{})
	kubeconfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, "", err
	}
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		return nil, "", err
	}
	return kubeconfig, namespace, nil
}

// GetDisplayName turns a project dir name in any of {snake, chain, camel}
// cases, hierarchical dot structure, or space-delimited into a
// space-delimited, title'd display name.
// Ex. "another-_AppOperator_againTwiceThrice More"
// ->  "Another App Operator Again Twice Thrice More"
func GetDisplayName(name string) string {
	for _, sep := range ".-_ " {
		splitName := strings.Split(name, string(sep))
		for i := 0; i < len(splitName); i++ {
			if splitName[i] == "" {
				splitName = append(splitName[:i], splitName[i+1:]...)
				i--
			} else {
				splitName[i] = strings.TrimSpace(splitName[i])
			}
		}
		name = strings.Join(splitName, " ")
	}
	splitName := strings.Split(name, " ")
	for i, word := range splitName {
		temp := word
		o := 0
		for j, r := range word {
			if unicode.IsUpper(r) {
				if j > 0 && !unicode.IsUpper(rune(word[j-1])) {
					index := j + o
					temp = temp[0:index] + " " + temp[index:]
					o++
				}
			}
		}
		splitName[i] = temp
	}
	return strings.TrimSpace(strings.Title(strings.Join(splitName, " ")))
}

// GetTypeMetaFromBytes gets the type and object metadata from b. b is assumed
// to be a single Kubernetes resource manifest.
func GetTypeMetaFromBytes(b []byte) (t metav1.TypeMeta, err error) {
	u := unstructured.Unstructured{}
	r := bytes.NewReader(b)
	dec := yaml.NewYAMLOrJSONDecoder(r, 8)
	// There is only one YAML doc if there are no more bytes to be read or EOF
	// is hit.
	if err := dec.Decode(&u); err == nil && r.Len() != 0 {
		return t, errors.New("error getting TypeMeta from bytes: more than one manifest in file")
	} else if err != nil && err != io.EOF {
		return t, fmt.Errorf("error getting TypeMeta from bytes: %v", err)
	}
	return metav1.TypeMeta{
		APIVersion: u.GetAPIVersion(),
		Kind:       u.GetKind(),
	}, nil
}

// dns1123LabelRegexp defines the character set allowed in a DNS 1123 label.
var dns1123LabelRegexp = regexp.MustCompile("[^a-zA-Z0-9]+")

// FormatOperatorNameDNS1123 ensures name is DNS1123 label-compliant by
// replacing all non-compliant UTF-8 characters with "-".
func FormatOperatorNameDNS1123(name string) string {
	if len(validation.IsDNS1123Label(name)) != 0 {
		return dns1123LabelRegexp.ReplaceAllString(name, "-")
	}
	return name
}

// TrimDNS1123Label trims a label to meet the DNS 1123 label length requirement
// by removing characters from the beginning of label such that len(label) <= 63.
func TrimDNS1123Label(label string) string {
	if len(label) > validation.DNS1123LabelMaxLength {
		return label[len(label)-validation.DNS1123LabelMaxLength:]
	}
	return label
}
