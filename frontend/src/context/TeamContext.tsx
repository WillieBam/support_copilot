import { createContext, useContext, type ReactNode } from 'react';
import { useTeamState } from '@/hooks/useTeamState';

type TeamContextType = ReturnType<typeof useTeamState>;

const TeamContext = createContext<TeamContextType | null>(null);

export const TeamProvider = ({ children, isSignedIn }: { children: ReactNode; isSignedIn: boolean }) => {
  const teamState = useTeamState(isSignedIn);
  return <TeamContext.Provider value={teamState}>{children}</TeamContext.Provider>;
};

export const useTeam = () => {
  const context = useContext(TeamContext);
  if (!context) {
    throw new Error('useTeam must be used within TeamProvider');
  }
  return context;
};