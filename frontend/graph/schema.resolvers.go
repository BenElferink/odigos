package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/odigos-io/odigos/api/k8sconsts"
	"github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/frontend/graph/model"
	"github.com/odigos-io/odigos/frontend/kube"
	"github.com/odigos-io/odigos/frontend/services"
	actionservices "github.com/odigos-io/odigos/frontend/services/actions"
	"github.com/odigos-io/odigos/frontend/services/describe/odigos_describe"
	"github.com/odigos-io/odigos/frontend/services/describe/source_describe"
	testconnection "github.com/odigos-io/odigos/frontend/services/test_connection"
	"github.com/odigos-io/odigos/k8sutils/pkg/env"
	"github.com/odigos-io/odigos/k8sutils/pkg/pro"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"

	"golang.org/x/sync/errgroup"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// APITokens is the resolver for the apiTokens field.
func (r *computePlatformResolver) APITokens(ctx context.Context, obj *model.ComputePlatform) ([]*model.APIToken, error) {
	ns := env.GetCurrentNamespace()

	// The result should always be 0 or 1:
	// If it's 0, it means this is the OSS version.
	// If it's 1, it means this is the Enterprise version.

	secret, err := kube.DefaultClient.CoreV1().Secrets(ns).Get(ctx, k8sconsts.OdigosProSecretName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return make([]*model.APIToken, 0), nil
		}
		return nil, err
	}

	token := string(secret.Data[k8sconsts.OdigosOnpremTokenSecretKey])

	// Extract the payload from the JWT
	tokenPayload, err := extractJWTPayload(token)
	if err != nil {
		// We don't want to return an error here, because the user may have provided a bad token.
		// Throwing this will prevent the entire CP from being fetched, and prevent the user from being able to update the token...
		// return nil, fmt.Errorf("failed to extract JWT payload: %w", err)

		return []*model.APIToken{
			{
				Token:     token,
				Name:      "ERROR",
				IssuedAt:  0,
				ExpiresAt: 0,
			},
		}, nil
	}

	// Extract values from the token payload
	aud, _ := tokenPayload["aud"].(string)
	iat, _ := tokenPayload["iat"].(float64)
	exp, _ := tokenPayload["exp"].(float64)

	// We need to return an array (even if it's just 1 token), because in the future we will have to support multiple platforms.
	return []*model.APIToken{
		{
			Token:     token,
			Name:      aud,
			IssuedAt:  int(iat) * 1000, // Convert to milliseconds
			ExpiresAt: int(exp) * 1000, // Convert to milliseconds
		},
	}, nil
}

// K8sActualNamespaces is the resolver for the k8sActualNamespaces field.
func (r *computePlatformResolver) K8sActualNamespaces(ctx context.Context, obj *model.ComputePlatform) ([]*model.K8sActualNamespace, error) {
	k8sNamespaces, err := services.GetK8SNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	k8sActualNamespaces := make([]*model.K8sActualNamespace, len(k8sNamespaces.Namespaces))
	for i, namespace := range k8sNamespaces.Namespaces {
		k8sActualNamespaces[i] = &namespace
	}

	return k8sActualNamespaces, nil
}

// K8sActualNamespace is the resolver for the k8sActualNamespace field.
func (r *computePlatformResolver) K8sActualNamespace(ctx context.Context, obj *model.ComputePlatform, name string) (*model.K8sActualNamespace, error) {
	namespace, err := kube.DefaultClient.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	nsName := namespace.Namespace

	// check if entire namespace is instrumented
	crd, err := services.GetSourceCRD(ctx, nsName, nsName, services.WorkloadKindNamespace)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return &model.K8sActualNamespace{}, err
	}
	instrumented := crd != nil && !crd.Spec.DisableInstrumentation

	return &model.K8sActualNamespace{
		Name:     namespace.Name,
		Selected: instrumented,
	}, nil
}

// Sources is the resolver for the sources field.
func (r *computePlatformResolver) Sources(ctx context.Context, obj *model.ComputePlatform, nextPage string, groupName string) (*model.PaginatedSources, error) {
	limit, _ := services.GetPageLimit(ctx)
	icList, err := kube.DefaultClient.OdigosClient.InstrumentationConfigs("").List(ctx, metav1.ListOptions{
		Limit:    int64(limit),
		Continue: nextPage,
	})

	if err != nil {
		if strings.Contains(err.Error(), "The provided continue parameter is too old") {
			// Retry without the continue token
			icList, err = kube.DefaultClient.OdigosClient.InstrumentationConfigs("").List(ctx, metav1.ListOptions{
				Limit: int64(limit),
			})

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Get Source objects to compare with groupName
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(len(icList.Items))

	// Keep order based on idx
	srcList := make([]*v1alpha1.Source, len(icList.Items))
	for idx, ic := range icList.Items {
		g.Go(func() error {
			src, err := services.GetSourceCRD(ctx, ic.Namespace, ic.OwnerReferences[0].Name, services.WorkloadKind(ic.OwnerReferences[0].Kind))
			if err != nil {
				return err
			}

			srcList[idx] = src
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	var actualSources []*model.K8sActualSource
	for idx, ic := range icList.Items {
		// If matches groupName, return the source
		if groupName == "" || srcList[idx].Labels[k8sconsts.SourceGroupLabelPrefix+groupName] == "true" {
			actualSources = append(actualSources, instrumentationConfigToActualSource(ic))
		}
	}

	return &model.PaginatedSources{
		NextPage: icList.GetContinue(),
		Items:    actualSources,
	}, nil
}

// Source is the resolver for the source field.
func (r *computePlatformResolver) Source(ctx context.Context, obj *model.ComputePlatform, sourceID model.K8sSourceID, groupName string) (*model.K8sActualSource, error) {
	ns := sourceID.Namespace
	kind := sourceID.Kind
	name := sourceID.Name

	ic, err := kube.DefaultClient.OdigosClient.InstrumentationConfigs(ns).Get(ctx, workload.CalculateWorkloadRuntimeObjectName(name, string(kind)), metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get InstrumentationConfig: %w", err)
	}
	if ic == nil {
		return nil, fmt.Errorf("InstrumentationConfig not found for %s/%s in namespace %s", kind, name, ns)
	}

	// Get Source object to compare with groupName
	src, err := services.GetSourceCRD(ctx, ns, name, services.WorkloadKind(kind))
	if err != nil {
		return nil, fmt.Errorf("failed to get Source: %w", err)
	}

	// If matches groupName, return the source
	if groupName == "" || src.Labels[k8sconsts.SourceGroupLabelPrefix+groupName] == "true" {
		payload := instrumentationConfigToActualSource(*ic)
		condition, _ := services.GetInstrumentationInstancesHealthCondition(ctx, ns, name, string(kind))
		if condition.Status != "" {
			payload.Conditions = append(payload.Conditions, &condition)
		}

		return payload, nil
	}
	return nil, fmt.Errorf("Source found but not included in group %s", groupName)
}

// Destinations is the resolver for the destinations field.
func (r *computePlatformResolver) Destinations(ctx context.Context, obj *model.ComputePlatform, groupName string) ([]*model.Destination, error) {
	ns := env.GetCurrentNamespace()

	dests, err := kube.DefaultClient.OdigosClient.Destinations(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var destinations []*model.Destination
	for _, dest := range dests.Items {
		// If matches groupName, return the destination
		if groupName == "" || services.ArrayContains(dest.Spec.SourceSelector.Groups, groupName) {
			secretFields, err := services.GetDestinationSecretFields(ctx, ns, &dest)
			if err != nil {
				return nil, err
			}

			// Convert the k8s destination format to the expected endpoint format
			endpointDest := services.K8sDestinationToEndpointFormat(dest, secretFields)
			destinations = append(destinations, &endpointDest)
		}
	}

	return destinations, nil
}

// Actions is the resolver for the actions field.
func (r *computePlatformResolver) Actions(ctx context.Context, obj *model.ComputePlatform) ([]*model.PipelineAction, error) {
	var response []*model.PipelineAction
	ns := env.GetCurrentNamespace()

	// K8sAttributes actions
	kaActions, err := kube.DefaultClient.ActionsClient.K8sAttributesResolvers(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range kaActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.PipelineAction{
			ID:         action.Name,
			Type:       action.Kind,
			Spec:       string(specStr),
			Conditions: convertConditions(action.Status.Conditions),
		})
	}

	// AddClusterInfos actions
	icaActions, err := kube.DefaultClient.ActionsClient.AddClusterInfos(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range icaActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.PipelineAction{
			ID:         action.Name,
			Type:       action.Kind,
			Spec:       string(specStr),
			Conditions: convertConditions(action.Status.Conditions),
		})
	}

	// DeleteAttributes actions
	daActions, err := kube.DefaultClient.ActionsClient.DeleteAttributes(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range daActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.PipelineAction{
			ID:         action.Name,
			Type:       action.Kind,
			Spec:       string(specStr),
			Conditions: convertConditions(action.Status.Conditions),
		})
	}

	// RenameAttributes actions
	raActions, err := kube.DefaultClient.ActionsClient.RenameAttributes(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range raActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.PipelineAction{
			ID:         action.Name,
			Type:       action.Kind,
			Spec:       string(specStr),
			Conditions: convertConditions(action.Status.Conditions),
		})
	}

	// ErrorSamplers actions
	esActions, err := kube.DefaultClient.ActionsClient.ErrorSamplers(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range esActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.PipelineAction{
			ID:         action.Name,
			Type:       action.Kind,
			Spec:       string(specStr),
			Conditions: convertConditions(action.Status.Conditions),
		})
	}

	// LatencySamplers actions
	lsActions, err := kube.DefaultClient.ActionsClient.LatencySamplers(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range lsActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.PipelineAction{
			ID:         action.Name,
			Type:       action.Kind,
			Spec:       string(specStr),
			Conditions: convertConditions(action.Status.Conditions),
		})
	}

	// ProbabilisticSamplers actions
	psActions, err := kube.DefaultClient.ActionsClient.ProbabilisticSamplers(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range psActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.PipelineAction{
			ID:         action.Name,
			Type:       action.Kind,
			Spec:       string(specStr),
			Conditions: convertConditions(action.Status.Conditions),
		})
	}

	// PiiMaskings actions
	piActions, err := kube.DefaultClient.ActionsClient.PiiMaskings(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range piActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.PipelineAction{
			ID:         action.Name,
			Type:       action.Kind,
			Spec:       string(specStr),
			Conditions: convertConditions(action.Status.Conditions),
		})
	}

	return response, nil
}

// InstrumentationRules is the resolver for the instrumentationRules field.
func (r *computePlatformResolver) InstrumentationRules(ctx context.Context, obj *model.ComputePlatform) ([]*model.InstrumentationRule, error) {
	return services.ListInstrumentationRules(ctx)
}

// Sources is the resolver for the sources field.
func (r *k8sActualNamespaceResolver) Sources(ctx context.Context, obj *model.K8sActualNamespace) ([]*model.K8sActualSource, error) {
	ns := obj.Name

	workloads, err := services.GetWorkloadsInNamespace(ctx, ns)
	if err != nil {
		return nil, err
	}

	sourceList, err := kube.DefaultClient.OdigosClient.Sources(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaceSource *v1alpha1.Source
	sourceObjects := make(map[string]*v1alpha1.Source)
	for _, source := range sourceList.Items {
		if services.WorkloadKind(source.Spec.Workload.Kind) == services.WorkloadKindNamespace {
			namespaceSource = &source
		} else {
			key := fmt.Sprintf("%s/%s/%s", source.Spec.Workload.Namespace, source.Spec.Workload.Name, source.Spec.Workload.Kind)
			sourceObjects[key] = &source
		}
	}

	// Convert workloads to []*model.K8sActualSource
	sources := make([]*model.K8sActualSource, len(workloads))
	for i, workload := range workloads {
		workloadSource := sourceObjects[fmt.Sprintf("%s/%s/%s", workload.Namespace, workload.Name, workload.Kind)]

		namespaceInstrumented := namespaceSource != nil && !namespaceSource.Spec.DisableInstrumentation
		sourceInstrumented := workloadSource != nil && !workloadSource.Spec.DisableInstrumentation
		isInstrumented := (namespaceInstrumented && (sourceInstrumented || workloadSource == nil)) || (!namespaceInstrumented && sourceInstrumented)

		sources[i] = &workload
		sources[i].Selected = &isInstrumented
	}

	return sources, nil
}

// UpdateAPIToken is the resolver for the updateApiToken field.
func (r *mutationResolver) UpdateAPIToken(ctx context.Context, token string) (bool, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return false, services.ErrorIsReadonly
	}

	ns := env.GetCurrentNamespace()
	err := pro.UpdateOdigosToken(ctx, kube.DefaultClient, ns, token)
	return err == nil, nil
}

// PersistK8sNamespace is the resolver for the persistK8sNamespace field.
func (r *mutationResolver) PersistK8sNamespace(ctx context.Context, namespace model.PersistNamespaceItemInput) (bool, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return false, services.ErrorIsReadonly
	}

	persistObjects := []model.PersistNamespaceSourceInput{}
	persistObjects = append(persistObjects, model.PersistNamespaceSourceInput{
		Name:     namespace.Name,
		Kind:     model.K8sResourceKind(services.WorkloadKindNamespace),
		Selected: namespace.FutureSelected,
	})

	err := services.SyncWorkloadsInNamespace(ctx, namespace.Name, persistObjects)
	if err != nil {
		return false, fmt.Errorf("failed to sync workloads: %v", err)
	}

	return true, nil
}

// PersistK8sSources is the resolver for the persistK8sSources field.
func (r *mutationResolver) PersistK8sSources(ctx context.Context, namespace string, sources []*model.PersistNamespaceSourceInput) (bool, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return false, services.ErrorIsReadonly
	}

	var persistObjects []model.PersistNamespaceSourceInput
	for _, source := range sources {
		persistObjects = append(persistObjects, model.PersistNamespaceSourceInput{
			Name:     source.Name,
			Kind:     source.Kind,
			Selected: source.Selected,
		})
	}

	err := services.SyncWorkloadsInNamespace(ctx, namespace, persistObjects)
	if err != nil {
		return false, fmt.Errorf("failed to sync workloads: %v", err)
	}

	return true, nil
}

// UpdateK8sActualSource is the resolver for the updateK8sActualSource field.
func (r *mutationResolver) UpdateK8sActualSource(ctx context.Context, sourceID model.K8sSourceID, patchSourceRequest model.PatchSourceRequestInput) (bool, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return false, services.ErrorIsReadonly
	}

	nsName := sourceID.Namespace
	workloadName := sourceID.Name
	workloadKind := services.WorkloadKind(sourceID.Kind)
	otelServiceName := patchSourceRequest.OtelServiceName

	source, err := services.GetSourceCRD(ctx, nsName, workloadName, workloadKind)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			// unexpected error occurred while trying to get the source
			return false, err
		}

		source, err = services.CreateSourceCRD(ctx, nsName, workloadName, workloadKind)
		if err != nil {
			// unexpected error occurred while trying to create the source
			return false, err
		}
	}

	_, err = services.UpdateSourceCRDSpec(ctx, nsName, source.Name, "otelServiceName", fmt.Sprintf("\"%s\"", otelServiceName))
	if err != nil {
		// unexpected error occurred while trying to update the source
		return false, err
	}

	return true, nil
}

// CreateNewDestination is the resolver for the createNewDestination field.
func (r *mutationResolver) CreateNewDestination(ctx context.Context, destination model.DestinationInput) (*model.Destination, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return nil, services.ErrorIsReadonly
	}

	ns := env.GetCurrentNamespace()
	destType := common.DestinationType(destination.Type)
	destName := destination.Name

	destTypeConfig, err := services.GetDestinationTypeConfig(destType)
	if err != nil {
		return nil, fmt.Errorf("destination type %s not found", destType)
	}

	// Convert fields to map[string]string
	fieldsMap := make(map[string]string)
	for _, field := range destination.Fields {
		fieldsMap[field.Key] = field.Value
	}

	errors := services.VerifyDestinationDataScheme(destType, destTypeConfig, fieldsMap)
	if len(errors) > 0 {
		return nil, fmt.Errorf("invalid destination data scheme: %v", errors)
	}

	dataField, secretFields := services.TransformFieldsToDataAndSecrets(destTypeConfig, fieldsMap)
	generateNamePrefix := "odigos.io.dest." + string(destType) + "-"

	k8sDestination := v1alpha1.Destination{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: generateNamePrefix,
		},
		Spec: v1alpha1.DestinationSpec{
			Type:            destType,
			DestinationName: destName,
			Data:            dataField,
			Signals:         services.ExportedSignalsObjectToSlice(destination.ExportedSignals),
		},
	}

	createSecret := len(secretFields) > 0
	if createSecret {
		secretRef, err := services.CreateDestinationSecret(ctx, destType, secretFields, ns)
		if err != nil {
			return nil, err
		}
		k8sDestination.Spec.SecretRef = secretRef
	}

	dest, err := kube.DefaultClient.OdigosClient.Destinations(ns).Create(ctx, &k8sDestination, metav1.CreateOptions{})
	if err != nil {
		if createSecret {
			kube.DefaultClient.CoreV1().Secrets(ns).Delete(ctx, destName, metav1.DeleteOptions{})
		}
		return nil, err
	}

	if dest.Spec.SecretRef != nil {
		err = services.AddDestinationOwnerReferenceToSecret(ctx, ns, dest)
		if err != nil {
			return nil, err
		}
	}

	secretFieldsMap, err := services.GetDestinationSecretFields(ctx, ns, dest)
	if err != nil {
		return nil, err
	}

	endpointDest := services.K8sDestinationToEndpointFormat(*dest, secretFieldsMap)
	return &endpointDest, nil
}

// UpdateDestination is the resolver for the updateDestination field.
func (r *mutationResolver) UpdateDestination(ctx context.Context, id string, destination model.DestinationInput) (*model.Destination, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return nil, services.ErrorIsReadonly
	}

	ns := env.GetCurrentNamespace()
	destType := common.DestinationType(destination.Type)
	destName := destination.Name

	// Get the destination type configuration
	destTypeConfig, err := services.GetDestinationTypeConfig(destType)
	if err != nil {
		return nil, fmt.Errorf("destination type %s not found: %v", destType, err)
	}

	// Convert fields from input to map[string]string
	fields := make(map[string]string)
	for _, field := range destination.Fields {
		fields[field.Key] = field.Value
	}

	// Validate the destination data schema
	validationErrors := services.VerifyDestinationDataScheme(destType, destTypeConfig, fields)
	if len(validationErrors) > 0 {
		var errMsg string
		for _, e := range validationErrors {
			errMsg += e.Error() + "; "
		}
		return nil, fmt.Errorf("validation errors: %s", errMsg)
	}

	// Separate data fields and secret fields
	dataFields, secretFields := services.TransformFieldsToDataAndSecrets(destTypeConfig, fields)

	// Retrieve the existing destination
	dest, err := kube.DefaultClient.OdigosClient.Destinations(ns).Get(ctx, id, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get destination: %v", err)
	}

	// Handle secrets
	destUpdateHasSecrets := len(secretFields) > 0
	destCurrentlyHasSecrets := dest.Spec.SecretRef != nil

	if !destUpdateHasSecrets && destCurrentlyHasSecrets {
		// Delete the secret if it's not needed anymore
		err := kube.DefaultClient.CoreV1().Secrets(ns).Delete(ctx, dest.Spec.SecretRef.Name, metav1.DeleteOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to delete secret: %v", err)
		}
		dest.Spec.SecretRef = nil
	} else if destUpdateHasSecrets && !destCurrentlyHasSecrets {
		// Create the secret if it was added in this update
		secretRef, err := services.CreateDestinationSecret(ctx, destType, secretFields, ns)
		if err != nil {
			return nil, fmt.Errorf("failed to create secret: %v", err)
		}
		dest.Spec.SecretRef = secretRef
		// Add owner reference to the secret
		err = services.AddDestinationOwnerReferenceToSecret(ctx, ns, dest)
		if err != nil {
			return nil, fmt.Errorf("failed to add owner reference to secret: %v", err)
		}
	} else if destUpdateHasSecrets && destCurrentlyHasSecrets {
		// Update the secret in case it is modified
		secret, err := kube.DefaultClient.CoreV1().Secrets(ns).Get(ctx, dest.Spec.SecretRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get secret: %v", err)
		}
		origSecret := secret.DeepCopy()

		secret.StringData = secretFields
		_, err = kube.DefaultClient.CoreV1().Secrets(ns).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			// Rollback secret if needed
			_, rollbackErr := kube.DefaultClient.CoreV1().Secrets(ns).Update(ctx, origSecret, metav1.UpdateOptions{})
			if rollbackErr != nil {
				fmt.Printf("Failed to rollback secret: %v\n", rollbackErr)
			}
			return nil, fmt.Errorf("failed to update secret: %v", err)
		}
	}

	// Update the destination specification
	dest.Spec.Type = destType
	dest.Spec.DestinationName = destName
	dest.Spec.Data = dataFields
	dest.Spec.Signals = services.ExportedSignalsObjectToSlice(destination.ExportedSignals)

	// Update the destination in Kubernetes
	updatedDest, err := kube.DefaultClient.OdigosClient.Destinations(ns).Update(ctx, dest, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update destination: %v", err)
	}

	// Get the secret fields for the updated destination
	secretFields, err = services.GetDestinationSecretFields(ctx, ns, updatedDest)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret fields: %v", err)
	}

	// Convert the updated destination to the GraphQL model
	resp := services.K8sDestinationToEndpointFormat(*updatedDest, secretFields)

	return &resp, nil
}

// DeleteDestination is the resolver for the deleteDestination field.
func (r *mutationResolver) DeleteDestination(ctx context.Context, id string) (bool, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return false, services.ErrorIsReadonly
	}

	ns := env.GetCurrentNamespace()
	err := kube.DefaultClient.OdigosClient.Destinations(ns).Delete(ctx, id, metav1.DeleteOptions{})

	if err != nil {
		return false, fmt.Errorf("failed to delete destination: %w", err)
	}

	return true, nil
}

// TestConnectionForDestination is the resolver for the testConnectionForDestination field.
func (r *mutationResolver) TestConnectionForDestination(ctx context.Context, destination model.DestinationInput) (*model.TestConnectionResponse, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return nil, services.ErrorIsReadonly
	}

	destType := common.DestinationType(destination.Type)

	destConfig, err := services.GetDestinationTypeConfig(destType)
	if err != nil {
		return nil, err
	}

	if !destConfig.Spec.TestConnectionSupported {
		return nil, fmt.Errorf("destination type %s does not support test connection", destination.Type)
	}

	configurer, err := testconnection.ConvertDestinationToConfigurer(destination)
	if err != nil {
		return nil, err
	}

	res := testconnection.TestConnection(ctx, configurer)
	if !res.Succeeded {
		return &model.TestConnectionResponse{
			Succeeded:       false,
			StatusCode:      res.StatusCode,
			DestinationType: (*string)(&res.DestinationType),
			Message:         &res.Message,
			Reason:          (*string)(&res.Reason),
		}, nil
	}

	return &model.TestConnectionResponse{
		Succeeded:       true,
		StatusCode:      200,
		DestinationType: (*string)(&res.DestinationType),
	}, nil
}

// CreateAction is the resolver for the createAction field.
func (r *mutationResolver) CreateAction(ctx context.Context, action model.ActionInput) (model.Action, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return nil, services.ErrorIsReadonly
	}

	switch action.Type {
	case actionservices.ActionTypeK8sAttributes:
		return actionservices.CreateK8sAttributes(ctx, action)
	case actionservices.ActionTypeAddClusterInfo:
		return actionservices.CreateAddClusterInfo(ctx, action)
	case actionservices.ActionTypeDeleteAttribute:
		return actionservices.CreateDeleteAttribute(ctx, action)
	case actionservices.ActionTypePiiMasking:
		return actionservices.CreatePiiMasking(ctx, action)
	case actionservices.ActionTypeErrorSampler:
		return actionservices.CreateErrorSampler(ctx, action)
	case actionservices.ActionTypeLatencySampler:
		return actionservices.CreateLatencySampler(ctx, action)
	case actionservices.ActionTypeProbabilisticSampler:
		return actionservices.CreateProbabilisticSampler(ctx, action)
	case actionservices.ActionTypeRenameAttribute:
		return actionservices.CreateRenameAttribute(ctx, action)
	default:
		return nil, fmt.Errorf("unsupported action type: %s", action.Type)
	}
}

// UpdateAction is the resolver for the updateAction field.
func (r *mutationResolver) UpdateAction(ctx context.Context, id string, action model.ActionInput) (model.Action, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return nil, services.ErrorIsReadonly
	}

	switch action.Type {
	case actionservices.ActionTypeK8sAttributes:
		return actionservices.UpdateK8sAttributes(ctx, id, action)
	case actionservices.ActionTypeAddClusterInfo:
		return actionservices.UpdateAddClusterInfo(ctx, id, action)
	case actionservices.ActionTypeDeleteAttribute:
		return actionservices.UpdateDeleteAttribute(ctx, id, action)
	case actionservices.ActionTypePiiMasking:
		return actionservices.UpdatePiiMasking(ctx, id, action)
	case actionservices.ActionTypeErrorSampler:
		return actionservices.UpdateErrorSampler(ctx, id, action)
	case actionservices.ActionTypeLatencySampler:
		return actionservices.UpdateLatencySampler(ctx, id, action)
	case actionservices.ActionTypeProbabilisticSampler:
		return actionservices.UpdateProbabilisticSampler(ctx, id, action)
	case actionservices.ActionTypeRenameAttribute:
		return actionservices.UpdateRenameAttribute(ctx, id, action)
	default:
		return nil, fmt.Errorf("unsupported action type: %s", action.Type)
	}
}

// DeleteAction is the resolver for the deleteAction field.
func (r *mutationResolver) DeleteAction(ctx context.Context, id string, actionType string) (bool, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return false, services.ErrorIsReadonly
	}

	switch actionType {
	case actionservices.ActionTypeK8sAttributes:
		err := actionservices.DeleteK8sAttributes(ctx, id)
		if err != nil {
			return false, fmt.Errorf("failed to delete K8sAttributes: %v", err)
		}

	case actionservices.ActionTypeAddClusterInfo:
		err := actionservices.DeleteAddClusterInfo(ctx, id)
		if err != nil {
			return false, fmt.Errorf("failed to delete AddClusterInfo: %v", err)
		}

	case actionservices.ActionTypeDeleteAttribute:
		err := actionservices.DeleteDeleteAttribute(ctx, id)
		if err != nil {
			return false, fmt.Errorf("failed to delete DeleteAttribute: %v", err)
		}
	case actionservices.ActionTypePiiMasking:
		err := actionservices.DeletePiiMasking(ctx, id)
		if err != nil {
			return false, fmt.Errorf("failed to delete PiiMasking: %v", err)
		}
	case actionservices.ActionTypeErrorSampler:
		err := actionservices.DeleteErrorSampler(ctx, id)
		if err != nil {
			return false, fmt.Errorf("failed to delete ErrorSampler: %v", err)
		}
	case actionservices.ActionTypeLatencySampler:
		err := actionservices.DeleteLatencySampler(ctx, id)
		if err != nil {
			return false, fmt.Errorf("failed to delete LatencySampler: %v", err)
		}
	case actionservices.ActionTypeProbabilisticSampler:
		err := actionservices.DeleteProbabilisticSampler(ctx, id)
		if err != nil {
			return false, fmt.Errorf("failed to delete ProbabilisticSampler: %v", err)
		}
	case actionservices.ActionTypeRenameAttribute:
		err := actionservices.DeleteRenameAttribute(ctx, id)
		if err != nil {
			return false, fmt.Errorf("failed to delete RenameAttribute: %v", err)
		}
	default:
		return false, fmt.Errorf("unsupported action type: %s", actionType)
	}

	// Return true if the deletion was successful
	return true, nil
}

// CreateInstrumentationRule is the resolver for the createInstrumentationRule field.
func (r *mutationResolver) CreateInstrumentationRule(ctx context.Context, instrumentationRule model.InstrumentationRuleInput) (*model.InstrumentationRule, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return nil, services.ErrorIsReadonly
	}

	return services.CreateInstrumentationRule(ctx, instrumentationRule)
}

// UpdateInstrumentationRule is the resolver for the updateInstrumentationRule field.
func (r *mutationResolver) UpdateInstrumentationRule(ctx context.Context, ruleID string, instrumentationRule model.InstrumentationRuleInput) (*model.InstrumentationRule, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return nil, services.ErrorIsReadonly
	}

	return services.UpdateInstrumentationRule(ctx, ruleID, instrumentationRule)
}

// DeleteInstrumentationRule is the resolver for the deleteInstrumentationRule field.
func (r *mutationResolver) DeleteInstrumentationRule(ctx context.Context, ruleID string) (bool, error) {
	isReadonly := services.IsReadonlyMode(ctx)
	if isReadonly {
		return false, services.ErrorIsReadonly
	}

	_, err := services.DeleteInstrumentationRule(ctx, ruleID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ComputePlatform is the resolver for the computePlatform field.
func (r *queryResolver) ComputePlatform(ctx context.Context) (*model.ComputePlatform, error) {
	return &model.ComputePlatform{
		ComputePlatformType: model.ComputePlatformTypeK8s,
	}, nil
}

// Config is the resolver for the config field.
func (r *queryResolver) Config(ctx context.Context) (*model.GetConfigResponse, error) {
	config := services.GetConfig(ctx)

	return &config, nil
}

// DestinationCategories is the resolver for the destinationCategories field.
func (r *queryResolver) DestinationCategories(ctx context.Context) (*model.GetDestinationCategories, error) {
	destTypes := services.GetDestinationCategories()

	return &destTypes, nil
}

// PotentialDestinations is the resolver for the potentialDestinations field.
func (r *queryResolver) PotentialDestinations(ctx context.Context) ([]*model.DestinationDetails, error) {
	potentialDestinations := services.PotentialDestinations(ctx)
	if potentialDestinations == nil {
		return nil, fmt.Errorf("failed to fetch potential destinations")
	}

	// Convert []destination_recognition.DestinationDetails to []*DestinationDetails
	var result []*model.DestinationDetails
	for _, dest := range potentialDestinations {

		fieldsString, err := json.Marshal(dest.Fields)
		if err != nil {
			return nil, fmt.Errorf("error marshalling fields: %v", err)
		}

		result = append(result, &model.DestinationDetails{
			Type:   string(dest.Type),
			Fields: string(fieldsString),
		})
	}

	return result, nil
}

// GetOverviewMetrics is the resolver for the getOverviewMetrics field.
func (r *queryResolver) GetOverviewMetrics(ctx context.Context) (*model.OverviewMetricsResponse, error) {
	if r.MetricsConsumer == nil {
		return nil, fmt.Errorf("metrics consumer not initialized")
	}

	sourcesMetrics := r.MetricsConsumer.GetSourcesMetrics()
	destinationsMetrics := r.MetricsConsumer.GetDestinationsMetrics()

	var sourcesResp []*model.SingleSourceMetricsResponse
	for sID, metric := range sourcesMetrics {
		sourcesResp = append(sourcesResp, &model.SingleSourceMetricsResponse{
			Namespace:     sID.Namespace,
			Kind:          string(sID.Kind),
			Name:          sID.Name,
			TotalDataSent: int(metric.TotalDataSent()),
			Throughput:    int(metric.TotalThroughput()),
		})
	}

	var destinationsResp []*model.SingleDestinationMetricsResponse
	for destId, metric := range destinationsMetrics {
		destinationsResp = append(destinationsResp, &model.SingleDestinationMetricsResponse{
			ID:            destId,
			TotalDataSent: int(metric.TotalDataSent()),
			Throughput:    int(metric.TotalThroughput()),
		})
	}

	return &model.OverviewMetricsResponse{
		Sources:      sourcesResp,
		Destinations: destinationsResp,
	}, nil
}

// DescribeOdigos is the resolver for the describeOdigos field.
func (r *queryResolver) DescribeOdigos(ctx context.Context) (*model.OdigosAnalyze, error) {
	return odigos_describe.GetOdigosDescription(ctx)
}

// DescribeSource is the resolver for the describeSource field.
func (r *queryResolver) DescribeSource(ctx context.Context, namespace string, kind string, name string) (*model.SourceAnalyze, error) {
	return source_describe.GetSourceDescription(ctx, namespace, kind, name)
}

// InstrumentationInstancesHealth is the resolver for the instrumentationInstancesHealth field.
func (r *queryResolver) InstrumentationInstancesHealth(ctx context.Context) ([]*model.InstrumentationInstanceHealth, error) {
	return services.GetInstrumentationInstancesHealthConditions(ctx)
}

// GroupNames is the resolver for the groupNames field.
func (r *queryResolver) GroupNames(ctx context.Context) ([]string, error) {
	ns := env.GetCurrentNamespace()

	dests, err := kube.DefaultClient.OdigosClient.Destinations(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var groupNames []string

	// Collect group names without duplicates
	seen := make(map[string]bool)
	for _, dest := range dests.Items {
		for _, group := range dest.Spec.SourceSelector.Groups {
			if _, exists := seen[group]; !exists {
				seen[group] = true
				groupNames = append(groupNames, group)
			}
		}
	}

	return groupNames, nil
}

// ComputePlatform returns ComputePlatformResolver implementation.
func (r *Resolver) ComputePlatform() ComputePlatformResolver { return &computePlatformResolver{r} }

// K8sActualNamespace returns K8sActualNamespaceResolver implementation.
func (r *Resolver) K8sActualNamespace() K8sActualNamespaceResolver {
	return &k8sActualNamespaceResolver{r}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type computePlatformResolver struct{ *Resolver }
type k8sActualNamespaceResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
