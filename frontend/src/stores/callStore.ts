import { create } from 'zustand';

export interface ActiveCall {
  id: string;
  channel_id: string;
  caller: string;
  callee: string;
  direction: string;
  status: string;
  started_at: string;
}

interface CallState {
  activeCalls: ActiveCall[];
  currentCall: ActiveCall | null;
  incomingCall: ActiveCall | null;
  softphoneOpen: boolean;
  registered: boolean;
  setActiveCalls: (calls: ActiveCall[]) => void;
  setCurrentCall: (call: ActiveCall | null) => void;
  setIncomingCall: (call: ActiveCall | null) => void;
  toggleSoftphone: () => void;
  setSoftphoneOpen: (v: boolean) => void;
  setRegistered: (v: boolean) => void;
}

export const useCallStore = create<CallState>((set) => ({
  activeCalls: [],
  currentCall: null,
  incomingCall: null,
  softphoneOpen: false,
  registered: false,
  setActiveCalls: (activeCalls) => set({ activeCalls }),
  setCurrentCall: (currentCall) => set({ currentCall }),
  setIncomingCall: (incomingCall) => set({ incomingCall }),
  toggleSoftphone: () => set((s) => ({ softphoneOpen: !s.softphoneOpen })),
  setSoftphoneOpen: (softphoneOpen) => set({ softphoneOpen }),
  setRegistered: (registered) => set({ registered }),
}));
