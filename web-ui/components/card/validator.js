import { Inter } from 'next/font/google'
import stakingStyles from '@/styles/Staking.module.css'
import { useState } from 'react'
import {
    Button, Grid,
    Center,
    Flex,
    InputGroup,
    InputLeftElement,
    Input,
    Switch,
    Box,
    Checkbox,
    IconButton,
    Icon,
    Image,
    Text,
} from '@chakra-ui/react'
import { ChevronLeftIcon, SearchIcon } from '@chakra-ui/icons'
import { AiFillStar, AiOutlineStar } from "react-icons/ai";


const ValidatorCard = (props) => {
    const [isStar, setIsStar] = useState(false)
    return (
        <Grid
            templateColumns='40% 15% 15% 15% 15%'
            borderBottom={'solid 1px rgba(77, 77, 77, 1)'}
        >
            <Box
                padding={'1em 0'}
                borderRight={'solid 2px #4D4D4D'}
                w={'90%'}
            >
                <Flex>
                    <Center gap={'5px'} fontSize='15px'>
                        <Checkbox colorScheme='orange' borderColor={'#E77728'} />
                        <IconButton
                            variant='ghost'
                            colorScheme='teal'
                            aria-label='Done'
                            fontSize='20px'
                            icon={isStar ? <Icon as={AiFillStar} /> : <Icon as={AiOutlineStar} />}
                            onClick={() => setIsStar(!isStar)}
                            _hover={{
                                backgroundClip: 'transparent'
                            }}
                        />
                        <text>
                            {props.index}
                        </text>
                        <Image src='/logo/qs_logo.svg' boxSize={'30px'} borderRadius={'50%'} />
                        <text>
                            {props.name}
                        </text>
                    </Center>
                </Flex>
            </Box>
            <Box
                padding={'1em 0'}
            >
                <Text className={`${stakingStyles.tableMainText}`}>
                    {props.votingPower}
                </Text>
                <Text className={`${stakingStyles.tableSubText}`}>
                    {props.votingPowerPercentage}
                </Text>
            </Box>
            <Flex justify={'start'}>
                <Center
                    padding={'1em 0'}
                >
                    <Text className={`${stakingStyles.tableMainText}`}>
                        {props.commission}
                    </Text>
                </Center>
            </Flex>
            <Flex justify={'start'} padding={'1em 0'}>
                <Center borderRadius={'10px'} backgroundColor={'#9E9E9E'} w={'70%'} padding={'.5em'}>
                    <Text className={`${stakingStyles.tableMainText}`}>
                        {props.record}
                    </Text>
                </Center>
            </Flex>
            <Flex justify={'start'} padding={'1em 0'}>
                <Center borderRadius={'10px'} backgroundColor={'#1E421E'} w={'70%'} padding={'.5em'}>
                    <Text className={`${stakingStyles.tableMainText}`}>
                        {`Level 0${props.prScore}`}
                    </Text>
                </Center>
            </Flex>
        </Grid>
    )
}

export default ValidatorCard