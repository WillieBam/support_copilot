import { useState } from 'react';
import { useTeam } from '@/context/TeamContext';

export const useTeamSelectorState = () => {
  const { memberships, activeMembership, activeTeamId, isOwner, selectTeam, reloadTeams } = useTeam();
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [isAddModalOpen, setIsAddModalOpen] = useState<boolean>(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState<boolean>(false);

  const toggleDropdown = () => setIsOpen((prev) => !prev);
  const closeDropdown = () => setIsOpen(false);

  const handleSelectTeam = (teamId: string) => {
    selectTeam(teamId);
    setIsOpen(false);
  };

  const openAddModal = () => {
    setIsOpen(false);
    setIsAddModalOpen(true);
  };

  const closeAddModal = () => setIsAddModalOpen(false);

  const openCreateModal = () => {
    setIsOpen(false);
    setIsCreateModalOpen(true);
  };

  const closeCreateModal = () => setIsCreateModalOpen(false);

  return {
    memberships,
    activeMembership,
    activeTeamId,
    isOwner,
    isOpen,
    isAddModalOpen,
    isCreateModalOpen,
    toggleDropdown,
    closeDropdown,
    handleSelectTeam,
    openAddModal,
    closeAddModal,
    openCreateModal,
    closeCreateModal,
    reloadTeams,
  };
};
