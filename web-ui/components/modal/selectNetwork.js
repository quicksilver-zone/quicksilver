import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalBody,
    ModalCloseButton,
    Flex,
    Image,
    Text,
    Center,
    Link,
    Box,
    Grid,
} from '@chakra-ui/react'
import stakingStyles from '@/styles/Staking.module.css'
import NetworkCard from '../card/network'

const SelectNetwork = (props) => {
    return (
        <Modal
            isOpen={props.isShow}
            onClose={() => props.setIsShow(false)}
            size='4xl'
        >
            <ModalOverlay
                bg='blackAlpha.300'
                backdropFilter='blur(10px)'
            />
            <ModalContent
                color='#FBFBFB'
                background={'linear-gradient(207deg, rgba(255, 255, 255, 0.09) -46.03%, rgba(255, 255, 255, 0.00) 128.13%), #0E0E0E'}
                borderRadius={'20px'}
            >
                <ModalHeader>
                    <Center justifyContent={'start'} gap={'10px'}>
                        <Image src='/logo/qs_logo.svg' />
                        <Text fontSize={'24px'}>
                            Switch Network
                        </Text>
                    </Center>
                    <ModalCloseButton
                        color='#E77728'
                        boxSize={'3em'}
                        fontSize={'1em'}
                        _hover={{
                            backgroundColor: 'transparent'
                        }}
                    />
                </ModalHeader>
                <ModalBody padding={'0'}>
                    <Box padding={'0 2em 1em 2em'}>
                        <Text className={`${stakingStyles.modal_m_size}`}>
                            Stake assets with your preferred validators and earn staking yield while receiving a qAsset for use in DeFi.
                            Select a network to get started Liquid Staking.
                        </Text>
                        <Text className={`${stakingStyles.modal_m_size}`}>
                            By selecting a network, you confirm you accept the Quicksilver
                            Terms of Service.
                        </Text>
                        <br />
                        <Text className={`${stakingStyles.modal_m_size}`}>
                            Don`t see your preferred network? <Link isExternal color='#7198FA'>Request a network</Link>
                        </Text>
                    </Box>
                    <Box backgroundColor={'rgba(14, 14, 14, 0.6)'} borderRadius={'0 0 20px 20px'}>
                        <Grid 
                            templateColumns={'repeat(2, 1fr)'} 
                            maxHeight={'325px'} 
                            rowGap={5} 
                            columnGap={10} 
                            overflowY={'scroll'}
                            padding={'2em 2em 2em 2em'}
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
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                            <NetworkCard />
                        </Grid>
                    </Box>
                </ModalBody>
            </ModalContent>
        </Modal>
    )
}

export default SelectNetwork