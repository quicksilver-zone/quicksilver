import {
  Text,
  Input,
  InputGroup,
  InputRightElement,
  Button,
} from '@chakra-ui/react';
import { ChangeEvent, useState } from 'react';

import { StatBox } from '@/components/Staking/modals/modalElements';

export const InputBox = ({
  label,
  token,
  value,
  onChange,
  onMaxClick,
  isMaxBtnLoading = false,
}: {
  label: string;
  token: string;
  value: number | string;
  onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onMaxClick: () => void;
  isMaxBtnLoading?: boolean;
}) => (
  <InputGroup mt={2}>
    <Input
      _active={{
        borderColor: 'complimentary.900',
      }}
      _selected={{
        borderColor: 'complimentary.900',
      }}
      _hover={{
        borderColor: 'complimentary.900',
      }}
      _focus={{
        borderColor: 'complimentary.900',
        boxShadow: '0 0 0 3px #FF8000',
      }}
      color="complimentary.900"
      textAlign={'right'}
      placeholder="amount"
      type="number"
      value={value}
      onChange={onChange}
    />
  </InputGroup>
);

export const useInputBox = (maxAmount?: number | string) => {
  const [amount, setAmount] = useState<number | string>('');
  const [max, setMax] = useState<number | string>(maxAmount || 0);

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (Number(e.target.value) > max) {
      setAmount(max);
      return;
    }

    if (e.target.value === '') {
      setAmount('');
      return;
    }

    setAmount(+Number(e.target.value).toFixed(6));
  };

  const renderInputBox = (
    label: string,
    token: string,
    onMaxClick?: () => void,
    isLoading?: boolean,
  ) => {
    return (
      <InputBox
        label={label}
        token={token}
        value={amount}
        isMaxBtnLoading={isLoading}
        onChange={(e) => handleInputChange(e)}
        onMaxClick={() => (onMaxClick ? onMaxClick() : setAmount(max))}
      />
    );
  };

  return { renderInputBox, amount, setAmount, setMax };
};
