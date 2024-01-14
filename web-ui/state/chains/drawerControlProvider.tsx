import { createContext, ReactNode, useContext } from 'react';

const DrawerControlContext = createContext({
  closeDrawer: () => {},
});

// Hook to use the context
export const useDrawerControl = () => useContext(DrawerControlContext);

interface DrawerControlProviderProps {
  children: ReactNode;
  closeDrawer: () => void;
}

export const DrawerControlProvider = ({ children, closeDrawer }: DrawerControlProviderProps) => {
  return <DrawerControlContext.Provider value={{ closeDrawer }}>{children}</DrawerControlContext.Provider>;
};
