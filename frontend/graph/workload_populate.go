package graph

import (
	"context"
	"sync"

	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/frontend/graph/loaders"
	"github.com/odigos-io/odigos/frontend/graph/model"
	"github.com/odigos-io/odigos/frontend/graph/status"
	"github.com/odigos-io/odigos/frontend/services"
	sourceutils "github.com/odigos-io/odigos/k8sutils/pkg/source"
)

// Serializes heavy workload queries so concurrent requests don't double memory usage.
// A single GET_WORKLOADS at 10K workloads uses most of the 512Mi budget; a second
// concurrent query would OOM. The mutex makes the second query wait until the first
// completes and its memory is reclaimed by GC.
var heavyWorkloadQueryMu sync.Mutex

// populateWorkloadFields pre-computes all resolver fields for a workload
// sequentially. This is called from the Workloads batch resolver to avoid
// gqlgen spawning a goroutine per field per workload.
// Errors are logged but don't fail the entire batch — individual workloads
// get nil/zero values for fields that fail to resolve.
func (r *queryResolver) populateWorkloadFields(ctx context.Context, l *loaders.Loaders, w *model.K8sWorkload) {
	id := *w.ID

	ic, _ := l.GetInstrumentationConfig(ctx, id)

	// serviceName
	if ic != nil {
		w.ServiceName = &ic.Spec.ServiceName
	}

	// rollbackOccurred
	if ic != nil {
		w.RollbackOccurred = ic.Status.RollbackOccurred
	}

	// markedForInstrumentation
	if sources, err := l.GetSources(ctx, id); err == nil {
		enabled, reason, err := sourceutils.IsObjectInstrumentedBySource(ctx, sources, nil)
		if err == nil {
			var marked *bool
			if enabled {
				marked = &enabled
			} else if reason.Reason == string("WorkloadSourceDisabled") {
				marked = &enabled
			}
			w.MarkedForInstrumentation = &model.K8sWorkloadMarkedForInstrumentation{
				MarkedForInstrumentation: marked,
				DecisionEnum:             string(reason.Reason),
				Message:                  reason.Message,
			}
		}
	}

	// dataStreamNames
	if sources, err := l.GetSources(ctx, id); err == nil {
		ptrNames := services.ExtractDataStreamsFromSource(sources.Workload, sources.Namespace)
		names := make([]string, len(ptrNames))
		for i, p := range ptrNames {
			names[i] = *p
		}
		w.DataStreamNames = names
	}

	// numberOfInstances
	if manifest, err := l.GetWorkloadManifest(ctx, id); err == nil && manifest != nil {
		count := int(manifest.AvailableReplicas)
		w.NumberOfInstances = &count
	}

	// runtimeInfo
	if ic != nil {
		completed := len(ic.Status.RuntimeDetailsByContainer) > 0
		uniqueLanguages := make(map[common.ProgrammingLanguage]struct{})
		containers := make([]*model.K8sWorkloadRuntimeInfoContainer, len(ic.Status.RuntimeDetailsByContainer))
		for i, container := range ic.Status.RuntimeDetailsByContainer {
			containers[i] = runtimeDetailsContainersToModel(&container)
			_, ignored := l.GetIgnoredContainers()[container.ContainerName]
			if container.Language != common.UnknownProgrammingLanguage && !ignored {
				uniqueLanguages[container.Language] = struct{}{}
			}
		}
		var detectedLanguages []model.ProgrammingLanguage
		if completed {
			detectedLanguages = make([]model.ProgrammingLanguage, 0, len(uniqueLanguages))
			for language := range uniqueLanguages {
				detectedLanguages = append(detectedLanguages, model.ProgrammingLanguage(language))
			}
		}
		w.RuntimeInfo = &model.K8sWorkloadRuntimeInfo{
			Completed:         completed,
			CompletedStatus:   status.CalculateRuntimeInspectionStatus(ic),
			DetectedLanguages: detectedLanguages,
			Containers:        containers,
		}
	}

	// containers
	if ic != nil {
		containerByName := make(map[string]*model.K8sWorkloadContainer)
		for i := range ic.Spec.Containers {
			container := &ic.Spec.Containers[i]
			if _, ok := containerByName[container.ContainerName]; !ok {
				containerByName[container.ContainerName] = &model.K8sWorkloadContainer{
					ContainerName: container.ContainerName,
				}
			}
			containerByName[container.ContainerName].AgentEnabled = agentEnabledContainersToModel(container)
			containerByName[container.ContainerName].AgentConfig = containerAgentConfigToAgentConfigModel(container)
		}
		for _, container := range ic.Status.RuntimeDetailsByContainer {
			if _, ok := containerByName[container.ContainerName]; !ok {
				containerByName[container.ContainerName] = &model.K8sWorkloadContainer{
					ContainerName: container.ContainerName,
				}
			}
			containerByName[container.ContainerName].RuntimeInfo = runtimeDetailsContainersToModel(&container)
		}
		for _, container := range ic.Spec.ContainersOverrides {
			if _, ok := containerByName[container.ContainerName]; !ok {
				containerByName[container.ContainerName] = &model.K8sWorkloadContainer{
					ContainerName: container.ContainerName,
				}
			}
			overrides := &model.K8sWorkloadContainerOverrides{
				ContainerName: container.ContainerName,
			}
			if container.RuntimeInfo != nil {
				overrides.RuntimeInfo = runtimeDetailsContainersToModel(container.RuntimeInfo)
			}
			containerByName[container.ContainerName].Overrides = overrides
		}
		w.Containers = make([]*model.K8sWorkloadContainer, 0, len(containerByName))
		for _, container := range containerByName {
			w.Containers = append(w.Containers, container)
		}
	}

	// Pod-dependent fields (conditions, workloadOdigosHealthStatus, podsAgentInjectionStatus) are NOT pre-computed here.
	// They require loading ALL pods (12K CachedPod structs) and computing status for each, which adds ~15MB+ of allocations that push the pod past its 512Mi limit when combined with gqlgen's serialization buffer.
	// These 3 fields are left for gqlgen's lazy resolver path (30K goroutines = ~75MB, which is affordable compared to the 170K goroutines from pre-Fix 10).
}
