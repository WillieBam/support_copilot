import { useState, useEffect, useCallback } from 'react';
import { fetchUserTeams } from '@/service/team/teamService';
import type { UserMembership, TeamRole } from '@/types/team';

export const useTeamState = (isSignedIn: boolean) => {
  const [memberships, setMemberships] = useState<UserMembership[]>([]);
  const [activeTeamId, setActiveTeamId] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // load user team memberships
  const loadTeams = useCallback(async () => {
    if (!isSignedIn) return;
    setIsLoading(true);
    setError(null);
    try {
      const data = await fetchUserTeams();
      const userMemberships = data.memberships || [];
      setMemberships(userMemberships);
      
      // default active team to first joined team if not selected
      if (userMemberships.length > 0 && !activeTeamId) {
        setActiveTeamId(userMemberships[0].team_id);
      }
    } catch (err: any) {
      setError(err?.response?.data?.error || 'failed to load user teams');
    } finally {
      setIsLoading(false);
    }
  }, [isSignedIn, activeTeamId]);

  useEffect(() => {
    loadTeams();
  }, [loadTeams]);

  const activeMembership = memberships.find((m) => m.team_id === activeTeamId) || null;
  const activeRole: TeamRole | null = activeMembership ? activeMembership.role : null;
  const isOwner = activeRole === 'owner';

  const selectTeam = (teamId: string) => {
    setActiveTeamId(teamId);
  };

  return {
    memberships,
    activeTeamId,
    activeMembership,
    activeRole,
    isOwner,
    isLoading,
    error,
    selectTeam,
    reloadTeams: loadTeams,
  };
};