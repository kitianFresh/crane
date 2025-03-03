package base

import (
	"fmt"

	analysisapi "github.com/gocrane/api/analysis/v1alpha1"

	"github.com/gocrane/crane/pkg/recommendation/framework"
	"github.com/gocrane/crane/pkg/utils"
)

// Filter out k8s resources that are not supported by the recommender.
func (br *BaseRecommender) Filter(ctx *framework.RecommendationContext) error {
	// 1. get object identity
	identity := ctx.Identity

	// 2. load recommender accepted kubernetes object
	accepted := br.Recommender.AcceptedResourceSelectors

	// 3. if not support, abort the recommendation flow
	supported := IsIdentitySupported(identity, accepted)
	if !supported {
		return fmt.Errorf("recommender %s is failed at fliter, your kubernetes resource is not supported for recommender %s.", br.Name(), br.Name())
	}

	return nil
}

// IsIdentitySupported check whether object identity fit resource selector.
func IsIdentitySupported(identity framework.ObjectIdentity, selectors []analysisapi.ResourceSelector) bool {
	supported := false
	for _, selector := range selectors {
		if len(selector.Name) == 0 {
			if selector.Kind == identity.Kind && selector.APIVersion == identity.APIVersion {
				labelMatch, _ := utils.LabelSelectorMatched(identity.Labels, selector.LabelSelector)
				if labelMatch {
					return true
				}
			}
		} else if selector.Kind == identity.Kind && selector.APIVersion == identity.APIVersion && selector.Name == identity.Name {
			labelMatch, _ := utils.LabelSelectorMatched(identity.Labels, selector.LabelSelector)
			if labelMatch {
				return true
			}
		}
	}

	return supported
}
