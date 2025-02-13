'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ROUTES } from '@/utils';
import { Stepper } from '@odigos/ui-components';
import { OnboardingStepperWrapper } from '@/components';
import SetupHeader from '@/components/lib-imports/setup-header';
import { DestinationSelectionForm, useSetupStore } from '@odigos/ui-containers';
import { ENTITY_TYPES, MOCK_DESTINATION_CATEGORIES, MOCK_POTENTIAL_DESTINATIONS, sleep } from '@odigos/ui-utils';

export default function Page() {
  const router = useRouter();
  const { configuredSources } = useSetupStore();

  const [isLoading, setIsLoading] = useState(false);
  const [isTested, setIsTested] = useState(false);

  return (
    <>
      <SetupHeader entityType={ENTITY_TYPES.DESTINATION} isLoading={isLoading} setIsLoading={setIsLoading} />

      <OnboardingStepperWrapper>
        <Stepper
          currentStep={3}
          data={[
            { stepNumber: 1, title: 'INSTALLATION' },
            { stepNumber: 2, title: 'SOURCES' },
            { stepNumber: 3, title: 'DESTINATIONS' },
          ]}
        />
      </OnboardingStepperWrapper>

      <DestinationSelectionForm
        categories={MOCK_DESTINATION_CATEGORIES}
        potentialDestinations={MOCK_POTENTIAL_DESTINATIONS}
        createDestination={async () => {}}
        isLoading={false}
        testLoading={false}
        testConnection={() => sleep(1500).then(() => setIsTested(true))}
        testResult={isTested ? { succeeded: false } : undefined}
        isSourcesListEmpty={!Object.values(configuredSources).some((sources) => !!sources.length)}
        goToSources={() => router.push(ROUTES.CHOOSE_SOURCES)}
      />
    </>
  );
}
