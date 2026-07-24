import { X, UserPlus, Loader2, Search, Check } from 'lucide-react';
import { useAddMemberModalState } from './useAddMemberModalState';

interface AddMemberModalProps {
  teamId: string;
  teamName: string;
  onClose: () => void;
}

export const AddMemberModal = ({ teamId, teamName, onClose }: AddMemberModalProps) => {
  const {
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
  } = useAddMemberModalState(teamId, onClose);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div className="bg-card border border-border rounded-[20px] p-6 w-full max-w-md shadow-2xl relative">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-muted-foreground hover:text-foreground p-1 transition-colors cursor-pointer"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="flex items-center gap-3 mb-4">
          <div className="p-2 rounded-xl bg-emerald-500/10 text-emerald-500">
            <UserPlus className="w-5 h-5" />
          </div>
          <div>
            <h2 className="text-lg font-bold text-foreground">Add Member</h2>
            <p className="text-xs text-muted-foreground">Team: {teamName}</p>
          </div>
        </div>

        {error && (
          <div className="mb-4 p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-500 text-xs font-medium">
            {error}
          </div>
        )}

        {successMsg && (
          <div className="mb-4 p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 text-xs font-medium">
            {successMsg}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-medium text-muted-foreground mb-1">
              Search User (by Email or Name)
            </label>

            {selectedUser ? (
              <div className="flex items-center justify-between p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/30">
                <div className="flex items-center gap-2">
                  <Check className="w-4 h-4 text-emerald-500" />
                  <div>
                    <p className="text-xs font-medium text-foreground">
                      {selectedUser.display_name || selectedUser.email}
                    </p>
                    <p className="text-[11px] text-muted-foreground">{selectedUser.email}</p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={handleClearSelection}
                  className="text-muted-foreground hover:text-foreground p-1 cursor-pointer"
                >
                  <X className="w-4 h-4" />
                </button>
              </div>
            ) : (
              <div className="relative">
                <div className="relative flex items-center">
                  <Search className="w-4 h-4 absolute left-3 text-muted-foreground pointer-events-none" />
                  <input
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="Search by email or name..."
                    className="w-full bg-muted border border-border rounded-xl pl-9 pr-9 py-2.5 text-sm text-foreground focus:outline-none focus:border-emerald-500 transition-colors"
                  />
                  {isSearching && (
                    <Loader2 className="w-4 h-4 absolute right-3 text-emerald-500 animate-spin" />
                  )}
                </div>

                {/* Dropdown Results */}
                {searchResults.length > 0 && (
                  <div className="absolute left-0 right-0 mt-1.5 max-h-48 overflow-y-auto rounded-xl border border-border bg-card shadow-xl z-50 py-1">
                    {searchResults.map((user) => (
                      <button
                        key={user.id}
                        type="button"
                        onClick={() => handleSelectUser(user)}
                        className="w-full text-left px-3 py-2 hover:bg-muted transition-colors flex flex-col cursor-pointer"
                      >
                        <span className="text-xs font-medium text-foreground">
                          {user.display_name || user.email}
                        </span>
                        <span className="text-[11px] text-muted-foreground">{user.email}</span>
                      </button>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !selectedUser}
              className="flex items-center gap-2 bg-emerald-500 hover:bg-emerald-600 text-white font-medium px-5 py-2 rounded-xl text-sm transition-colors disabled:opacity-50 cursor-pointer"
            >
              {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Add Member'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};