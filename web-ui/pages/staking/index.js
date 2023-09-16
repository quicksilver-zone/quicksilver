import { Inter } from 'next/font/google'
import { useState } from 'react'
import {
    Flex,
} from '@chakra-ui/react'
import SideBar from "@/components/layout/Sidebar";
import Layout from "@/components/layout";
import StakingPannel from '@/components/staking/staking';
import ValidatorPanel from '@/components/staking/validator';
import PageHead from '@/components/layout/PageHead';

export default function Staking() {
    const [balances, setBalances] = useState([])
    const [step, setStep] = useState(1)

    return (
        <>
            <PageHead pageTitle="Staking | Quicksilver" />
                {
                    step === 1 ? <StakingPannel setStep={setStep} />
                        : <ValidatorPanel setStep={setStep} />
                }
        </>
    )
}
