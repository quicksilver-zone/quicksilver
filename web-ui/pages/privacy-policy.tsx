import { VStack, Heading, Text, SlideFade, Container, Divider } from '@chakra-ui/react';
import Head from 'next/head';

const PrivacyPolicyPage = () => {
  return (
    <SlideFade offsetY={'200px'} in={true}>
      <Container
        flexDir={'column'}
        top={20}
        mt={{ base: 10, md: 0 }}
        zIndex={2}
        position="relative"
        justifyContent="center"
        alignItems="center"
        maxW="6xl"
      >
        <Head>
          <title>Privacy Policy </title>
          <meta name="viewport" content="width=device-width, initial-scale=1.0" />
          <link rel="icon" href="/quicksilver/img/favicon.png" />
        </Head>
        <VStack pb={12} spacing={4} align="stretch" m={8}>
          <Heading color="white" as="h1" size="xl" textAlign="left">
            Privacy Policy
          </Heading>
          <Text fontSize="md" color="gray.200">
            Last updated: January 8th, 2024
          </Text>
          <Divider />
          <Text fontSize="md" color="gray.200" mt={4}>
            Thank you for choosing to be part of our community at QuickSilver. We are committed to protecting your personal information and
            your right to privacy. If you have any questions or concerns about this privacy notice, or our practices with regards to your
            personal information, please contact us at [Contact Information].
          </Text>
          <Heading as="h2" size="lg" mt={4}>
            Information We Collect
          </Heading>
          <Text fontSize="md" color="gray.200">
            When you visit our website (https://app.quicksilver.zone), and more generally, use any of our services (the
            &ldquo;Services&rdquo;, which include the Website), we appreciate that you are trusting us with your personal information. We
            take your privacy very seriously. In this privacy notice, we seek to explain to you in the clearest way possible what
            information we collect, how we use it, and what rights you have in relation to it.
          </Text>
          <Heading as="h2" size="lg" mt={4}>
            Cookies and Tracking Technologies
          </Heading>
          <Text fontSize="md" color="gray.200">
            We use cookies and similar tracking technologies (like web beacons and pixels) to access or store user interactions.
          </Text>
          <Heading as="h2" size="lg" mt={4}>
            Use of Your Information
          </Heading>
          <Text fontSize="md" color="gray.200">
            We use personal wallet information (NOT PRIVATE KEYS) collected via our Website for a variety of business purposes described
            below. We process your wallet addres for these purposes in reliance on our legitimate business interests, in order to provide
            liquid staking services to you.
          </Text>
          <Text fontSize="md" color="gray.200">
            - To facilitate account creation and authentication and otherwise manage user accounts. - To send administrative information to
            you, such as changes to our terms, conditions, and policies. - To enforce our terms, conditions, and policies for business
            purposes, to comply with legal and regulatory requirements, or in connection with our contract. - To respond to legal requests
            and prevent harm. - For other business purposes, such as data analysis, identifying usage trends, determining the effectiveness
            of our promotional campaigns, and to evaluate and improve our Website, products, marketing, and your experience.
          </Text>
          <Heading as="h2" size="lg" mt={4}>
            Information Sharing and Disclosure
          </Heading>
          <Text fontSize="md" color="gray.200">
            We may process or share your data that we hold based on the following legal basis:
          </Text>
          <Text fontSize="md" color="gray.200">
            - Consent: We may process your data if you have given us specific consent to use your personal information for a specific
            purpose. - Legitimate Interests: We may process your data when it is reasonably necessary to achieve our legitimate business
            interests. - Performance of a Contract: Where we have entered into a contract with you, we may process your personal information
            to fulfill the terms of our contract. - Legal Obligations: We may disclose your information where we are legally required to do
            so in order to comply with applicable law, governmental requests, a judicial proceeding, court order, or legal process, such as
            in response to a court order or a subpoena (including in response to public authorities to meet national security or law
            enforcement requirements). - Vital Interests: We may disclose your information where we believe it is necessary to investigate,
            prevent, or take action regarding potential violations of our policies, suspected fraud, situations involving potential threats
            to the safety of any person, and illegal activities, or as evidence in litigation in which we are involved.
          </Text>
          <Heading as="h2" size="lg" mt={4}>
            Data Security
          </Heading>
          <Text fontSize="md" color="gray.200">
            We have implemented appropriate technical and organizational security measures designed to protect the security of any personal
            information we process. However, despite our safeguards and efforts to secure your information, no electronic transmission over
            the Internet or information storage technology can be guaranteed to be 100% secure, so we cannot promise or guarantee that
            hackers, cybercriminals, or other unauthorized third parties will not be able to defeat our security and improperly collect,
            access, steal, or modify your information.
          </Text>
          {/* Add more paragraphs as needed */}
        </VStack>
      </Container>
    </SlideFade>
  );
};

export default PrivacyPolicyPage;
