import {
  Badge,
  Box,
  Center,
  Flex,
  Icon,
  Spinner,
  Stack,
  Text,
  useColorModeValue,
} from '@chakra-ui/react';
import styled from '@emotion/styled';
import { cosmos } from 'interchain-query';
import { ProposalStatus } from 'interchain-query/cosmos/gov/v1beta1/gov';
import { IconType } from 'react-icons';
import {
  AiFillCheckCircle,
  AiFillCloseCircle,
  AiFillMinusCircle,
} from 'react-icons/ai';
import ReactMarkdown from 'react-markdown';

import { VoteColor } from './ProposalCard';

/* eslint-disable */ // @ts-ignore
const MarkdownStyled = styled(ReactMarkdown)`
  color: white;

  a {
    color: white;
  }
`;

export const Loader = () => (
  <Center
    w="100%"
    h="200px"
    bgColor="rgba(214, 219, 220, 0.1)"
    borderRadius="xl"
  >
    <Spinner thickness="3px" color="complimentary.900" size="lg" speed="0.4s" />
  </Center>
);

export const DisconnectedContent = () => (
  <Center
    w="100%"
    h="100px"
    bgColor="rgba(214, 219, 220, 0.1)"
    borderRadius="xl"
  >
    <Text fontSize="lg" color="white">
      Please connect your wallet to see the proposals
    </Text>
  </Center>
);

export const TimeDisplay = ({
  title,
  time,
}: {
  title: string;
  time: string;
}) => (
  <Stack spacing="0.5">
    <Text fontSize="sm" fontWeight="semibold" color="white">
      {title}
    </Text>
    <Text fontWeight="semibold" color="white" letterSpacing="tight">
      {time}
    </Text>
  </Stack>
);

const VoteType = cosmos.gov.v1beta1.VoteOption;

export enum VoteOption {
  YES = 'YES',
  NO = 'NO',
  NWV = 'NWV',
  ABSTAIN = 'ABSTAIN',
}

export const VoteRatio = ({
  type,
  ratio,
  amount,
  token,
}: {
  type: keyof typeof VoteOption;
  ratio: string;
  amount: string;
  token: string;
}) => (
  <Box
    py={2}
    px={4}
    border="1px solid"
    borderRadius="md"
    borderColor="complimentary.900"
    bgColor="primary.900"
    w="200px"
  >
    <Text color={VoteColor[type]} fontWeight="bold">
      {type} {ratio}
    </Text>
    <Text fontSize="sm" fontWeight="semibold" color="white">
      {amount} {token}
    </Text>
  </Box>
);

export const VoteResult = ({ voteOption }: { voteOption: number }) => {
  let optionConfig: {
    color: string;
    icon: IconType;
    option: string;
  } = {
    color: VoteColor.YES,
    icon: AiFillCheckCircle,
    option: 'Yes',
  };

  switch (voteOption) {
    case VoteType.VOTE_OPTION_YES:
      break;
    case VoteType.VOTE_OPTION_NO:
      optionConfig = {
        color: VoteColor.NO,
        icon: AiFillCloseCircle,
        option: 'No',
      };
      break;
    case VoteType.VOTE_OPTION_NO_WITH_VETO:
      optionConfig = {
        color: VoteColor.NWV,
        icon: AiFillCloseCircle,
        option: 'NoWithVeto',
      };
      break;
    case VoteType.VOTE_OPTION_ABSTAIN:
      optionConfig = {
        color: VoteColor.ABSTAIN,
        icon: AiFillMinusCircle,
        option: 'Abstain',
      };
      break;
    default:
      break;
  }

  return (
    <>
      <Text
        fontSize="sm"
        fontWeight="semibold"
        color={useColorModeValue('gray.600', 'gray.400')}
      >
        You Voted
      </Text>
      <Flex color={optionConfig.color} alignItems="center" gap="2px">
        <Icon
          as={optionConfig.icon}
          fontSize="lg"
          transform="translateY(1px)"
        />
        <Text fontSize="sm" fontWeight="bold">
          {optionConfig.option}
        </Text>
      </Flex>
    </>
  );
};

export const NewLineText = ({ text }: { text: string }) => {
  let count = 0;
  return (
    <>
      {text.split('\\n').map((str) => (
        <Box lineHeight="taller" fontSize="sm" key={count++}>
          <MarkdownStyled linkTarget="_blank" className="markdown-text">
            {str}
          </MarkdownStyled>
        </Box>
      ))}
    </>
  );
};

export const StatusBadge = ({ status }: { status: number }) => {
  let statusConfig: {
    color: string;
    name: string;
  } = {
    color: 'purple',
    name: 'Deposit Period',
  };

  switch (status) {
    case ProposalStatus.PROPOSAL_STATUS_DEPOSIT_PERIOD:
      break;
    case ProposalStatus.PROPOSAL_STATUS_VOTING_PERIOD:
      statusConfig = {
        color: 'twitter',
        name: 'Voting Period',
      };
      break;
    case ProposalStatus.PROPOSAL_STATUS_PASSED:
      statusConfig = {
        color: 'green',
        name: 'Passed',
      };
      break;
    case ProposalStatus.PROPOSAL_STATUS_REJECTED:
      statusConfig = {
        color: 'red',
        name: 'Rejected',
      };
      break;
    default:
      break;
  }

  return (
    <Badge colorScheme={statusConfig.color} variant="subtle" borderRadius={4}>
      <Flex alignItems="center">{statusConfig.name}</Flex>
    </Badge>
  );
};

export const VotedBadge = () => (
  <Badge colorScheme="purple" variant="solid" borderRadius={4} h="min-content">
    Voted
  </Badge>
);
