import { Users, Loader2, X } from 'lucide-react';
import { useCreateTeamModalState } from './useCreateTeamModalState';

interface CreateTeamModalProps {
  reloadTeams: () => Promise<void>;
  onClose: () => void;
}

export const CreateTeamModal = ({ reloadTeams, onClose }: CreateTeamModalProps) => {
  const {
    teamName,
    setTeamName,
    isSubmitting,
    error,
    successMsg,
    handleSubmit,
  } = useCreateTeamModalState({ reloadTeams, onClose });

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center bg-background/80 backdrop-blur-sm p-4">
      <div className="relative w-full max-w-md bg-card border border-border rounded-[20px] shadow-2xl p-6">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="mb-6">
          <div className="w-12 h-12 bg-emerald-500/10 rounded-2xl flex items-center justify-center mb-4">
            <Users className="w-6 h-6 text-emerald-500" />
          </div>
          <h2 className="text-xl font-semibold">Create Team</h2>
          <p className="text-sm text-muted-foreground mt-1">
            Create a new team to collaborate with others.
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <div className="flex items-center justify-between mb-1">
              <label className="text-sm font-medium">Team Name</label>
              <span className="text-xs text-muted-foreground">
                {teamName.length}/20
              </span>
            </div>
            <input
              type="text"
              value={teamName}
              onChange={(e) => setTeamName(e.target.value)}
              placeholder="e.g. Support Geniuses"
              maxLength={20}
              className="w-full bg-background border border-border rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500/50"
              disabled={isSubmitting || !!successMsg}
            />
          </div>

          {error && (
            <div className="p-3 bg-red-500/10 border border-red-500/20 text-red-500 rounded-xl text-sm">
              {error}
            </div>
          )}

          {successMsg && (
            <div className="p-3 bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 rounded-xl text-sm">
              {successMsg}
            </div>
          )}

          <div className="flex justify-end gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              disabled={isSubmitting || !!successMsg}
              className="px-4 py-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !!successMsg || !teamName.trim()}
              className="flex items-center justify-center min-w-[120px] bg-emerald-500 hover:bg-emerald-600 text-white px-4 py-2 rounded-xl text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer shadow-sm"
            >
              {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Create Team'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
