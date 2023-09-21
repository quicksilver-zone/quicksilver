import { Box, Container } from '@chakra-ui/react';
import Head from 'next/head';

import { Header } from '@/components';
import { SideHeader } from '@/components';

export default function Home() {
  return (
    <>
      <Box
        w="100vw"
        h="100vh"
        bgImage="url('https://s3-alpha-sig.figma.com/quicksilver-app-v2/img/555d/db64/f5bf65e93a15603069e8e865d5f6d60d?Expires=1694995200&Signature=fYfmbqDdOGRYtSeEsOkavPhhkaNQK1UFFfICaUaM1k9OVEpACsoWOcK2upjRW7Tfs-pPTJBuQuvcmF9gBjosh5-Al2xTWHYzDlR~CYJNzsXcseIEnVf7H8lCdJqhZY-T0r~lmbJK5-CmbulWfOaubc-wyY3C-oM3b1RanGV1TqmPZto5bbHwf56jDYqK86HedVMXbUCOlzkeBw2R93AkmNDMOdDbKa9rIKqxil64DuQQAfIFxWm1Rc69Jc1-4K-bunsS~kfz8bSET6TIGmR15nCo~ibfISG72YYKAa7zz6XqUY6GKmmG-Yhj9XyyYb7Jy02r5axNei3DRD78SBe~6w__&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4')"
        bgSize="fit"
        bgPosition="right center"
        bgAttachment="fixed"
        bgRepeat="no-repeat"
        bgColor="#000000"
      >
        <Header chainName="quicksilver" />
        <SideHeader />
        <Container justifyContent="center" alignItems="center" maxW="5xl">
          <Head>
            <title>DeFi</title>
            <meta
              name="viewport"
              content="width=device-width, initial-scale=1.0"
            />
            <link rel="icon" href="/quicksilver-app-v2/img/favicon.png" />
          </Head>
        </Container>
      </Box>
    </>
  );
}
