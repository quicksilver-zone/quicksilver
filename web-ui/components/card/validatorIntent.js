import {
    Flex,
    Image,
    Text,
    Center,
    Box,
} from '@chakra-ui/react'
import stakingStyles from '@/styles/Staking.module.css'

const ValidatorIntentCard = (props) => {
    return (
        <Flex
            padding={'1em'}
            justify={'space-between'}
        >
            <Center gap={'10px'}>
                <Image src='/atom.svg' boxSize={'100%'} />
                <Text className={`${stakingStyles.switch_network_modal_sub_text}`}>
                    ATOM
                </Text>
            </Center>
            <Box>
                <Text className={`${stakingStyles.tableMainText}`}>
                    12.5 %
                </Text>
            </Box>
        </Flex>
    )
}

export default ValidatorIntentCard