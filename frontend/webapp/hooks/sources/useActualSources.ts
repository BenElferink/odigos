import { useState, useCallback } from 'react';
import { usePersistSource } from '../sources';
import { useNamespace } from '../compute-platform';
import { useUpdateSource } from './useUpdateSource';
import { useComputePlatform } from '../compute-platform';
import {
  PatchSourceRequestInput,
  PersistSourcesArray,
  WorkloadId,
} from '@/types';

export function useActualSources() {
  const { data, refetch } = useComputePlatform();
  const { persistSource, error: sourceError } = usePersistSource();
  const { updateSource, error: updateError } = useUpdateSource();

  const { persistNamespace } = useNamespace(undefined);
  const [isPolling, setIsPolling] = useState(false);

  const startPolling = useCallback(async () => {
    setIsPolling(true);
    const maxRetries = 5;
    const retryInterval = 1000; // Poll every second
    let retries = 0;

    while (retries < maxRetries) {
      await new Promise((resolve) => setTimeout(resolve, retryInterval));
      refetch();
      retries++;
    }

    setIsPolling(false);
  }, [refetch]);

  const createSourcesForNamespace = async (
    namespaceName: string,
    sources: PersistSourcesArray[]
  ) => {
    await persistSource(namespaceName, sources);

    startPolling();
    if (sourceError) {
      throw new Error(`Error creating sources for namespace: ${namespaceName}`);
    }
  };

  const deleteSourcesForNamespace = async (
    namespaceName: string,
    sources: PersistSourcesArray[]
  ) => {
    await persistSource(namespaceName, sources);

    startPolling();
    if (sourceError) {
      throw new Error(`Error creating sources for namespace: ${namespaceName}`);
    }
  };

  const updateActualSource = async (
    sourceId: WorkloadId,
    patchRequest: PatchSourceRequestInput
  ) => {
    try {
      await updateSource(sourceId, patchRequest);
      refetch();
    } catch (error) {
      console.error('Error updating source:', error);
      throw error;
    }
  };

  const persistNamespaceItems = async (namespaceItems) => {
    for (const namespace of namespaceItems) {
      await persistNamespace(namespace);
    }
  };

  return {
    sources: data?.computePlatform.k8sActualSources || [],
    deleteSourcesForNamespace,
    createSourcesForNamespace,
    persistNamespaceItems,
    updateActualSource,
    isPolling,
  };
}
