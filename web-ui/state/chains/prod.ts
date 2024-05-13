export const ibcDenomWithdrawMapping = {
    quicksilver: {
      qATOM: 'qatom',
      qOSMO: 'qosmo',
      qSTARS: 'qstars',
      qREGEN: 'qregen',
      qSOMM: 'qsomm',
      qJUNO: 'qjuno',
      qDYDX: 'qdydx',
      qSAGA: 'qsaga'
    }
  };

  export const ibcDenomDepositMapping = {
    osmosis: {
      qATOM: 'ibc/FA602364BEC305A696CBDF987058E99D8B479F0318E47314C49173E8838C5BAC',
      qOSMO: 'ibc/42D24879D4569CE6477B7E88206ADBFE47C222C6CAD51A54083E4A72594269FC',
      qSTARS: 'ibc/46C83BB054E12E189882B5284542DB605D94C99827E367C9192CF0579CD5BC83',
      qREGEN: 'ibc/79A676508A2ECA1021EDDC7BB9CF70CEEC9514C478DA526A5A8B3E78506C2206',
      qSOMM: 'ibc/EAF76AD1EEF7B16D167D87711FB26ABE881AC7D9F7E6D0CF313D5FA530417208',
      qJUNO: 'ibc/B4E18E61E1505C2F371B621E49B09E983F6A138F251A7B5286A6BDF739FD0D54',
      qDYDX: 'ibc/273C593E51ACE56F1F2BDB3E03A5CB81BB208B894BCAA642676A32C3454E8C27',
      qSAGA: 'ibc/F2D400F2728E9DA06EAE2AFAB289931A69EDDA5A661578C66A3177EDFE3C0D13'
    },
    umee: {
      qATOM: 'ibc/454725EA4029BAA99C293904336DE9A4B84E2BF7D83B9C56EE6B03E8A65FB5A1',
      qOSMO: 'ibc/F0D60708ACC09F2BDFF531D17477AE5F218220943A4792256DEF3F836E875D27',
      qSTARS: 'ibc/31946162F3E898B9E3A21792DD2AC740F2E82E7B92769BDF239C3DDA1726BB9F',
      qREGEN: 'ibc/16F0C7E49C2FE3A99E92A20DBCF4006B38ABC4E29F7F37829AD40F2C585BE835',
      qSOMM: 'ibc/ACF9DA139FE5BC8F95AC4A12B0B6D7710274DEDAC57284B881BEE1896F40642D',
      qJUNO: 'ibc/CA0BEF2524A37205009210EFCFB09585FBA9648C5F065FA078944A5C6704E8DC',
      qDYDX: 'ibc/41F3C94FAB3FB2D6D2B1F130A78697B07D729D1F50DA132C18F7963413A2DCF6',
      qSAGA: 'ibc/9B4BDA7382D0CF8C48A9D7496449D626DDF99AF640325978B5BD1AD4A9ED274C'
    },
  };

export const networks = [
    {
      value: 'ATOM',
      logo: '/img/networks/atom.svg',
      qlogo: '/img/networks/qatom.svg',
      name: 'Cosmos Hub',
      chainName: 'cosmoshub',
      chainId: 'cosmoshub-4',
    },
    {
      value: 'OSMO',
      logo: '/img/networks/osmo.svg',
      qlogo: '/img/networks/qosmo.svg',
      name: 'Osmosis',
      chainName: 'osmosis',
      chainId: 'osmosis-1',
    },
    {
      value: 'DYDX',
      logo: '/img/networks/dydx.svg',
      qlogo: '/img/networks/qdydx.svg',
      name: 'Dydx',
      chainName: 'dydx',
    chainId: 'dydx-mainnet-1',
    },
    {
      value: 'STARS',
      logo: '/img/networks/stargaze.svg',
      qlogo: '/img/networks/qstars.svg',
      name: 'Stargaze',
      chainName: 'stargaze',
      chainId: 'stargaze-1',
    },
    {
      value: 'REGEN',
      logo: '/img/networks/regen.svg',
      qlogo: '/img/networks/qregen.svg',
      name: 'Regen',
      chainName: 'regen',
      chainId: 'regen-1',
    },
    {
      value: 'SOMM',
      logo: '/img/networks/somm.svg',
      qlogo: '/img/networks/qsomm.svg',
      name: 'Sommelier',
      chainName: 'sommelier',
      chainId: 'sommelier-3',
    },
    {
      value: 'JUNO',
      logo: '/img/networks/juno.svg',
      qlogo: '/img/networks/qjuno.svg',
      name: 'Juno',
      chainName: 'juno',
      chainId: 'juno-1',
      },
    {
        value: 'SAGA',
        logo: '/img/networks/saga.svg',
        qlogo: '/img/networks/qsaga.svg',
        name: 'Saga',
        chainName: 'saga',
        chainId: 'ssc-1',
    },
    
  ];

  export const testNetworks = [
    {
      value: 'ATOM',
      logo: '/img/networks/atom.svg',
      qlogo: '/img/networks/qatom.svg',
      name: 'Cosmos Hub',
      chainName: 'cosmoshub',
      chainId: 'provider',
    },
    {
      value: 'OSMO',
      logo: '/img/networks/osmo.svg',
      qlogo: '/img/networks/qosmo.svg',
      name: 'Osmosis',
      chainName: 'osmosistestnet',
      chainId: 'osmo-test-5',
    },
    {
      value: 'STARS',
      logo: '/img/networks/stargaze.svg',
      qlogo: '/img/networks/qstars.svg',
      name: 'Stargaze',
      chainName: 'stargaze',
      chainId: 'elgafar-1',
    },
    {
      value: 'REGEN',
      logo: '/img/networks/regen.svg',
      qlogo: '/img/networks/qregen.svg',
      name: 'Regen',
      chainName: 'regen',
      chainId: 'regen-redwood-1',
    },
    {
      value: 'SOMM',
      logo: '/img/networks/somm.svg',
      qlogo: '/img/networks/qsomm.svg',
      name: 'Sommelier',
      chainName: 'Sommelier',
      chainId: 'sommelier-3',
    },
    {
        value: 'JUNO',
        logo: '/img/networks/juno.svg',
        qlogo: '/img/networks/qjuno.svg',
        name: 'Juno',
        chainName: 'juno',
        chainId: 'juno-1',
        },
  ];
