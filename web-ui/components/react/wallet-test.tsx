import { Button } from '@chakra-ui/react';
import { useChain, useWalletClient } from '@cosmos-kit/react';
import { useEffect } from 'react';

export default function WalletTest() {
  const { openView } = useChain('quicksilver');
  const { status, client } = useWalletClient();

  useEffect(() => {
    if (status === 'Done') {
      client?.enable?.(['cosmoshub-4', 'osmosis-1', 'regen-1', 'sommelier-3', 'stargaze-1']);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status]);

  return (
    <div style={{ textAlign: 'center', margin: '4rem' }}>
      <Button onClick={openView}>Connect</Button>
    </div>
  );
}
