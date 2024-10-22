package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/common/consts"
	"github.com/odigos-io/odigos/frontend/endpoints"
	"github.com/odigos-io/odigos/frontend/graph/model"
	"github.com/odigos-io/odigos/frontend/kube"
	"github.com/odigos-io/odigos/frontend/services"
	actionservices "github.com/odigos-io/odigos/frontend/services/actions"
	testconnection "github.com/odigos-io/odigos/frontend/services/test_connection"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// K8sActualNamespace is the resolver for the k8sActualNamespace field.
func (r *computePlatformResolver) K8sActualNamespace(ctx context.Context, obj *model.ComputePlatform, name string) (*model.K8sActualNamespace, error) {
	namespaceActualSources, err := services.GetWorkloadsInNamespace(ctx, name, nil)
	if err != nil {
		return nil, err
	}

	// Convert namespaceActualSources to []*model.K8sActualSource
	namespaceActualSourcesPointers := make([]*model.K8sActualSource, len(namespaceActualSources))
	for i, source := range namespaceActualSources {
		namespaceActualSourcesPointers[i] = &source
	}

	namespace, err := kube.DefaultClient.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	nsInstrumented := workload.GetInstrumentationLabelValue(namespace.GetLabels())

	return &model.K8sActualNamespace{
		Name:                        name,
		InstrumentationLabelEnabled: nsInstrumented,
		K8sActualSources:            namespaceActualSourcesPointers,
	}, nil
}

// K8sActualNamespaces is the resolver for the k8sActualNamespaces field.
func (r *computePlatformResolver) K8sActualNamespaces(ctx context.Context, obj *model.ComputePlatform) ([]*model.K8sActualNamespace, error) {
	namespacesResponse := services.GetK8SNamespaces(ctx)

	K8sActualNamespaces := make([]*model.K8sActualNamespace, len(namespacesResponse.Namespaces))
	for i, namespace := range namespacesResponse.Namespaces {

		namespace, err := kube.DefaultClient.CoreV1().Namespaces().Get(ctx, namespace.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		nsInstrumented := workload.GetInstrumentationLabelValue(namespace.GetLabels())

		K8sActualNamespaces[i] = &model.K8sActualNamespace{
			Name:                        namespace.Name,
			InstrumentationLabelEnabled: nsInstrumented,
		}
	}

	return K8sActualNamespaces, nil
}

// K8sActualSource is the resolver for the k8sActualSource field.
func (r *computePlatformResolver) K8sActualSource(ctx context.Context, obj *model.ComputePlatform, name *string, namespace *string, kind *string) (*model.K8sActualSource, error) {
	source, err := services.GetActualSource(ctx, *namespace, *kind, *name)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, nil
	}
	k8sActualSource := k8sSourceToGql(source)

	return k8sActualSource, nil
}

// K8sActualSources is the resolver for the k8sActualSources field.
func (r *computePlatformResolver) K8sActualSources(ctx context.Context, obj *model.ComputePlatform) ([]*model.K8sActualSource, error) {
	instrumentedApplications, err := kube.DefaultClient.OdigosClient.InstrumentedApplications("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Initialize an empty list of K8sActualSource
	var actualSources []*model.K8sActualSource

	// Convert each instrumented application to the K8sActualSource type
	for _, app := range instrumentedApplications.Items {
		actualSource := instrumentedApplicationToActualSource(app)

		owner, _ := services.GetWorkload(ctx, actualSource.Namespace, string(actualSource.Kind), actualSource.Name)
		if owner == nil {

			continue
		}
		ownerAnnotations := owner.GetAnnotations()
		var reportedName string
		if ownerAnnotations != nil {
			reportedName = ownerAnnotations[consts.OdigosReportedNameAnnotation]
		}
		actualSource.ReportedName = &reportedName
		actualSources = append(actualSources, actualSource)
	}

	return actualSources, nil
}

// Destinations is the resolver for the destinations field.
func (r *computePlatformResolver) Destinations(ctx context.Context, obj *model.ComputePlatform) ([]*model.Destination, error) {
	odigosns := consts.DefaultOdigosNamespace
	dests, err := kube.DefaultClient.OdigosClient.Destinations(odigosns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var destinations []*model.Destination
	for _, dest := range dests.Items {
		secretFields, err := services.GetDestinationSecretFields(ctx, odigosns, &dest)
		if err != nil {
			return nil, err
		}

		// Convert the k8s destination format to the expected endpoint format
		endpointDest := services.K8sDestinationToEndpointFormat(dest, secretFields)
		destinations = append(destinations, &endpointDest)
	}

	return destinations, nil
}

// Actions is the resolver for the actions field.
func (r *computePlatformResolver) Actions(ctx context.Context, obj *model.ComputePlatform) ([]*model.IcaInstanceResponse, error) {
	var response []*model.IcaInstanceResponse
	odigosns := consts.DefaultOdigosNamespace

	// AddClusterInfos actions
	icaActions, err := kube.DefaultClient.ActionsClient.AddClusterInfos(odigosns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range icaActions.Items {
		specStr, err := json.Marshal(action.Spec) // Convert spec to JSON string
		if err != nil {
			return nil, err
		}
		response = append(response, &model.IcaInstanceResponse{
			ID:   action.Name,
			Type: action.Kind,
			Spec: string(specStr), // Return the JSON string
		})
	}

	// DeleteAttributes actions
	daActions, err := kube.DefaultClient.ActionsClient.DeleteAttributes(odigosns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range daActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.IcaInstanceResponse{
			ID:   action.Name,
			Type: action.Kind,
			Spec: string(specStr),
		})
	}

	// RenameAttributes actions
	raActions, err := kube.DefaultClient.ActionsClient.RenameAttributes(odigosns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range raActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.IcaInstanceResponse{
			ID:   action.Name,
			Type: action.Kind,
			Spec: string(specStr),
		})
	}

	// ErrorSamplers actions
	esActions, err := kube.DefaultClient.ActionsClient.ErrorSamplers(odigosns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range esActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.IcaInstanceResponse{
			ID:   action.Name,
			Type: action.Kind,
			Spec: string(specStr),
		})
	}

	// LatencySamplers actions
	lsActions, err := kube.DefaultClient.ActionsClient.LatencySamplers(odigosns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range lsActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.IcaInstanceResponse{
			ID:   action.Name,
			Type: action.Kind,
			Spec: string(specStr),
		})
	}

	// ProbabilisticSamplers actions
	psActions, err := kube.DefaultClient.ActionsClient.ProbabilisticSamplers(odigosns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range psActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.IcaInstanceResponse{
			ID:   action.Name,
			Type: action.Kind,
			Spec: string(specStr),
		})
	}

	// PiiMaskings actions
	piActions, err := kube.DefaultClient.ActionsClient.PiiMaskings(odigosns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, action := range piActions.Items {
		specStr, err := json.Marshal(action.Spec)
		if err != nil {
			return nil, err
		}
		response = append(response, &model.IcaInstanceResponse{
			ID:   action.Name,
			Type: action.Kind,
			Spec: string(specStr),
		})
	}

	return response, nil
}

// Type is the resolver for the type field.
func (r *destinationResolver) Type(ctx context.Context, obj *model.Destination) (string, error) {
	panic(fmt.Errorf("not implemented: Type - type"))
}

// Conditions is the resolver for the conditions field.
func (r *destinationResolver) Conditions(ctx context.Context, obj *model.Destination) ([]*model.Condition, error) {
	panic(fmt.Errorf("not implemented: Conditions - conditions"))
}

// K8sActualSources is the resolver for the k8sActualSources field.
func (r *k8sActualNamespaceResolver) K8sActualSources(ctx context.Context, obj *model.K8sActualNamespace, instrumentationLabeled *bool) ([]*model.K8sActualSource, error) {
	namespaceActualSources, err := services.GetWorkloadsInNamespace(ctx, obj.Name, instrumentationLabeled)
	if err != nil {
		return nil, err
	}

	// Convert namespaceActualSources to []*model.K8sActualSource
	namespaceActualSourcesPointers := make([]*model.K8sActualSource, len(namespaceActualSources))
	for i, source := range namespaceActualSources {
		namespaceActualSourcesPointers[i] = &source
	}

	return namespaceActualSourcesPointers, nil
}

// CreateNewDestination is the resolver for the createNewDestination field.
func (r *mutationResolver) CreateNewDestination(ctx context.Context, destination model.DestinationInput) (*model.Destination, error) {
	odigosns := consts.DefaultOdigosNamespace

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
		secretRef, err := services.CreateDestinationSecret(ctx, destType, secretFields, odigosns)
		if err != nil {
			return nil, err
		}
		k8sDestination.Spec.SecretRef = secretRef
	}

	dest, err := kube.DefaultClient.OdigosClient.Destinations(odigosns).Create(ctx, &k8sDestination, metav1.CreateOptions{})
	if err != nil {
		if createSecret {
			kube.DefaultClient.CoreV1().Secrets(odigosns).Delete(ctx, destName, metav1.DeleteOptions{})
		}
		return nil, err
	}

	if dest.Spec.SecretRef != nil {
		err = services.AddDestinationOwnerReferenceToSecret(ctx, odigosns, dest)
		if err != nil {
			return nil, err
		}
	}

	secretFieldsMap, err := services.GetDestinationSecretFields(ctx, odigosns, dest)
	if err != nil {
		return nil, err
	}

	endpointDest := services.K8sDestinationToEndpointFormat(*dest, secretFieldsMap)
	return &endpointDest, nil
}

// PersistK8sNamespace is the resolver for the persistK8sNamespace field.
func (r *mutationResolver) PersistK8sNamespace(ctx context.Context, namespace model.PersistNamespaceItemInput) (bool, error) {
	jsonMergePayload := services.GetJsonMergePatchForInstrumentationLabel(namespace.FutureSelected)
	_, err := kube.DefaultClient.CoreV1().Namespaces().Patch(ctx, namespace.Name, types.MergePatchType, jsonMergePayload, metav1.PatchOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to patch namespace: %v", err)
	}

	return true, nil
}

// PersistK8sSources is the resolver for the persistK8sSources field.
func (r *mutationResolver) PersistK8sSources(ctx context.Context, namespace string, sources []*model.PersistNamespaceSourceInput) (bool, error) {
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

// TestConnectionForDestination is the resolver for the testConnectionForDestination field.
func (r *mutationResolver) TestConnectionForDestination(ctx context.Context, destination model.DestinationInput) (*model.TestConnectionResponse, error) {
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

// UpdateK8sActualSource is the resolver for the updateK8sActualSource field.
func (r *mutationResolver) UpdateK8sActualSource(ctx context.Context, sourceID model.K8sSourceID, patchSourceRequest model.PatchSourceRequestInput) (bool, error) {
	ns := sourceID.Namespace
	kind := string(sourceID.Kind)
	name := sourceID.Name

	request := patchSourceRequest

	// Handle ReportedName update
	if request.ReportedName != nil {
		if err := services.UpdateReportedName(ctx, ns, kind, name, *request.ReportedName); err != nil {
			return false, err
		}
	}

	return true, nil
}

// UpdateDestination is the resolver for the updateDestination field.
func (r *mutationResolver) UpdateDestination(ctx context.Context, id string, destination model.DestinationInput) (*model.Destination, error) {
	odigosns := consts.DefaultOdigosNamespace

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
	dest, err := kube.DefaultClient.OdigosClient.Destinations(odigosns).Get(ctx, id, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get destination: %v", err)
	}

	// Handle secrets
	destUpdateHasSecrets := len(secretFields) > 0
	destCurrentlyHasSecrets := dest.Spec.SecretRef != nil

	if !destUpdateHasSecrets && destCurrentlyHasSecrets {
		// Delete the secret if it's not needed anymore
		err := kube.DefaultClient.CoreV1().Secrets(odigosns).Delete(ctx, dest.Spec.SecretRef.Name, metav1.DeleteOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to delete secret: %v", err)
		}
		dest.Spec.SecretRef = nil
	} else if destUpdateHasSecrets && !destCurrentlyHasSecrets {
		// Create the secret if it was added in this update
		secretRef, err := services.CreateDestinationSecret(ctx, destType, secretFields, odigosns)
		if err != nil {
			return nil, fmt.Errorf("failed to create secret: %v", err)
		}
		dest.Spec.SecretRef = secretRef
		// Add owner reference to the secret
		err = services.AddDestinationOwnerReferenceToSecret(ctx, odigosns, dest)
		if err != nil {
			return nil, fmt.Errorf("failed to add owner reference to secret: %v", err)
		}
	} else if destUpdateHasSecrets && destCurrentlyHasSecrets {
		// Update the secret in case it is modified
		secret, err := kube.DefaultClient.CoreV1().Secrets(odigosns).Get(ctx, dest.Spec.SecretRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get secret: %v", err)
		}
		origSecret := secret.DeepCopy()

		secret.StringData = secretFields
		_, err = kube.DefaultClient.CoreV1().Secrets(odigosns).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			// Rollback secret if needed
			_, rollbackErr := kube.DefaultClient.CoreV1().Secrets(odigosns).Update(ctx, origSecret, metav1.UpdateOptions{})
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
	updatedDest, err := kube.DefaultClient.OdigosClient.Destinations(odigosns).Update(ctx, dest, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update destination: %v", err)
	}

	// Get the secret fields for the updated destination
	secretFields, err = services.GetDestinationSecretFields(ctx, odigosns, updatedDest)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret fields: %v", err)
	}

	// Convert the updated destination to the GraphQL model
	resp := services.K8sDestinationToEndpointFormat(*updatedDest, secretFields)

	return &resp, nil
}

// CreateAction is the resolver for the createAction field.
func (r *mutationResolver) CreateAction(ctx context.Context, action model.ActionInput) (model.Action, error) {
	switch action.Type {
	case actionservices.ActionTypeAddClusterInfo:
		return actionservices.CreateAddClusterInfo(ctx, action)

	case actionservices.ActionTypeDeleteAttribute:
		return actionservices.CreateDeleteAttribute(ctx, action)

	default:
		return nil, fmt.Errorf("unsupported action type: %s", action.Type)
	}

}

// UpdateAction is the resolver for the updateAction field.
func (r *mutationResolver) UpdateAction(ctx context.Context, id string, action model.ActionInput) (model.Action, error) {
	switch action.Type {
	case actionservices.ActionTypeAddClusterInfo:
		return actionservices.UpdateAddClusterInfo(ctx, id, action)

	case actionservices.ActionTypeDeleteAttribute:
		return actionservices.UpdateDeleteAttribute(ctx, id, action)

	default:
		return nil, fmt.Errorf("unsupported action type: %s", action.Type)
	}

}

// DeleteAction is the resolver for the deleteAction field.
func (r *mutationResolver) DeleteAction(ctx context.Context, id string, actionType string) (bool, error) {
	switch actionType {
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

	default:
		return false, fmt.Errorf("unsupported action type: %s", actionType)
	}

	// Return true if the deletion was successful
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
	response := endpoints.GetConfig(ctx)

	gqlResponse := &model.GetConfigResponse{
		Installation: model.InstallationStatus(response.Installation),
	}

	return gqlResponse, nil
}

// DestinationTypes is the resolver for the destinationTypes field.
func (r *queryResolver) DestinationTypes(ctx context.Context) (*model.GetDestinationTypesResponse, error) {
	destTypes := services.GetDestinationTypes()

	return &destTypes, nil
}

// DestinationTypeDetails is the resolver for the destinationTypeDetails field.
func (r *queryResolver) DestinationTypeDetails(ctx context.Context, typeArg string) (*model.GetDestinationDetailsResponse, error) {
	destType := common.DestinationType(typeArg)
	destTypeConfig, err := services.GetDestinationTypeConfig(destType)
	if err != nil {
		return nil, fmt.Errorf("destination type %s not found", destType)
	}

	var resp model.GetDestinationDetailsResponse
	for _, field := range destTypeConfig.Spec.Fields {
		componentPropsJSON, err := json.Marshal(field.ComponentProps)
		if err != nil {
			return nil, fmt.Errorf("error marshalling component properties: %v", err)
		}

		resp.Fields = append(resp.Fields, &model.Field{
			Name:                field.Name,
			DisplayName:         field.DisplayName,
			ComponentType:       field.ComponentType,
			ComponentProperties: string(componentPropsJSON),
			InitialValue:        &field.InitialValue,
		})

	}

	return &resp, nil
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

// ComputePlatform returns ComputePlatformResolver implementation.
func (r *Resolver) ComputePlatform() ComputePlatformResolver { return &computePlatformResolver{r} }

// Destination returns DestinationResolver implementation.
func (r *Resolver) Destination() DestinationResolver { return &destinationResolver{r} }

// K8sActualNamespace returns K8sActualNamespaceResolver implementation.
func (r *Resolver) K8sActualNamespace() K8sActualNamespaceResolver {
	return &k8sActualNamespaceResolver{r}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type computePlatformResolver struct{ *Resolver }
type destinationResolver struct{ *Resolver }
type k8sActualNamespaceResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
