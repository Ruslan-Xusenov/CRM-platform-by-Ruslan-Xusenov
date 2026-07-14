import { create } from 'zustand';

export interface Lead {
  id: string;
  title: string;
  status: string;
  source?: string;
  budget?: number;
  currency: string;
  contact_name?: string;
  contact_phone?: string;
  contact_email?: string;
  company_name?: string;
  description?: string;
  custom_fields: Record<string, unknown>;
  assigned_to?: string;
  created_at: string;
}

export interface Contact {
  id: string;
  first_name: string;
  last_name: string;
  email?: string;
  phone?: string;
  position?: string;
  company_id?: string;
  created_at: string;
}

export interface Deal {
  id: string;
  title: string;
  pipeline_id?: string;
  stage_id?: string;
  contact_id?: string;
  amount?: number;
  currency: string;
  probability: number;
  created_at: string;
}

interface CRMState {
  leads: Lead[];
  contacts: Contact[];
  deals: Deal[];
  totalLeads: number;
  totalContacts: number;
  totalDeals: number;
  loading: boolean;
  setLeads: (leads: Lead[], total: number) => void;
  setContacts: (contacts: Contact[], total: number) => void;
  setDeals: (deals: Deal[], total: number) => void;
  addLead: (lead: Lead) => void;
  updateLead: (lead: Lead) => void;
  removeLead: (id: string) => void;
  setLoading: (v: boolean) => void;
}

export const useCRMStore = create<CRMState>((set) => ({
  leads: [],
  contacts: [],
  deals: [],
  totalLeads: 0,
  totalContacts: 0,
  totalDeals: 0,
  loading: false,
  setLeads: (leads, total) => set({ leads, totalLeads: total }),
  setContacts: (contacts, total) => set({ contacts, totalContacts: total }),
  setDeals: (deals, total) => set({ deals, totalDeals: total }),
  addLead: (lead) => set((s) => ({ leads: [lead, ...s.leads] })),
  updateLead: (lead) => set((s) => ({ leads: s.leads.map((l) => (l.id === lead.id ? lead : l)) })),
  removeLead: (id) => set((s) => ({ leads: s.leads.filter((l) => l.id !== id) })),
  setLoading: (loading) => set({ loading }),
}));
