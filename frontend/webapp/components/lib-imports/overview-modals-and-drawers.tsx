import React, { useCallback } from 'react';
import { EntityTypes } from '@odigos/ui-kit/types';
import { RuleFormContextProvider } from '@odigos/ui-kit/contexts';
import { useDrawerStore, useModalStore } from '@odigos/ui-kit/store';
import { ActionDrawer, DestinationDrawer, SourceDrawer } from '@odigos/ui-kit/containers';
import {
  AddActionDrawer,
  AddDestinationDrawer,
  AddRuleDrawer,
  AddSourceDrawer,
  AddActionFormContextProvider,
  AddDestinationFormContextProvider,
  AddSourceFormContextProvider,
  EditRuleDrawer,
} from '@odigos/ui-kit/containers/v2';
import {
  useActionCRUD,
  useDescribe,
  useDestinationCategories,
  useDestinationCategoriesLegacy,
  useDestinationCRUD,
  useInstrumentationRuleCRUD,
  useNamespace,
  usePotentialDestinations,
  useSourceCRUD,
  useTestConnection,
  useWorkloadUtils,
} from '@/hooks';

const OverviewModalsAndDrawers = () => {
  const { currentModal, setCurrentModal } = useModalStore();
  const { drawerType, drawerEntityId } = useDrawerStore();

  const { fetchDescribeSource } = useDescribe();
  const { testConnection } = useTestConnection();
  const { fetchNamespacesWithWorkloads } = useNamespace();
  const { categories } = useDestinationCategoriesLegacy();
  const { getDestinationCategories } = useDestinationCategories();
  const { getPotentialDestinations } = usePotentialDestinations();
  const { createActionV2, updateAction, deleteAction } = useActionCRUD();
  const { restartWorkloads, restartPod, recoverFromRollback } = useWorkloadUtils();
  const { createDestination, updateDestination, deleteDestination } = useDestinationCRUD();
  const { createInstrumentationRuleV2, updateInstrumentationRule, deleteInstrumentationRule } = useInstrumentationRuleCRUD();
  const { persistSources, persistSourcesV2, updateSource, fetchSourceById, fetchSourceLibraries, fetchPeerSources } = useSourceCRUD();

  const handleCloseModal = useCallback(() => setCurrentModal(''), [setCurrentModal]);

  return (
    <>
      {/* add drawers (v2) */}
      {currentModal === EntityTypes.Source && (
        <AddSourceFormContextProvider fetchNamespacesWithWorkloads={fetchNamespacesWithWorkloads}>
          <AddSourceDrawer onClose={handleCloseModal} persistSources={persistSourcesV2} withOverlay />
        </AddSourceFormContextProvider>
      )}
      {currentModal === EntityTypes.Destination && (
        <AddDestinationFormContextProvider>
          <AddDestinationDrawer
            onClose={handleCloseModal}
            getDestinationCategories={getDestinationCategories}
            getPotentialDestinations={getPotentialDestinations}
            testConnection={testConnection}
            createDestination={createDestination}
            updateDestination={updateDestination}
            withOverlay
          />
        </AddDestinationFormContextProvider>
      )}
      {currentModal === EntityTypes.InstrumentationRule && (
        <RuleFormContextProvider>
          <AddRuleDrawer onClose={handleCloseModal} createInstrumentationRule={createInstrumentationRuleV2} withOverlay />
        </RuleFormContextProvider>
      )}
      {currentModal === EntityTypes.Action && (
        <AddActionFormContextProvider>
          <AddActionDrawer onClose={handleCloseModal} createAction={createActionV2} withOverlay />
        </AddActionFormContextProvider>
      )}

      {/* edit drawers */}
      <SourceDrawer
        persistSources={persistSources}
        restartWorkloads={restartWorkloads}
        restartPod={restartPod}
        recoverFromRollback={recoverFromRollback}
        updateSource={updateSource}
        fetchSourceById={fetchSourceById}
        fetchSourceDescribe={fetchDescribeSource}
        fetchSourceLibraries={fetchSourceLibraries}
        fetchPeerSources={fetchPeerSources}
      />
      <DestinationDrawer categories={categories} updateDestination={updateDestination} deleteDestination={deleteDestination} testConnection={testConnection} />
      {drawerType === EntityTypes.InstrumentationRule && drawerEntityId && (
        <RuleFormContextProvider>
          <EditRuleDrawer onClose={handleCloseModal} ruleId={drawerEntityId} updateInstrumentationRule={updateInstrumentationRule} deleteInstrumentationRule={deleteInstrumentationRule} />
        </RuleFormContextProvider>
      )}
      <ActionDrawer updateAction={updateAction} deleteAction={deleteAction} />
    </>
  );
};

export { OverviewModalsAndDrawers };
