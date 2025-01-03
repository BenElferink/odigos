import { useCallback } from 'react';
import { useMutation } from '@apollo/client';
import { ACTION, getSseTargetFromId } from '@/utils';
import { useAppStore, useNotificationStore } from '@/store';
import { PERSIST_SOURCE, UPDATE_K8S_ACTUAL_SOURCE } from '@/graphql';
import { useComputePlatform, useNamespace } from '../compute-platform';
import { OVERVIEW_ENTITY_TYPES, type WorkloadId, type PatchSourceRequestInput, type K8sActualSource, NOTIFICATION_TYPE } from '@/types';

interface Params {
  onSuccess?: (type: string) => void;
  onError?: (type: string) => void;
}

export const useSourceCRUD = (params?: Params) => {
  const removeNotifications = useNotificationStore((store) => store.removeNotifications);
  const { configuredSources, setConfiguredSources } = useAppStore();

  const { persistNamespace } = useNamespace();
  const { data, refetch } = useComputePlatform();
  const { addNotification } = useNotificationStore();

  const startPolling = useCallback(async () => {
    let retries = 0;
    const maxRetries = 5;
    const retryInterval = 1 * 1000; // time in milliseconds

    while (retries < maxRetries) {
      await new Promise((resolve) => setTimeout(resolve, retryInterval));
      refetch();
      retries++;
    }
  }, [refetch]);

  const notifyUser = (type: NOTIFICATION_TYPE, title: string, message: string, id?: WorkloadId) => {
    addNotification({
      type,
      title,
      message,
      crdType: OVERVIEW_ENTITY_TYPES.SOURCE,
      target: id ? getSseTargetFromId(id, OVERVIEW_ENTITY_TYPES.SOURCE) : undefined,
    });
  };

  const handleError = (title: string, message: string, id?: WorkloadId) => {
    notifyUser(NOTIFICATION_TYPE.ERROR, title, message, id);
    params?.onError?.(title);
  };

  const handleComplete = (title: string, message: string, id?: WorkloadId) => {
    notifyUser(NOTIFICATION_TYPE.SUCCESS, title, message, id);
    startPolling();
    params?.onSuccess?.(title);
  };

  const [createOrDeleteSources, cdState] = useMutation<{ persistK8sSources: boolean }>(PERSIST_SOURCE, {
    onError: (error, req) => {
      const { selected } = req?.variables?.sources?.[0] || {};
      const action = selected ? ACTION.CREATE : ACTION.DELETE;

      handleError(action, error.message);
    },
    onCompleted: (res, req) => {
      const namespace = req?.variables?.namespace;
      const { name, kind, selected } = req?.variables?.sources?.[0] || {};

      const count = req?.variables?.sources.length;
      const action = selected ? ACTION.CREATE : ACTION.DELETE;
      const fromOrIn = selected ? 'in' : 'from';

      if (count > 1) {
        handleComplete(action, `${count} sources were ${action.toLowerCase()}d ${fromOrIn} "${namespace}"`);
      } else {
        const id = { kind, name, namespace };
        if (!selected) removeNotifications(getSseTargetFromId(id, OVERVIEW_ENTITY_TYPES.SOURCE));
        if (!selected) setConfiguredSources({ ...configuredSources, [namespace]: configuredSources[namespace]?.filter((source) => source.name !== name) || [] });
        handleComplete(action, `source "${name}" was ${action.toLowerCase()}d ${fromOrIn} "${namespace}"`, selected ? id : undefined);
      }
    },
  });

  const [updateSource, uState] = useMutation<{ updateK8sActualSource: boolean }>(UPDATE_K8S_ACTUAL_SOURCE, {
    onError: (error) => handleError(ACTION.UPDATE, error.message),
    onCompleted: (res, req) => {
      const id = req?.variables?.sourceId;
      const name = id?.name;
      handleComplete(ACTION.UPDATE, `source "${name}" was updated`, id);
    },
  });

  const persistNamespaces = async (items: { [key: string]: boolean }) => {
    for (const [namespace, futureSelected] of Object.entries(items)) {
      await persistNamespace({ name: namespace, futureSelected });
    }
  };

  const persistSources = async (items: { [key: string]: K8sActualSource[] }, selected: boolean) => {
    for (const [namespace, sources] of Object.entries(items)) {
      await createOrDeleteSources({
        variables: {
          namespace,
          sources: sources.map((source) => ({
            kind: source.kind,
            name: source.name,
            selected,
          })),
        },
      });
    }
  };

  return {
    loading: cdState.loading || uState.loading,
    sources: data?.computePlatform.k8sActualSources || [],

    createSources: async (selectAppsList: { [key: string]: K8sActualSource[] }, futureSelectAppsList: { [key: string]: boolean }) => {
      await persistNamespaces(futureSelectAppsList);
      await persistSources(selectAppsList, true);
    },
    updateSource: async (sourceId: WorkloadId, patchSourceRequest: PatchSourceRequestInput) => await updateSource({ variables: { sourceId, patchSourceRequest } }),
    deleteSources: async (selectAppsList: { [key: string]: K8sActualSource[] }) => await persistSources(selectAppsList, false),
  };
};
