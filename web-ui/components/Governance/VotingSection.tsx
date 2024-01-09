import { ChevronDownIcon } from '@chakra-ui/icons';
import {
  Box,
  Flex,
  Heading,
  Input,
  Text,
  InputGroup,
  InputLeftElement,
  Stack,
  useDisclosure,
  Button,
  ButtonGroup,
  Spacer,
  MenuList,
  Menu,
  MenuButton,
  MenuItem,
  AccordionIcon,
} from '@chakra-ui/react';
import { ChainName } from '@cosmos-kit/core';
import { useChain } from '@cosmos-kit/react';
import { Proposal } from 'interchain-query/cosmos/gov/v1/gov';
import React, { useMemo, useState } from 'react';
import { FaSearch } from 'react-icons/fa';

import { useVotingData } from '@/hooks';

import { DisconnectedContent, Loader } from './common';
import { ProposalCard } from './ProposalCard';
import { ProposalModal } from './ProposalModal';
import { decodeUint8Arr } from '@/utils';

function RotateIcon({ isOpen }: { isOpen: boolean }) {
  return (
    <ChevronDownIcon
      color="complimentary.900"
      transform={isOpen ? 'rotate(180deg)' : 'none'}
      transition="transform 0.2s"
      h="25px"
      w="25px"
    />
  );
}

export const VotingSection = ({ chainName }: { chainName: ChainName }) => {
  const [selectedProposal, setSelectedProposal] = useState<Proposal>();
  const [selectedPeriodOption, setSelectedPeriodOption] = useState('All Periods');
  const [selectedProposalOption, setSelectedProposalOption] = useState('All Proposals');

  const { address } = useChain(chainName);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const { data, isLoading, refetch } = useVotingData(chainName);
  const [searchTerm, setSearchTerm] = useState('');

  const filteredProposals = useMemo(() => {
    if (!data?.proposals) return [];

    return data.proposals.filter((proposal) => {
      const decodedContent = decodeUint8Arr(proposal.messages[0].value);
      const contentToSearch = decodedContent.toLowerCase();
      const titleMatches = proposal.title.toLowerCase().includes(searchTerm.toLowerCase());
      const contentMatches = contentToSearch.includes(searchTerm.toLowerCase());

      let periodMatches = true;
      let proposalMatches = true;

      // Constants for proposal status (these values might be different in your application)
      const VOTING_PERIOD_STATUS = 2; // Example value for 'Voting Period'
      const PASSED_STATUS = 3; // Example value for 'Passed'
      const REJECTED_STATUS = 4; // Example value for 'Rejected'

      // Filter by period
      switch (selectedPeriodOption) {
        case 'Voting Period':
          periodMatches = proposal.status === VOTING_PERIOD_STATUS;
          break;
        case 'Passed':
          periodMatches = proposal.status === PASSED_STATUS;
          break;
        case 'Rejected':
          periodMatches = proposal.status === REJECTED_STATUS;
          break;
        default:
          periodMatches = true;
      }

      // Filter by proposal type (e.g., voted)
      if (selectedProposalOption === 'Voted') {
        const userVote = data?.votes?.[proposal.id.toString()];
        proposalMatches = userVote !== undefined;
      }

      return titleMatches && contentMatches && periodMatches && proposalMatches;
    });
  }, [data, searchTerm, selectedPeriodOption, selectedProposalOption]);

  const content = address ? (
    <Stack spacing={4}>
      {filteredProposals.map((proposal) => (
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
        <Flex mb={4} alignContent="center" alignItems="center" justifyContent={'space-between'} w="100%" flexDirection={'row'}>
          <InputGroup>
            <Input
              textAlign="right"
              type="text"
              color="white"
              borderColor="complimentary.1000"
              placeholder="proposal content..."
              fontWeight="light"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              width="35%"
              borderRadius={'4px'}
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
            />
            <InputLeftElement pointerEvents="none">
              <FaSearch color="orange" />
            </InputLeftElement>
          </InputGroup>
          <Spacer />
          <ButtonGroup spacing={6}>
            <Menu>
              {({ isOpen }) => (
                <>
                  <MenuButton
                    _hover={{
                      bgColor: 'rgba(255,128,0, 0.25)',
                    }}
                    px={2}
                    color="white"
                    w="150px"
                    as={Button}
                    variant="outline"
                    rightIcon={<RotateIcon isOpen={isOpen} />}
                  >
                    {selectedPeriodOption}
                  </MenuButton>
                  <MenuList minW="150px" borderColor="black" bgColor="#181818">
                    <MenuItem
                      borderRadius={'5px'}
                      _hover={{
                        bgColor: 'rgba(255,128,0, 0.25)',
                      }}
                      color="white"
                      bgColor="#181818"
                      onClick={() => setSelectedPeriodOption('All Periods')}
                    >
                      All Periods
                    </MenuItem>
                    <MenuItem
                      borderRadius={'5px'}
                      _hover={{
                        bgColor: 'rgba(255,128,0, 0.25)',
                      }}
                      color="white"
                      bgColor="#181818"
                      onClick={() => setSelectedPeriodOption('Voting Period')}
                    >
                      Voting Period
                    </MenuItem>
                    <MenuItem
                      borderRadius={'5px'}
                      _hover={{
                        bgColor: 'rgba(255,128,0, 0.25)',
                      }}
                      color="white"
                      bgColor="#181818"
                      onClick={() => setSelectedPeriodOption('Passed')}
                    >
                      Passed
                    </MenuItem>
                    <MenuItem
                      borderRadius={'5px'}
                      _hover={{
                        bgColor: 'rgba(255,128,0, 0.25)',
                      }}
                      color="white"
                      bgColor="#181818"
                      onClick={() => setSelectedPeriodOption('Rejected')}
                    >
                      Rejected
                    </MenuItem>
                  </MenuList>
                </>
              )}
            </Menu>

            <Menu>
              {({ isOpen }) => (
                <>
                  <MenuButton
                    _hover={{
                      bgColor: 'rgba(255,128,0, 0.25)',
                    }}
                    px={2}
                    color="white"
                    w="150px"
                    as={Button}
                    variant="outline"
                    rightIcon={<RotateIcon isOpen={isOpen} />}
                  >
                    {selectedProposalOption}
                  </MenuButton>

                  <MenuList borderColor="black" bgColor="#181818" minW="150px">
                    <MenuItem
                      borderRadius={'5px'}
                      _hover={{
                        bgColor: 'rgba(255,128,0, 0.25)',
                      }}
                      color="white"
                      bgColor="#181818"
                      onClick={() => setSelectedProposalOption('All Proposals')}
                    >
                      All Proposals
                    </MenuItem>
                    <MenuItem
                      borderRadius={'5px'}
                      _hover={{
                        bgColor: 'rgba(255,128,0, 0.25)',
                      }}
                      color="white"
                      bgColor="#181818"
                      onClick={() => setSelectedProposalOption('Voted')}
                    >
                      Voted
                    </MenuItem>
                  </MenuList>
                </>
              )}
            </Menu>
          </ButtonGroup>
        </Flex>
        <Box
          pr={2}
          maxHeight="2xl"
          overflowY="scroll"
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
        >
          {isLoading ? <Loader /> : content}
        </Box>
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
