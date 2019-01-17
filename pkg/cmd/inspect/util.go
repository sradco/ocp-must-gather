package inspect

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"

	configv1 "github.com/openshift/api/config/v1"
)

// resourceContext is used to keep track of previously seen objects
type resourceContext struct {
	visited sets.String
}

func NewResourceContext() *resourceContext {
	return &resourceContext{
		visited: sets.NewString(),
	}
}

func objectReferenceToResourceInfo(clientGetter genericclioptions.RESTClientGetter, ref *configv1.ObjectReference) (*resource.Info, error) {
	b := resource.NewBuilder(clientGetter).
		Unstructured().
		ResourceTypeOrNameArgs(false, fmt.Sprintf("%s/%s", ref.Resource, ref.Name)).
		NamespaceParam(ref.Namespace).
		Flatten().
		Latest()

	infos, err := b.Do().Infos()
	if err != nil {
		return nil, err
	}

	return infos[0], nil
}

// infoToContextKey receives a resource.Info and returns a unique string for use in keeping track of objects previously seen
func infoToContextKey(info *resource.Info) string {
	return fmt.Sprintf("%s/%s/%s/%s", info.Namespace, info.ResourceMapping().GroupVersionKind.Group, info.ResourceMapping().Resource.Resource, info.Name)
}

// objectRefToContextKey is a variant of infoToContextKey that receives a configv1.ObjectReference and returns a unique string for use in keeping track of object references previously seen
func objectRefToContextKey(objRef *configv1.ObjectReference) string {
	return fmt.Sprintf("%s/%s/%s/%s", objRef.Namespace, objRef.Group, objRef.Resource, objRef.Name)
}
