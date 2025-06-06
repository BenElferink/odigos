import { useEffect } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useDataStreamStore } from '@odigos/ui-kit/store';
import { ROUTES, SKIP_TO_SUMMERY_QUERY_PARAM } from '@/utils';
import { DEFAULT_DATA_STREAM_NAME } from '@odigos/ui-kit/constants';

export const useSetupHelpers = () => {
  const router = useRouter();
  const pathname = usePathname();
  const params = useSearchParams();

  const isSkipToSummary = !!params.get(SKIP_TO_SUMMERY_QUERY_PARAM);
  const skipToSummaryQuerystring = isSkipToSummary ? `?${SKIP_TO_SUMMERY_QUERY_PARAM}=true` : '';

  // If we do not want to show the "go to summary button" in setup pages, we have to pass "undefined" as prop
  const onClickSummary = isSkipToSummary ? () => router.push(ROUTES.SETUP_SUMMARY) : undefined;
  const onClickRouteFromSummary = (path: string) => router.push(path + skipToSummaryQuerystring);

  const { selectedStreamName, setSelectedStreamName } = useDataStreamStore();

  useEffect(() => {
    if (!selectedStreamName && pathname !== ROUTES.CHOOSE_STREAM) setSelectedStreamName(DEFAULT_DATA_STREAM_NAME);
  }, [selectedStreamName]);

  return { onClickSummary, onClickRouteFromSummary };
};
