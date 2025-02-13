import React, { useState } from 'react';
import { ActionDrawer, ActionModal, DestinationDrawer, DestinationModal, InstrumentationRuleDrawer, InstrumentationRuleModal, SourceDrawer, SourceModal } from '@odigos/ui-containers';
import { MOCK_ACTIONS, MOCK_DESCRIBE_SOURCE, MOCK_DESTINATION_CATEGORIES, MOCK_DESTINATIONS, MOCK_INSTRUMENTATION_RULES, MOCK_POTENTIAL_DESTINATIONS, MOCK_SOURCES, sleep } from '@odigos/ui-utils';

const OverviewModalsAndDrawers = () => {
  const [selectedNamespace, setSelectedNamespace] = useState('');
  const [isTested, setIsTested] = useState(false);

  return (
    <>
      {/* modals */}
      <SourceModal
        namespaces={[
          { name: 'default', selected: false },
          { name: 'kube-public', selected: false },
          { name: 'tracing', selected: false },
        ]}
        namespace={
          !!selectedNamespace
            ? {
                name: selectedNamespace,
                selected: false,
                sources: MOCK_SOURCES,
              }
            : undefined
        }
        namespacesLoading={false}
        selectedNamespace={selectedNamespace}
        setSelectedNamespace={setSelectedNamespace}
        persistSources={async () => {}}
      />
      <DestinationModal
        isOnboarding={false}
        categories={MOCK_DESTINATION_CATEGORIES}
        potentialDestinations={MOCK_POTENTIAL_DESTINATIONS}
        createDestination={async () => {}}
        testLoading={false}
        testConnection={() => sleep(1500).then(() => setIsTested(true))}
        testResult={isTested ? { succeeded: false } : undefined}
      />
      <InstrumentationRuleModal isEnterprise={false} createInstrumentationRule={() => {}} />
      <ActionModal createAction={() => {}} />

      {/* drawers */}
      <SourceDrawer sources={MOCK_SOURCES} persistSources={async () => {}} updateSource={async () => {}} describe={MOCK_DESCRIBE_SOURCE} />
      <DestinationDrawer
        categories={MOCK_DESTINATION_CATEGORIES}
        destinations={MOCK_DESTINATIONS}
        updateDestination={async () => {}}
        deleteDestination={async () => {}}
        testLoading={false}
        testConnection={() => sleep(1500).then(() => setIsTested(true))}
        testResult={isTested ? { succeeded: false } : undefined}
      />
      <InstrumentationRuleDrawer instrumentationRules={MOCK_INSTRUMENTATION_RULES} updateInstrumentationRule={() => {}} deleteInstrumentationRule={() => {}} />
      <ActionDrawer actions={MOCK_ACTIONS} updateAction={() => {}} deleteAction={() => {}} />
    </>
  );
};

export default OverviewModalsAndDrawers;
