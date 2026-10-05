package util

import (
	"time"

	log "github.com/sirupsen/logrus"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/argoproj/argo/pkg/apis/workflow"
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	wflisters "github.com/argoproj/argo/pkg/client/listers/workflow/v1alpha1"
	unstructutil "github.com/argoproj/argo/util/unstructured"
)

func newDynamicInformer(cfg *rest.Config, resource, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	dclient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	groupVersionResource := schema.GroupVersionResource{
		Group:    workflow.Group,
		Version:  "v1alpha1",
		Resource: resource,
	}
	return unstructutil.NewFilteredUnstructuredInformer(groupVersionResource, dclient, namespace, resyncPeriod, indexers, tweakListOptions)
}

// NewCronWorkflowInformer returns an unstructured informer for CronWorkflows. Individual objects that fail typed conversion do not prevent the informer from listing valid resources.
func NewCronWorkflowInformer(cfg *rest.Config, namespace string, resyncPeriod time.Duration, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return newDynamicInformer(cfg, workflow.CronWorkflowPlural, namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, tweakListOptions)
}

// NewWorkflowTemplateInformer returns an unstructured informer for WorkflowTemplates. Individual objects that fail typed conversion do not prevent the informer from listing valid resources.
func NewWorkflowTemplateInformer(cfg *rest.Config, namespace string, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return newDynamicInformer(cfg, workflow.WorkflowTemplatePlural, namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, nil)
}

// NewClusterWorkflowTemplateInformer returns an unstructured informer for ClusterWorkflowTemplates. Individual objects that fail typed conversion do not prevent the informer from listing valid resources.
func NewClusterWorkflowTemplateInformer(cfg *rest.Config, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return newDynamicInformer(cfg, workflow.ClusterWorkflowTemplatePlural, "", resyncPeriod, cache.Indexers{}, nil)
}

func objectFromUnstructured(un *unstructured.Unstructured, obj runtime.Object) error {
	return runtime.DefaultUnstructuredConverter.FromUnstructured(un.Object, obj)
}

// CronWorkflowFromUnstructured converts an unstructured object to a CronWorkflow.
func CronWorkflowFromUnstructured(un *unstructured.Unstructured) (*wfv1.CronWorkflow, error) {
	var cronWorkflow wfv1.CronWorkflow
	err := objectFromUnstructured(un, &cronWorkflow)
	return &cronWorkflow, err
}

// WorkflowTemplateFromUnstructured converts an unstructured object to a WorkflowTemplate.
func WorkflowTemplateFromUnstructured(un *unstructured.Unstructured) (*wfv1.WorkflowTemplate, error) {
	var workflowTemplate wfv1.WorkflowTemplate
	err := objectFromUnstructured(un, &workflowTemplate)
	return &workflowTemplate, err
}

// ClusterWorkflowTemplateFromUnstructured converts an unstructured object to a ClusterWorkflowTemplate.
func ClusterWorkflowTemplateFromUnstructured(un *unstructured.Unstructured) (*wfv1.ClusterWorkflowTemplate, error) {
	var clusterWorkflowTemplate wfv1.ClusterWorkflowTemplate
	err := objectFromUnstructured(un, &clusterWorkflowTemplate)
	return &clusterWorkflowTemplate, err
}

// WorkflowTemplateInformer combines an unstructured informer with a typed lister.
type WorkflowTemplateInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() wflisters.WorkflowTemplateLister
}

type workflowTemplateInformer struct {
	informer cache.SharedIndexInformer
	lister   wflisters.WorkflowTemplateLister
}

func NewWorkflowTemplateInformerWithLister(cfg *rest.Config, namespace string, resyncPeriod time.Duration) WorkflowTemplateInformer {
	informer := NewWorkflowTemplateInformer(cfg, namespace, resyncPeriod)
	return &workflowTemplateInformer{
		informer: informer,
		lister:   newWorkflowTemplateLister(informer),
	}
}

func (i *workflowTemplateInformer) Informer() cache.SharedIndexInformer {
	return i.informer
}

func (i *workflowTemplateInformer) Lister() wflisters.WorkflowTemplateLister {
	return i.lister
}

// ClusterWorkflowTemplateInformer combines an unstructured informer with a typed lister.
type ClusterWorkflowTemplateInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() wflisters.ClusterWorkflowTemplateLister
}

type clusterWorkflowTemplateInformer struct {
	informer cache.SharedIndexInformer
	lister   wflisters.ClusterWorkflowTemplateLister
}

func NewClusterWorkflowTemplateInformerWithLister(cfg *rest.Config, resyncPeriod time.Duration) ClusterWorkflowTemplateInformer {
	informer := NewClusterWorkflowTemplateInformer(cfg, resyncPeriod)
	return &clusterWorkflowTemplateInformer{
		informer: informer,
		lister:   newClusterWorkflowTemplateLister(informer),
	}
}

func (i *clusterWorkflowTemplateInformer) Informer() cache.SharedIndexInformer {
	return i.informer
}

func (i *clusterWorkflowTemplateInformer) Lister() wflisters.ClusterWorkflowTemplateLister {
	return i.lister
}

type unstructuredWorkflowTemplateLister struct {
	informer cache.SharedIndexInformer
}

func newWorkflowTemplateLister(informer cache.SharedIndexInformer) wflisters.WorkflowTemplateLister {
	return &unstructuredWorkflowTemplateLister{informer: informer}
}

func (l *unstructuredWorkflowTemplateLister) List(selector labels.Selector) ([]*wfv1.WorkflowTemplate, error) {
	workflowTemplates := make([]*wfv1.WorkflowTemplate, 0)
	for _, object := range l.informer.GetStore().List() {
		unstructuredObject := object.(*unstructured.Unstructured)
		if !selector.Matches(labels.Set(unstructuredObject.GetLabels())) {
			continue
		}
		workflowTemplate, err := WorkflowTemplateFromUnstructured(unstructuredObject)
		if err != nil {
			log.Warnf("Failed to unmarshal WorkflowTemplate %v: %v", unstructuredObject, err)
			continue
		}
		workflowTemplates = append(workflowTemplates, workflowTemplate)
	}
	return workflowTemplates, nil
}

func (l *unstructuredWorkflowTemplateLister) WorkflowTemplates(namespace string) wflisters.WorkflowTemplateNamespaceLister {
	return &unstructuredWorkflowTemplateNamespaceLister{informer: l.informer, namespace: namespace}
}

type unstructuredWorkflowTemplateNamespaceLister struct {
	informer  cache.SharedIndexInformer
	namespace string
}

func (l *unstructuredWorkflowTemplateNamespaceLister) List(selector labels.Selector) ([]*wfv1.WorkflowTemplate, error) {
	workflowTemplates := make([]*wfv1.WorkflowTemplate, 0)
	for _, object := range l.informer.GetStore().List() {
		unstructuredObject := object.(*unstructured.Unstructured)
		if unstructuredObject.GetNamespace() != l.namespace || !selector.Matches(labels.Set(unstructuredObject.GetLabels())) {
			continue
		}
		workflowTemplate, err := WorkflowTemplateFromUnstructured(unstructuredObject)
		if err != nil {
			log.Warnf("Failed to unmarshal WorkflowTemplate %v: %v", unstructuredObject, err)
			continue
		}
		workflowTemplates = append(workflowTemplates, workflowTemplate)
	}
	return workflowTemplates, nil
}

func (l *unstructuredWorkflowTemplateNamespaceLister) Get(name string) (*wfv1.WorkflowTemplate, error) {
	object, exists, err := l.informer.GetStore().GetByKey(l.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, newNotFoundError(workflow.WorkflowTemplateSingular, name)
	}
	return WorkflowTemplateFromUnstructured(object.(*unstructured.Unstructured))
}

type unstructuredClusterWorkflowTemplateLister struct {
	informer cache.SharedIndexInformer
}

func newClusterWorkflowTemplateLister(informer cache.SharedIndexInformer) wflisters.ClusterWorkflowTemplateLister {
	return &unstructuredClusterWorkflowTemplateLister{informer: informer}
}

func (l *unstructuredClusterWorkflowTemplateLister) List(selector labels.Selector) ([]*wfv1.ClusterWorkflowTemplate, error) {
	clusterWorkflowTemplates := make([]*wfv1.ClusterWorkflowTemplate, 0)
	for _, object := range l.informer.GetStore().List() {
		unstructuredObject := object.(*unstructured.Unstructured)
		if !selector.Matches(labels.Set(unstructuredObject.GetLabels())) {
			continue
		}
		clusterWorkflowTemplate, err := ClusterWorkflowTemplateFromUnstructured(unstructuredObject)
		if err != nil {
			log.Warnf("Failed to unmarshal ClusterWorkflowTemplate %v: %v", unstructuredObject, err)
			continue
		}
		clusterWorkflowTemplates = append(clusterWorkflowTemplates, clusterWorkflowTemplate)
	}
	return clusterWorkflowTemplates, nil
}

func (l *unstructuredClusterWorkflowTemplateLister) Get(name string) (*wfv1.ClusterWorkflowTemplate, error) {
	object, exists, err := l.informer.GetStore().GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, newNotFoundError(workflow.ClusterWorkflowTemplateSingular, name)
	}
	return ClusterWorkflowTemplateFromUnstructured(object.(*unstructured.Unstructured))
}

func newNotFoundError(resource, name string) error {
	return apierr.NewNotFound(schema.GroupResource{Group: workflow.Group, Resource: resource}, name)
}
