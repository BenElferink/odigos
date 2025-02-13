'use client';

import { useEffect } from 'react';
import { ROUTES } from '@/utils';
import { sleep } from '@odigos/ui-utils';
import { useRouter } from 'next/navigation';
import { CenterThis, FadeLoader } from '@odigos/ui-components';

export default function App() {
  const router = useRouter();

  useEffect(() => {
    sleep(1500).then(() => router.push(ROUTES.OVERVIEW));
  }, []);

  return (
    <CenterThis style={{ height: '100%' }}>
      <FadeLoader scale={2} />
    </CenterThis>
  );
}
