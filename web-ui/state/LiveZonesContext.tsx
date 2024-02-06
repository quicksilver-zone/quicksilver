import axios from 'axios';
import { Zone } from 'quicksilverjs/dist/codegen/quicksilver/interchainstaking/v1/interchainstaking';
import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';

interface LiveZonesContextType {
  liveNetworks: string[];
}

interface LiveZonesProviderProps {
  children: ReactNode;
}

const LiveZonesContext = createContext<LiveZonesContextType>({ liveNetworks: [] });

export const useLiveZones = (): LiveZonesContextType => useContext(LiveZonesContext);

export const LiveZonesProvider: React.FC<LiveZonesProviderProps> = ({ children }) => {
  const [liveNetworks, setLiveNetworks] = useState([]);

  useEffect(() => {
    const fetchLiveZones = async () => {
      try {
        const response = await axios.get(`${process.env.NEXT_PUBLIC_QUICKSILVER_API}/quicksilver/interchainstaking/v1/zones`);
        const liveZones = response.data.zones.map((zone: Zone) => zone.chain_id);
        setLiveNetworks(liveZones);
      } catch (error) {
        console.error('Failed to fetch live zones:', error);
      }
    };

    fetchLiveZones();
  }, []);

  return <LiveZonesContext.Provider value={{ liveNetworks }}>{children}</LiveZonesContext.Provider>;
};
