import {
    Flex,
    Image,
    Text,
    Center,
    Heading,
    Box,
} from '@chakra-ui/react'
import stakingStyles from '@/styles/Staking.module.css'

const NetworkCard = (props) => {
    return (
        <Flex
            background='linear-gradient(93deg, rgba(0, 0, 0, 0.39) 41.57%, rgba(0, 0, 0, 0.39) 102.36%, rgba(165, 162, 162, 0.29) 105.28%, rgba(120, 117, 117, 0.40) 108.95%, rgba(231, 227, 227, 0.14) 110.37%, rgba(171, 171, 171, 0.09) 123.66%)'
            borderRadius={'8px'}
            border='1px solid var(--neutral-stroke, rgba(255, 255, 255, 0.20))'
            padding={'1em 2em'}
            justify={'space-between'}
        >
            <Center gap={'10px'}>
                <Image src='/atom.svg' boxSize={'100%'}/>
                <Box>
                    <Heading as='h6' size={'md'}>  
                        Cosmoshub
                    </Heading>
                    <Text className={`${stakingStyles.switch_network_modal_sub_text}`}>
                        ATOM
                    </Text>
                </Box>
            </Center>
            <Box>
                <Heading as='h6' size={'md'}>
                    25%
                </Heading>
                <Text className={`${stakingStyles.switch_network_modal_sub_text}`}>
                    APY
                </Text>
            </Box>
        </Flex>
    )
}

export default NetworkCard