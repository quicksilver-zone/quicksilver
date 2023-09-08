import {
  Box,
  Center,
  Flex,
  Grid,
  GridItem,
  Spacer,
  Text,
  useColorMode,
} from '@chakra-ui/react';
import dayjs from 'dayjs';
import { cosmos } from 'interchain-query';
import { Proposal } from 'interchain-query/cosmos/gov/v1beta1/gov';
import React, { useMemo } from 'react';

import { Votes } from '@/hooks';
import { getPercentage } from '@/utils';

import {
  StatusBadge,
  VotedBadge,
} from './common';

enum VoteOption {
  YES = 'YES',
  NO = 'NO',
  NWV = 'NWV',
  ABSTAIN = 'ABSTAIN',
}

const ProposalStatus =
  cosmos.gov.v1beta1.ProposalStatus;

export const VoteColor: {
  [key in VoteOption]: string;
} = {
  [VoteOption.YES]: '#17a572',
  [VoteOption.NO]: '#ce4256',
  [VoteOption.NWV]: '#ff5b6d',
  [VoteOption.ABSTAIN]: '#546198',
};

export const ProposalCard = ({
  proposal,
  handleClick,
  votes,
}: {
  proposal: Proposal;
  handleClick: () => void;
  votes: Votes | undefined;
}) => {
  const { colorMode } = useColorMode();

  const totalVotes = useMemo(() => {
    if (!proposal.finalTallyResult) return 0;
    const total = Object.values(
      proposal.finalTallyResult,
    ).reduce(
      (prev, cur) => prev + Number(cur),
      0,
    );
    return total ? total : 0;
  }, [proposal]);

  const isVoted =
    votes &&
    votes[proposal.proposalId.toString()];

  const getTitleFromProposal = (
    proposal: Proposal,
  ): string | undefined => {
    if (
      proposal.content &&
      'title' in proposal.content
    ) {
      return proposal.content.title;
    }
    return undefined;
  };

  const title = getTitleFromProposal(proposal);

  return (
    <Grid
      h="120px"
      py={4}
      templateColumns="repeat(13, 1fr)"
      bgColor="rgba(255,255,255,0.1)"
      backdropFilter="blur(30px)"
      borderColor="gray.400"
      borderRadius={10}
      transition="all 0.2s linear"
      _hover={{
        backgroundColor: 'rgba(255,255,255,0.25)',
        cursor: 'pointer',
      }}
      onClick={handleClick}
    >
      <GridItem colSpan={2}>
        <Center color="white" w="100%" h="100%">
          #{' '}
          {proposal.proposalId
            .toString()
            .padStart(6, '0')}
        </Center>
      </GridItem>
      <GridItem colSpan={9} py={2}>
        <Flex flexDirection="column" h="100%">
          <Flex gap={2} alignItems="center">
            <Text color="white" fontSize="lg">
              {title || ''}
            </Text>
            {isVoted && <VotedBadge />}
          </Flex>
          <Spacer />
          <Flex flexDirection="column" h="44%">
            <Flex
              alignItems="center"
              fontSize="sm"
            >
              <Text color="white">
                {proposal.status ===
                ProposalStatus.PROPOSAL_STATUS_DEPOSIT_PERIOD
                  ? 'Deposit'
                  : 'Voting'}
                &nbsp;end time: &nbsp;
              </Text>
              <Text
                color="white"
                fontWeight="semibold"
              >
                {dayjs(
                  proposal.status ===
                    ProposalStatus.PROPOSAL_STATUS_DEPOSIT_PERIOD
                    ? proposal.depositEndTime
                    : proposal.votingEndTime,
                ).format('YYYY-MM-DD hh:mm')}
              </Text>
            </Flex>
            <Spacer />
            {totalVotes ? (
              <Flex gap="1px">
                <Box
                  w={getPercentage(
                    proposal.finalTallyResult
                      ?.yes,
                    totalVotes,
                  )}
                  h="3px"
                  bgColor={VoteColor.YES}
                />
                <Box
                  w={getPercentage(
                    proposal.finalTallyResult?.no,
                    totalVotes,
                  )}
                  h="3px"
                  bgColor={VoteColor.NO}
                />
                <Box
                  w={getPercentage(
                    proposal.finalTallyResult
                      ?.noWithVeto,
                    totalVotes,
                  )}
                  h="3px"
                  bgColor={VoteColor.NWV}
                />
                <Box
                  w={getPercentage(
                    proposal.finalTallyResult
                      ?.abstain,
                    totalVotes,
                  )}
                  h="3px"
                  bgColor={VoteColor.ABSTAIN}
                />
              </Flex>
            ) : (
              <Box
                w="100%"
                h="3px"
                bgColor={
                  colorMode === 'light'
                    ? 'gray.200'
                    : 'gray.600'
                }
              />
            )}
          </Flex>
        </Flex>
      </GridItem>
      <GridItem colSpan={2}>
        <Flex
          w="100%"
          h="100%"
          alignItems="center"
          px={4}
          justifyContent="center"
        >
          <StatusBadge status={proposal.status} />
        </Flex>
      </GridItem>
    </Grid>
  );
};
