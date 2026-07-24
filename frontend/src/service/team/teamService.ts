import apiClient from '@/service/apiClient';
import type {UserWithTeams, AddMemberPayload, TeamMember} from '@/types/team';
import type {UserSearchResult} from '@/types/user';

export const fetchUserTeams = async (): Promise<UserWithTeams> => {
    const response = await apiClient.get<UserWithTeams>('/api/teams/me');
    return response.data;
}

export const addTeamMember = async (teamId: string, payload: AddMemberPayload): Promise<void> => {
    await apiClient.post(`/api/teams/${teamId}/members`, payload)
}

export const createTeam = async (teamName: string): Promise<void> => {
    await apiClient.post('/api/teams', { team_name: teamName })
}

export const fetchTeamMembers = async (teamId: string): Promise<TeamMember[]> => {
    const response = await apiClient.get<TeamMember[]>(`/api/teams/${teamId}/members`);
    return response.data;
}

export const removeTeamMember = async (teamId: string, userId: string): Promise<void> => {
    await apiClient.delete(`/api/teams/${teamId}/members/${userId}`);
}

export const searchUsers = async (query: string): Promise<UserSearchResult[]> => {
    const response = await apiClient.get<UserSearchResult[]>('/api/users/search', {
        params: { q: query },
    });
    return response.data;
}