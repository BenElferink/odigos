import styled from "styled-components";

export const SourcesOptionMenuWrapper = styled.section`
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 40px 0 0 0;
  @media screen and (max-width: 1500px) {
    flex-wrap: wrap;
    width: 70%;
  }
`;

export const FilterMenuWrapper = styled.div`
  display: flex;
  gap: 16px;
  align-items: center;
`;

export const DropdownWrapper = styled.div`
  display: flex;
  position: inherit;
  align-items: center;
  gap: 12px;
`;

export const CheckboxWrapper = styled.div`
  display: flex;
  gap: 8px;
  align-items: center;
  min-width: 180px;
`;

export const SwitcherWrapper = styled.div`
  min-width: 90px;
  margin-left: 24px;
`;
