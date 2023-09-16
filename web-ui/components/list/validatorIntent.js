import {
    Flex,
    Text,
    Center,
    Button,
    VStack,
    StackDivider,
} from "@chakra-ui/react"
import ValidatorIntentCard from "../card/validatorIntent"

const ValidatorIntent = () => {
    return (
        <>
            <Flex justify={'space-between'} padding={'1em 2em'}>
                <Center>
                    <Text fontSize={'20px'}>
                        Validator List
                    </Text>
                </Center>
                <Button
                    variant={'ghost'}
                    color={'rgba(255, 133, 0, 1)'}
                    _hover={{ backgroundColor: 'transparent' }}
                    fontWeight={'400'}
                    fontSize={'14px'}
                >
                    Edit Stake Allocation
                </Button>
            </Flex>
            <VStack
                margin={'0 2em 1em 2em'}
                divider={<StackDivider height={'1px'} borderColor='rgba(255, 255, 255, 0.20)' />}
                align='stretch'
                borderRadius={'10px'}
                border='0.5px solid var(--neutral-stroke, rgba(255, 255, 255, 0.20))'
                backgroundColor={'#171717'}
                maxHeight={'325px'}
                gap={0}
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
                <ValidatorIntentCard />
                <ValidatorIntentCard />
                <ValidatorIntentCard />
                <ValidatorIntentCard />
                <ValidatorIntentCard />
                <ValidatorIntentCard />

            </VStack>
        </>
    )
}

export default ValidatorIntent