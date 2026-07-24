import { useState, useEffect, type FormEvent } from 'react';
import { addTeamMember, searchUsers } from '@/service/team/teamService';
import type { UserSearchResult } from '@/types/user';

export const useAddMemberModalState = (teamId: string, onClose: () => void) => {
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [searchResults, setSearchResults] = useState<UserSearchResult[]>([]);
  const [selectedUser, setSelectedUser] = useState<UserSearchResult | null>(null);
  const [isSearching, setIsSearching] = useState<boolean>(false);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  useEffect(() => {
    if (searchQuery.trim().length < 2) {
      setSearchResults([]);
      return;
    }

    const timer = setTimeout(async () => {
      setIsSearching(true);
      try {
        const results = await searchUsers(searchQuery.trim());
        setSearchResults(results);
      } catch (err: any) {
        console.error('Failed to search users:', err);
      } finally {
        setIsSearching(false);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [searchQuery]);

  const handleSelectUser = (user: UserSearchResult) => {
    setSelectedUser(user);
    setSearchQuery('');
    setSearchResults([]);
    setError(null);
  };

  const handleClearSelection = () => {
    setSelectedUser(null);
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!selectedUser) {
      setError('Please search and select a user to add');
      return;
    }

    setIsSubmitting(true);
    setError(null);
    setSuccessMsg(null);

    try {
      await addTeamMember(teamId, { user_id: selectedUser.id });
      setSuccessMsg(`Added member successfully`);
      setSelectedUser(null);
      setTimeout(() => onClose(), 1200);
    } catch (err: any) {
      setError(err?.response?.data?.error || 'failed to add team member');
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    searchQuery,
    setSearchQuery,
    searchResults,
    selectedUser,
    isSearching,
    isSubmitting,
    error,
    successMsg,
    handleSelectUser,
    handleClearSelection,
    handleSubmit,
  };
};