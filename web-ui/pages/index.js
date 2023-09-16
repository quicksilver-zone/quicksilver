import Image from 'next/image'
import { Inter } from 'next/font/google'
import Layout from '@/components/layout'
import Link from 'next/link'

import styles from '@/styles/Home.module.css'
import PageHead from '@/components/layout/PageHead'

const inter = Inter({ subsets: ['latin'] })

export default function Home() {
  return (
    <>
      <PageHead pageTitle="Quicksilver DApp" />
      <div>heloo</div>
    </>
  )
}
