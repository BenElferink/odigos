import { useLazyQuery } from '@apollo/client';
import { TEST_CONNECTION_MUTATION } from '@/graphql';
import { useNotificationStore } from '@odigos/ui-kit/store';
import { Crud, StatusType, type DestinationFormData } from '@odigos/ui-kit/types';

interface TestConnectionResponse {
  succeeded: boolean;
  statusCode: number;
  destinationType: string;
  message: string;
  reason: string;
}

export const useTestConnection = () => {
  const { addNotification } = useNotificationStore();

  const notifyUser = (type: StatusType, title: string, message: string, hideFromHistory?: boolean) => {
    addNotification({ type, title, message, hideFromHistory });
  };

  const [testConnection] = useLazyQuery<{ testConnectionForDestination: TestConnectionResponse }, { destination: DestinationFormData }>(TEST_CONNECTION_MUTATION, {
    onError: (error) => notifyUser(StatusType.Error, error.name || Crud.Read, error.cause?.message || error.message),
  });

  return {
    testConnection,
  };
};
