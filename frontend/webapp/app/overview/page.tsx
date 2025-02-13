'use client';

import React from 'react';
import { MainContent } from '@/components';
import OverviewHeader from '@/components/lib-imports/overview-header';
import { DataFlow, DataFlowActionsMenu, MultiSourceControl } from '@odigos/ui-containers';
import OverviewModalsAndDrawers from '@/components/lib-imports/overview-modals-and-drawers';
import { MOCK_ACTIONS, MOCK_DESTINATIONS, MOCK_INSTRUMENTATION_RULES, MOCK_SOURCES } from '@odigos/ui-utils';

export default function Page() {
  return (
    <>
      <OverviewHeader />
      <MainContent>
        <DataFlowActionsMenu
          namespaces={[{ name: 'default' }, { name: 'kube-public' }, { name: 'tracing' }]}
          sources={MOCK_SOURCES}
          destinations={MOCK_DESTINATIONS}
          actions={MOCK_ACTIONS}
          instrumentationRules={MOCK_INSTRUMENTATION_RULES}
        />
        <DataFlow
          heightToRemove='176px'
          sources={MOCK_SOURCES}
          sourcesLoading={false}
          sourcesTotalCount={MOCK_SOURCES.length}
          destinations={MOCK_DESTINATIONS}
          destinationsLoading={false}
          destinationsTotalCount={MOCK_DESTINATIONS.length}
          actions={MOCK_ACTIONS}
          actionsLoading={false}
          actionsTotalCount={MOCK_ACTIONS.length}
          instrumentationRules={MOCK_INSTRUMENTATION_RULES}
          instrumentationRulesLoading={false}
          instrumentationRulesTotalCount={MOCK_INSTRUMENTATION_RULES.length}
          metrics={{ sources: [], destinations: [] }}
        />
        <MultiSourceControl totalSourceCount={MOCK_SOURCES.length} uninstrumentSources={() => {}} />
      </MainContent>
      <OverviewModalsAndDrawers />
    </>
  );
}
