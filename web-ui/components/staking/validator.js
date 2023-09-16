import { Inter } from 'next/font/google'
import stakingStyles from '@/styles/Staking.module.css'
import { useState } from 'react'
import {
    Button, Grid,
    Center,
    Flex,
    InputGroup,
    InputLeftElement,
    Text,
    Input,
    Switch,
    Box,
    Icon,
    HStack,
    VStack
} from '@chakra-ui/react'
import { ChevronLeftIcon, SearchIcon } from '@chakra-ui/icons'
import ValidatorCard from '@/components/card/validator'
import { GoPencil } from "react-icons/go";
import StakingModal from '../modal/staking'

const ValidatorPanel = (props) => {
    const [pannelMode, setPannelMode] = useState(0)
    const [show, setShow] = useState(false)
    return (
        <Center w={'100%'} padding={'10vh'} minH={'100%'}>
            <Box className={`${stakingStyles.staking_container}`} w={'100%'} h={'100%'}>
                <Flex direction={'column'} justify={'space-between'} h={'100%'}>
                    <Box h={'100%'}>
                        <Flex justify={'space-between'}>
                            <Flex justify={'start'}>
                                <Button variant={'ghost'} padding={0} _hover={{ backgroundColor: 'Boxansparent' }} onClick={() => props.setStep(1)}>
                                    <ChevronLeftIcon boxSize={'100%'} color={'#E77728'} />
                                </Button>
                                <button
                                    className={`${stakingStyles.pannel_mode_btn} ${pannelMode === 0 && stakingStyles.in_color}`}
                                    id={`${stakingStyles.stake_btn}`}
                                    onClick={() => setPannelMode(0)}
                                >
                                    All Validators
                                </button>
                                <Center>
                                    <div className={`${stakingStyles.verticalLine}`} style={{ height: '80%' }} />
                                </Center>
                                <button
                                    className={`${stakingStyles.pannel_mode_btn} ${pannelMode === 1 && stakingStyles.in_color}`}
                                    id={`${stakingStyles.unstake_btn}`}
                                    onClick={() => setPannelMode(1)}
                                >
                                    Favourites
                                </button>
                            </Flex>
                            <InputGroup w={'30%'}>
                                <InputLeftElement pointerEvents='none'>
                                    <SearchIcon color='gray.300' />
                                </InputLeftElement>
                                <Input
                                    type='tel'
                                    placeholder='Search validator'
                                    borderRadius={'20px'}
                                    _focus={{ borderColor: '#E77728', boxShadow: 'none' }}
                                />
                            </InputGroup>
                            <Center gap={'10px'} >
                                <Switch size={'lg'} colorScheme='orange' />
                                <Text>
                                    Show inactive validators
                                </Text>
                            </Center>
                        </Flex>
                        <Box margin={'3em 0'} h={'100%'}>
                            <Grid templateColumns='40% 15% 15% 15% 15%' padding={'0 1em'} fontWeight='bold' fontSize={'14px'}>
                                <Box>VALIDATOR</Box>
                                <Box>VOTING POWER</Box>
                                <Box>COMMISSION</Box>
                                <Box>VOTING RECORD</Box>
                                <Box>PR SCORE</Box>
                            </Grid>
                            <Box
                                bg={'rgba(20, 20, 20, 1)'}
                                borderRadius='10px'
                                padding={'0 1em'}
                                boxShadow={'0px 0px 5px 0px rgba(255, 255, 255, 0.50)'}
                                h={'100%'}
                                maxH={'calc(48vh - 2em)'}
                                overflowY={'scroll'}
                                sx={{
                                    '&::-webkit-scrollbar': {
                                        width: '5px',
                                    },
                                    '&::-webkit-scrollbar-thumb': {
                                        background: '#E77728',
                                        borderRadius: '20px',
                                    },
                                }}
                            >
                                <ValidatorCard
                                    index={1}
                                    name={'Lavender.Five Nodes'}
                                    votingPower={'12,793,452'}
                                    votingPowerPercentage={`${'6.18'} %`}
                                    commission={`${'22.35'} %`}
                                    record={'12/65'}
                                    prScore={1}
                                />
                                <ValidatorCard
                                    index={1}
                                    name={'Lavender.Five Nodes'}
                                    votingPower={'12,793,452'}
                                    votingPowerPercentage={`${'6.18'} %`}
                                    commission={`${'22.35'} %`}
                                    record={'12/65'}
                                    prScore={1}
                                />
                                <ValidatorCard
                                    index={1}
                                    name={'Lavender.Five Nodes'}
                                    votingPower={'12,793,452'}
                                    votingPowerPercentage={`${'6.18'} %`}
                                    commission={`${'22.35'} %`}
                                    record={'12/65'}
                                    prScore={1}
                                />
                                <ValidatorCard
                                    index={1}
                                    name={'Lavender.Five Nodes'}
                                    votingPower={'12,793,452'}
                                    votingPowerPercentage={`${'6.18'} %`}
                                    commission={`${'22.35'} %`}
                                    record={'12/65'}
                                    prScore={1}
                                />
                                <ValidatorCard
                                    index={1}
                                    name={'Lavender.Five Nodes'}
                                    votingPower={'12,793,452'}
                                    votingPowerPercentage={`${'6.18'} %`}
                                    commission={`${'22.35'} %`}
                                    record={'12/65'}
                                    prScore={1}
                                />
                                <ValidatorCard
                                    index={1}
                                    name={'Lavender.Five Nodes'}
                                    votingPower={'12,793,452'}
                                    votingPowerPercentage={`${'6.18'} %`}
                                    commission={`${'22.35'} %`}
                                    record={'12/65'}
                                    prScore={1}
                                />
                                <ValidatorCard
                                    index={1}
                                    name={'Lavender.Five Nodes'}
                                    votingPower={'12,793,452'}
                                    votingPowerPercentage={`${'6.18'} %`}
                                    commission={`${'22.35'} %`}
                                    record={'12/65'}
                                    prScore={1}
                                />
                                <ValidatorCard
                                    index={1}
                                    name={'Lavender.Five Nodes'}
                                    votingPower={'12,793,452'}
                                    votingPowerPercentage={`${'6.18'} %`}
                                    commission={`${'22.35'} %`}
                                    record={'12/65'}
                                    prScore={1}
                                />
                                <ValidatorCard
                                    index={1}
                                    name={'Lavender.Five Nodes'}
                                    votingPower={'12,793,452'}
                                    votingPowerPercentage={`${'6.18'} %`}
                                    commission={`${'22.35'} %`}
                                    record={'12/65'}
                                    prScore={1}
                                />
                            </Box>
                        </Box>
                    </Box>
                    <Center w={'100%'}>
                    <Box
                        bg={'rgba(20, 20, 20, 1)'}
                        borderRadius='10px'
                        padding={'2em 1em'}
                        boxShadow={'0px 0px 5px 0px rgba(255, 255, 255, 0.50)'}
                        h={'100%'}
                        w={'100%'}
                    >
                        <Flex justify={'space-between'}>
                            <Center gap={'10px'}>
                                <Text className={`${stakingStyles.tableMainText}`}>
                                    10.123123 ATOM
                                </Text>
                                <Center>
                                    <Icon as={GoPencil} color={'#E77728'} boxSize={'19px'} />
                                </Center>
                                <div className={`${stakingStyles.verticalLine}`} style={{ height: '100%', margin: '0 10px' }} />
                                <Box>
                                    <Text className={`${stakingStyles.tableMainText}`}>
                                        0 Validators Selected
                                    </Text>
                                    <Text className={`${stakingStyles.tableSubText}`}>
                                        Select between 1 to 8 validators.
                                    </Text>
                                </Box>
                            </Center>
                            <Button
                                w={'20%'}
                                padding={'1.5em 0'}
                                color={'black'}
                                backgroundColor={'rgba(255, 133, 0, 1)'}
                                _hover={{
                                    backgroundColor: '#ba5c1a'
                                }}
                                onClick={() => setShow(true)}
                            >
                                Next
                            </Button>
                        </Flex>
                    </Box>
                    </Center>
                </Flex>
            </Box>

            <StakingModal
                isShow={show}
                setIsShow={setShow}
            />
        </Center>
    )
}

export default ValidatorPanel