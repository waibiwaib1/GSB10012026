package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

func TestWorkflowTemplateListerSkipsMalformedObjects(t *testing.T) {
	informer := cache.NewSharedIndexInformer(nil, &unstructured.Unstructured{}, 0, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	require.NoError(t, informer.GetStore().Add(&unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "argoproj.io/v1alpha1",
		"kind":       "WorkflowTemplate",
		"metadata": map[string]interface{}{
			"name":      "valid",
			"namespace": "argo",
		},
		"spec": map[string]interface{}{
			"templates": []interface{}{map[string]interface{}{"name": "valid-template"}},
		},
	}}))
	require.NoError(t, informer.GetStore().Add(&unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "argoproj.io/v1alpha1",
		"kind":       "WorkflowTemplate",
		"metadata": map[string]interface{}{
			"name":      "malformed",
			"namespace": "argo",
		},
		"spec": map[string]interface{}{"arguments": map[string]interface{}{"parameters": map[string]interface{}{"name": "invalid"}}},
	}}))

	lister := newWorkflowTemplateLister(informer)
	result, err := lister.WorkflowTemplates("argo").List(labels.Everything())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "valid", result[0].Name)

	_, err = lister.WorkflowTemplates("argo").Get("malformed")
	assert.Error(t, err)
}

func TestClusterWorkflowTemplateListerSkipsMalformedObjects(t *testing.T) {
	informer := cache.NewSharedIndexInformer(nil, &unstructured.Unstructured{}, 0, cache.Indexers{})
	require.NoError(t, informer.GetStore().Add(&unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "argoproj.io/v1alpha1",
		"kind":       "ClusterWorkflowTemplate",
		"metadata": map[string]interface{}{
			"name": "valid",
		},
		"spec": map[string]interface{}{
			"templates": []interface{}{map[string]interface{}{"name": "valid-template"}},
		},
	}}))
	require.NoError(t, informer.GetStore().Add(&unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "argoproj.io/v1alpha1",
		"kind":       "ClusterWorkflowTemplate",
		"metadata": map[string]interface{}{
			"name": "malformed",
		},
		"spec": map[string]interface{}{"arguments": map[string]interface{}{"parameters": map[string]interface{}{"name": "invalid"}}},
	}}))

	lister := newClusterWorkflowTemplateLister(informer)
	result, err := lister.List(labels.Everything())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "valid", result[0].Name)

	_, err = lister.Get("malformed")
	assert.Error(t, err)
}

func TestCronWorkflowFromUnstructured(t *testing.T) {
	unstructuredCronWorkflow := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "argoproj.io/v1alpha1",
		"kind":       "CronWorkflow",
		"metadata": map[string]interface{}{
			"name":      "cron",
			"namespace": "argo",
		},
		"spec": map[string]interface{}{
			"schedule": "0 * * * *",
			"workflowSpec": map[string]interface{}{
				"templates": []interface{}{map[string]interface{}{"name": "cron-template"}},
			},
		},
	}}

	cronWorkflow, err := CronWorkflowFromUnstructured(unstructuredCronWorkflow)
	require.NoError(t, err)
	assert.Equal(t, metav1.ObjectMeta{Name: "cron", Namespace: "argo"}, cronWorkflow.ObjectMeta)
	assert.Equal(t, "0 * * * *", cronWorkflow.Spec.Schedule)
	assert.IsType(t, wfv1.WorkflowSpec{}, cronWorkflow.Spec.WorkflowSpec)
}
