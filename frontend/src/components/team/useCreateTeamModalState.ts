import { useState, type FormEvent } from 'react';
import { createTeam } from '@/service/team/teamService';

interface UseCreateTeamModalStateProps {
  reloadTeams: () => Promise<void>;
  onClose: () => void;
}

export const useCreateTeamModalState = ({ reloadTeams, onClose }: UseCreateTeamModalStateProps) => {
  const [teamName, setTeamName] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!teamName.trim()) {
      setError('Team name cannot be empty');
      return;
    }
    if (teamName.length > 20) {
      setError('Team name cannot exceed 20 characters');
      return;
    }
    
    setError(null);
    setIsSubmitting(true);
    
    try {
      await createTeam(teamName);
      setSuccessMsg('Team created successfully!');
      await reloadTeams();
      setTimeout(() => {
        onClose();
      }, 1200);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create team');
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    teamName,
    setTeamName,
    isSubmitting,
    error,
    successMsg,
    handleSubmit,
  };
};
