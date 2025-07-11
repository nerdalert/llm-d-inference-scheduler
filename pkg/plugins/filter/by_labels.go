package filter

import (
	"context"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/scheduling/framework"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/scheduling/types"
)

const (
	// ByLabelsFilterType is the type of the ByLabelsFilter
	ByLabelsFilterType = "by-labels"
)

// compile-time type assertion
var _ framework.Filter = &ByLabels{}

// NewByLabel returns a new filter instance, configured with the provided
// name and label selector.
func NewByLabel(name string, selector *metav1.LabelSelector) (framework.Filter, error) {
	if name == "" {
		return nil, errors.New("ByLabels: missing filter name")
	}
	labelSelector, err := metav1.LabelSelectorAsSelector(selector)
	if err != nil {
		return nil, err
	}

	return &ByLabels{
		name:     name,
		selector: labelSelector,
	}, nil
}

// ByLabels filters out pods that do not match its label selector criteria
type ByLabels struct {
	name     string
	selector labels.Selector
}

// Type returns the type of the filter
func (blf *ByLabels) Type() string {
	return ByLabelsFilterType
}

// Name returns the name of the instance of the filter.
func (blf *ByLabels) Name() string {
	return blf.name
}

// WithName sets the name of the filter.
func (blf *ByLabels) WithName(name string) *ByLabels {
	blf.name = name
	return blf
}

// Filter filters out all pods that do not satisfy the label selector
func (blf *ByLabels) Filter(_ context.Context, _ *types.CycleState, _ *types.LLMRequest, pods []types.Pod) []types.Pod {
	filtered := []types.Pod{}

	for _, pod := range pods {
		labels := labels.Set(pod.GetPod().Labels)
		if blf.selector.Matches(labels) {
			filtered = append(filtered, pod)
		}
	}
	return filtered
}
