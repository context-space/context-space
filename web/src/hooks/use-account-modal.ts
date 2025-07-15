import type { TabType } from "@/components/auth/account-tabs/types"
import { atom, useAtom } from "jotai"

// Global state for account modal
interface AccountModalState {
  isOpen: boolean
  activeTab: TabType
}

const accountModalAtom = atom<AccountModalState>({
  isOpen: false,
  activeTab: "profile",
})

export function useAccountModal() {
  const [state, setState] = useAtom(accountModalAtom)

  const openModal = (tab: TabType = "profile") => {
    setState({ isOpen: true, activeTab: tab })
  }

  const closeModal = () => {
    setState(prev => ({ ...prev, isOpen: false }))
  }

  const setActiveTab = (tab: TabType) => {
    setState(prev => ({ ...prev, activeTab: tab }))
  }

  return {
    isOpen: state.isOpen,
    activeTab: state.activeTab,
    openModal,
    closeModal,
    setActiveTab,
  }
}
