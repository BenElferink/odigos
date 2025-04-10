import React from 'react';
import { ActionDrawer, ActionModal, DestinationDrawer, DestinationModal, InstrumentationRuleDrawer, InstrumentationRuleModal, SourceDrawer, SourceModal } from '@odigos/ui-kit/containers';
import {
  useActionCRUD,
  useConfig,
  useDescribe,
  useDestinationCategories,
  useDestinationCRUD,
  useInstrumentationRuleCRUD,
  useNamespace,
  usePotentialDestinations,
  useSourceCRUD,
  useTestConnection,
} from '@/hooks';

const OverviewModalsAndDrawers = () => {
  const { isEnterprise } = useConfig();

  const { fetchDescribeSource } = useDescribe();
  const { fetchSingleNamespace } = useNamespace();
  const { categories } = useDestinationCategories();
  const { persistSources, updateSource } = useSourceCRUD();
  const { potentialDestinations } = usePotentialDestinations();
  const { createAction, updateAction, deleteAction } = useActionCRUD();
  const { testConnection } = useTestConnection();
  const { createDestination, updateDestination, deleteDestination } = useDestinationCRUD();
  const { createInstrumentationRule, updateInstrumentationRule, deleteInstrumentationRule } = useInstrumentationRuleCRUD();

  return (
    <>
      {/* modals */}
      <SourceModal groupNames={groupNames} fetchSingleNamespace={fetchSingleNamespace} persistSources={persistSources} />
      <DestinationModal isOnboarding={false} categories={categories} potentialDestinations={potentialDestinations} createDestination={createDestination} testConnection={testConnection} />
      <InstrumentationRuleModal isEnterprise={isEnterprise} createInstrumentationRule={createInstrumentationRule} />
      <ActionModal createAction={createAction} />

      {/* drawers */}
      <SourceDrawer persistSources={persistSources} updateSource={updateSource} fetchDescribeSource={fetchDescribeSource} />
      <DestinationDrawer categories={categories} updateDestination={updateDestination} deleteDestination={deleteDestination} testConnection={testConnection} />
      <InstrumentationRuleDrawer updateInstrumentationRule={updateInstrumentationRule} deleteInstrumentationRule={deleteInstrumentationRule} />
      <ActionDrawer updateAction={updateAction} deleteAction={deleteAction} />
    </>
  );
};

export { OverviewModalsAndDrawers };
