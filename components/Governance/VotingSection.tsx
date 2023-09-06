import {
  Box,
  Heading,
  Stack,
  useDisclosure,
} from '@chakra-ui/react';
import { ChainName } from '@cosmos-kit/core';
import { useChain } from '@cosmos-kit/react';
import { Proposal } from 'interchain-query/cosmos/gov/v1/gov';
import React, { useState } from 'react';

import { useVotingData } from '@/hooks';

import {
  DisconnectedContent,
  Loader,
} from './common';
import { ProposalCard } from './ProposalCard';
import { ProposalModal } from './ProposalModal';

export const VotingSection = ({
  chainName,
}: {
  chainName: ChainName;
}) => {
  const [selectedProposal, setSelectedProposal] =
    useState<Proposal>();

  const { address } = useChain(chainName);
  const { isOpen, onOpen, onClose } =
    useDisclosure();
  const { data, isLoading, refetch } =
    useVotingData(chainName);

  const content = address ? (
    <Stack spacing={4}>
      {data?.proposals?.map((proposal) => (
        <ProposalCard
          proposal={proposal}
          votes={data?.votes}
          handleClick={() => {
            onOpen();
            setSelectedProposal(proposal);
          }}
          key={proposal.submitTime?.getTime()}
        />
      ))}
    </Stack>
  ) : (
    <DisconnectedContent />
  );

  return (
    <>
      <Box mb={16}>
        <Heading as="h1" size="md" mb={4}>
          Proposals
        </Heading>
        {isLoading ? <Loader /> : content}
      </Box>
      {selectedProposal && (
        <ProposalModal
          proposal={selectedProposal}
          quorum={data?.quorum}
          bondedTokens={data?.bondedTokens}
          isOpen={isOpen}
          chainName={chainName}
          onClose={onClose}
          votes={data?.votes}
          updateVotes={refetch}
        />
      )}
    </>
  );
};
