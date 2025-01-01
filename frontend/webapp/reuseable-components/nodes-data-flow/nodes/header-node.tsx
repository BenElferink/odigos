import React, { useMemo } from 'react';
import { PlusIcon } from '@/assets';
import styled from 'styled-components';
import { useSourceCRUD } from '@/hooks';
import type { Node, NodeProps } from '@xyflow/react';
import { NODE_TYPES, OVERVIEW_ENTITY_TYPES } from '@/types';
import { useAppStore, useModalStore, usePendingStore } from '@/store';
import { Badge, Checkbox, IconButton, Text } from '@/reuseable-components';

interface Props
  extends NodeProps<
    Node<
      {
        nodeWidth: number;
        icon: string;
        title: string;
        tagValue: number;
      },
      NODE_TYPES.HEADER
    >
  > {}

const Container = styled.div<{ $nodeWidth: Props['data']['nodeWidth'] }>`
  width: ${({ $nodeWidth }) => `${$nodeWidth}px`};
  padding: 12px 0px 16px 0px;
  gap: 8px;
  display: flex;
  align-items: center;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
`;

const Title = styled(Text)`
  color: ${({ theme }) => theme.text.grey};
`;

const ActionsWrapper = styled.div`
  margin-left: auto;
  margin-right: 16px;
`;

const HeaderNode: React.FC<Props> = ({ id: nodeId, data }) => {
  const { nodeWidth, title, icon: Icon, tagValue } = data;
  const entityType = nodeId.split('-')[0] as OVERVIEW_ENTITY_TYPES;
  const isSources = entityType === 'source';

  const { configuredSources, setConfiguredSources } = useAppStore();
  const { setCurrentModal } = useModalStore();
  const { isThisPending } = usePendingStore();
  const { sources } = useSourceCRUD();

  const totalSelectedSources = useMemo(() => {
    let num = 0;

    Object.values(configuredSources).forEach((selectedSources) => {
      num += selectedSources.length;
    });

    return num;
  }, [configuredSources]);

  const renderActions = () => {
    if (!isSources || !sources.length) return null;

    const onSelect = (bool: boolean) => {
      if (bool) {
        const payload = {};

        sources.forEach((source) => {
          const id = { namespace: source.namespace, name: source.name, kind: source.kind };
          const isPending = isThisPending({ entityType, entityId: id });

          if (!isPending) {
            if (!payload[source.namespace]) {
              payload[source.namespace] = [source];
            } else {
              payload[source.namespace].push(source);
            }
          }
        });

        setConfiguredSources(payload);
      } else {
        setConfiguredSources({});
      }
    };

    return (
      <ActionsWrapper>
        <Checkbox value={sources.length === totalSelectedSources} onChange={onSelect} />
      </ActionsWrapper>
    );
  };

  return (
    <Container $nodeWidth={nodeWidth} className='nowheel nodrag'>
      {Icon && <Icon />}
      <Title size={14}>{title}</Title>
      <Badge label={tagValue} />
      <Badge
        filled={!!tagValue}
        label={
          <IconButton size={24} onClick={() => setCurrentModal(entityType)}>
            <PlusIcon />
          </IconButton>
        }
      />

      {renderActions()}
    </Container>
  );
};

export default HeaderNode;
