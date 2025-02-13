'use client';

import React, { useRef, useState } from 'react';
import { Stepper } from '@odigos/ui-components';
import { OnboardingStepperWrapper } from '@/components';
import { ENTITY_TYPES, MOCK_SOURCES } from '@odigos/ui-utils';
import SetupHeader from '@/components/lib-imports/setup-header';
import { SourceSelectionForm, type SourceSelectionFormRef } from '@odigos/ui-containers';

export default function Page() {
  const [selectedNamespace, setSelectedNamespace] = useState('');
  const onSelectNamespace = (ns: string) => setSelectedNamespace((prev) => (prev === ns ? '' : ns));

  const formRef = useRef<SourceSelectionFormRef>(null);

  return (
    <>
      <SetupHeader entityType={ENTITY_TYPES.SOURCE} formRef={formRef} />

      <OnboardingStepperWrapper>
        <Stepper
          currentStep={2}
          data={[
            { stepNumber: 1, title: 'INSTALLATION' },
            { stepNumber: 2, title: 'SOURCES' },
            { stepNumber: 3, title: 'DESTINATIONS' },
          ]}
        />
      </OnboardingStepperWrapper>

      <SourceSelectionForm
        ref={formRef}
        componentType='FAST'
        isModal={false}
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
        onSelectNamespace={onSelectNamespace}
      />
    </>
  );
}
