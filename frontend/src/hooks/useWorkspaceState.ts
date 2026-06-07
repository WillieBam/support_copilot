import { useState } from 'react'

export function useWorkspaceState() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true)
  const toggleSidebar = () => setIsSidebarOpen((prev) => !prev)

  return { isSidebarOpen, toggleSidebar }
}
