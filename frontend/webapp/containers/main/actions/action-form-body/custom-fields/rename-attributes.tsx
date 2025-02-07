import React, { useMemo } from 'react';
import { safeJsonParse } from '@odigos/ui-utils';
import type { RenameAttributesSpec } from '@/types';
import { KeyValueInputsList } from '@odigos/ui-components';

type Props = {
  value: string;
  setValue: (value: string) => void;
  errorMessage?: string;
};

type Parsed = RenameAttributesSpec;

const RenameAttributes: React.FC<Props> = ({ value, setValue, errorMessage }) => {
  const mappedValue = useMemo(() => Object.entries(safeJsonParse<Parsed>(value, { renames: {} }).renames).map(([k, v]) => ({ key: k, value: v })), [value]);

  const handleChange = (
    arr: {
      key: string;
      value: string;
    }[],
  ) => {
    const payload: Parsed = {
      renames: {},
    };

    arr.forEach((obj) => {
      payload.renames[obj.key] = obj.value;
    });

    const str = !!arr.length ? JSON.stringify(payload) : '';

    setValue(str);
  };

  return <KeyValueInputsList title='Attributes to rename' value={mappedValue} onChange={handleChange} required errorMessage={errorMessage} />;
};

export default RenameAttributes;
