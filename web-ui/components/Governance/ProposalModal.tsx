import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Button,
  useDisclosure,
  Flex,
  Text,
  Box,
  Center,
  Divider,
  Heading,
  useColorMode,
  useColorModeValue,
} from '@chakra-ui/react';
import { cosmos } from 'interchain-query';
import { Proposal } from 'interchain-query/cosmos/gov/v1/gov';
import React, { useMemo, useState } from 'react';
import { PieChart } from 'react-minimal-pie-chart';

import { Votes } from '@/hooks';
import { decodeUint8Arr, exponentiate, formatDate, getCoin, getExponent, getPercentage } from '@/utils';

import { VoteResult, TimeDisplay, VoteRatio, NewLineText, StatusBadge, VoteOption } from './common';
import { VoteColor } from './ProposalCard';
import { VoteModal } from './VoteModal';

const ProposalStatus = cosmos.gov.v1beta1.ProposalStatus;

export const ProposalModal = ({
  isOpen,
  onClose,
  proposal,
  chainName,
  quorum,
  bondedTokens,
  votes,
  updateVotes,
}: {
  isOpen: boolean;
  onClose: () => void;
  proposal: Proposal;
  chainName: string;
  quorum: number | undefined;
  bondedTokens: string | undefined;
  votes: Votes | undefined;
  updateVotes: () => void;
}) => {
  const [showMore, setShowMore] = useState(false);
  const voteModalControl = useDisclosure();
  const { colorMode } = useColorMode();

  const coin = getCoin(chainName);
  const exponent = getExponent(chainName);

  const chartData = [
    {
      title: 'YES',
      value: Number(proposal.finalTallyResult?.yesCount),
      color: VoteColor.YES,
    },
    {
      title: 'NO',
      value: Number(proposal.finalTallyResult?.noCount),
      color: VoteColor.NO,
    },
    {
      title: 'NWV',
      value: Number(proposal.finalTallyResult?.noWithVetoCount),
      color: VoteColor.NWV,
    },
    {
      title: 'ABSTAIN',
      value: Number(proposal.finalTallyResult?.abstainCount),
      color: VoteColor.ABSTAIN,
    },
  ];

  const emptyChartData = [
    {
      title: 'NO VOTES YET',
      value: 100,
      color: VoteColor.ABSTAIN,
    },
  ];

  const totalVotes = useMemo(() => {
    if (!proposal.finalTallyResult) return 0;
    const total = Object.values(proposal.finalTallyResult).reduce((prev, cur) => prev + Number(cur), 0);
    return total ? total : 0;
  }, [proposal]);

  const vote = votes && proposal.finalTallyResult && votes?.[proposal.id.toString()];

  const isDepositPeriod = proposal.status === ProposalStatus.PROPOSAL_STATUS_DEPOSIT_PERIOD;

  const isVotingPeriod = proposal.status === ProposalStatus.PROPOSAL_STATUS_VOTING_PERIOD;

  const turnout = totalVotes / Number(bondedTokens);

  const minStakedTokens = quorum && exponentiate(quorum * Number(bondedTokens), -exponent).toFixed(6);

  const uint8ArrayValue = proposal.messages[0].value;
  const propinfo = decodeUint8Arr(uint8ArrayValue);

  console.log(propinfo);

  const getTitleFromDecoded = (decodedStr: string) => {
    return decodedStr.slice(0, 250).match(/[A-Z][A-Za-z].*(?=\u0012)/)?.[0];
  };

  const getDescriptionFromProposal = (decodedData: string): string => {
    const lines = decodedData.split('\n');
    return lines.slice(4).join('\n') || '';
  };

  const title = getTitleFromDecoded(propinfo);
  const description = getDescriptionFromProposal(propinfo);

  const descriptionRenderer = () => {
    if (!description) return '';

    if (description.length > 200) {
      return showMore ? description : `${description.slice(0, 200)}...`;
    }

    return description;
  };

  const renderedDescription = descriptionRenderer();

  return (
    <>
      <VoteModal
        chainName={chainName}
        modalControl={voteModalControl}
        updateVotes={updateVotes}
        title={title || ''}
        vote={vote || 0}
        proposalId={proposal.id}
      />

      <Modal
        isOpen={isOpen}
        onClose={() => {
          onClose();
          setShowMore(false);
        }}
        isCentered
        size="3xl"
      >
        <ModalOverlay />
        <>
          <ModalContent
            pr={2}
            sx={{
              '&::-webkit-scrollbar': {
                width: '8px',
              },
              '&::-webkit-scrollbar-thumb': {
                backgroundColor: 'complimentary.900',
                borderRadius: '4px',
              },
              '&::-webkit-scrollbar-track': {
                backgroundColor: 'rgba(255,128,0, 0.25)',
                borderRadius: '10px',
              },
            }}
            bgColor="#1A1A1A"
            maxH="80vh"
            overflowY="scroll"
            px={2}
          >
            <ModalHeader>
              <Flex gap={2} mt={1} mb={2} alignItems="center">
                <Center h="min-content" transform="translateY(1px)">
                  <StatusBadge status={proposal.status} />
                </Center>
                {vote && <VoteResult voteOption={vote} />}
              </Flex>
              <Text color="white">{`#${proposal.id} ${title}`}</Text>
            </ModalHeader>
            <ModalCloseButton color="white" />
            <ModalBody>
              <Flex justifyContent="space-between" alignItems="center">
                <TimeDisplay title="Submit Time" time={formatDate(proposal.submitTime)} />
                <TimeDisplay title="Voting Starts" time={isDepositPeriod ? 'Not Specified Yet' : formatDate(proposal.votingStartTime)} />
                <TimeDisplay title="Voting Ends" time={isDepositPeriod ? 'Not Specified Yet' : formatDate(proposal.votingEndTime)} />
                <Button
                  isDisabled={!isVotingPeriod}
                  _hover={{
                    bgColor: '#181818',
                  }}
                  w="140px"
                  onClick={voteModalControl.onOpen}
                >
                  {vote ? 'Edit Vote' : 'Vote'}
                </Button>
              </Flex>
              <Center my={4} />
              <Divider bgColor="complimentary.500" />
              <Box mt={4}>
                <Heading fontSize="sm" color="white" mb={4}>
                  Vote Details
                </Heading>

                <Flex bgColor="blackAlpha.100" py={6} borderRadius="lg">
                  <Center px={6}>
                    <PieChart
                      data={totalVotes ? chartData : emptyChartData}
                      lineWidth={14}
                      paddingAngle={totalVotes ? 1 : 0}
                      style={{
                        height: '160px',
                        width: '160px',
                      }}
                      label={({ dataEntry }) => {
                        const { value, title, percentage } = dataEntry;
                        if (!totalVotes) return title;

                        const maxValue = Math.max(...chartData.map((item) => item.value));

                        if (value !== maxValue) return '';
                        return `${title} ${percentage.toFixed(2)}%`;
                      }}
                      labelStyle={{
                        fontSize: '10px',
                        fill: totalVotes ? chartData.sort((a, b) => b.value - a.value)[0].color : VoteColor.ABSTAIN,
                        fontWeight: 'bold',
                      }}
                      labelPosition={0}
                    />
                  </Center>

                  <Box pr={2}>
                    <Text
                      color={turnout > (quorum || 0) ? 'green.500' : 'gray.400'}
                      borderColor={turnout > (quorum || 0) ? 'green.500' : 'gray.400'}
                      fontWeight="bold"
                      border="1px solid"
                      w="fit-content"
                      px={2}
                      borderRadius="4px"
                    >
                      Turnout: {(turnout * 100).toFixed(2)}%
                    </Text>
                    {quorum && (
                      <Text color="white" fontSize="sm" my={2} fontWeight="semibold">
                        {`Minimum of staked ${minStakedTokens} ${coin.symbol}(${quorum * 100}%) need to vote
                    for this proposal to pass.`}
                      </Text>
                    )}
                    <Flex wrap="wrap" gap={4}>
                      <VoteRatio
                        type={VoteOption.YES}
                        ratio={getPercentage(proposal.finalTallyResult?.yesCount, totalVotes)}
                        amount={exponentiate(proposal.finalTallyResult?.yesCount, -exponent).toFixed(2)}
                        token={coin.symbol}
                      />
                      <VoteRatio
                        type={VoteOption.NO}
                        ratio={getPercentage(proposal.finalTallyResult?.noCount, totalVotes)}
                        amount={exponentiate(proposal.finalTallyResult?.noCount, -exponent).toFixed(2)}
                        token={coin.symbol}
                      />
                      <VoteRatio
                        type={VoteOption.NWV}
                        ratio={getPercentage(proposal.finalTallyResult?.noWithVetoCount, totalVotes)}
                        amount={exponentiate(proposal.finalTallyResult?.noWithVetoCount, -exponent).toFixed(2)}
                        token={coin.symbol}
                      />
                      <VoteRatio
                        type={VoteOption.ABSTAIN}
                        ratio={getPercentage(proposal.finalTallyResult?.abstainCount, totalVotes)}
                        amount={exponentiate(proposal.finalTallyResult?.abstainCount, -exponent).toFixed(2)}
                        token={coin.symbol}
                      />
                    </Flex>
                  </Box>
                </Flex>
              </Box>

              <Box mt={4}>
                <Heading fontSize="sm" color="white">
                  Description
                </Heading>
                <NewLineText text={renderedDescription} />
                {description && description.length > 200 && (
                  <Button
                    _hover={{
                      bgColor: '#181818',
                    }}
                    onClick={() => setShowMore(!showMore)}
                    size="sm"
                  >
                    {showMore ? 'Show less' : 'Show more'}
                  </Button>
                )}
              </Box>
            </ModalBody>

            <ModalFooter>
              <Button _hover={{ bgColor: 'complimentary.900' }} mt={-4} color="white" variant="ghost" onClick={onClose}>
                Close
              </Button>
            </ModalFooter>
          </ModalContent>
        </>
      </Modal>
    </>
  );
};
