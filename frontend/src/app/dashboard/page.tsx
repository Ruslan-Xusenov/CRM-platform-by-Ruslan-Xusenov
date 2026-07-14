'use client';

import React, { useEffect, useState } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { leadsAPI, contactsAPI, dealsAPI, callsAPI } from '@/lib/api';

export default function DashboardPage() {
  const user = useAuthStore((s) => s.user);
  const [stats, setStats] = useState({ leads: 0, contacts: 0, deals: 0, calls: 0 });
  const [recentLeads, setRecentLeads] = useState<any[]>([]);

  useEffect(() => {
    const load = async () => {
      try {
        const [l, c, d, cl] = await Promise.all([
          leadsAPI.list({ page: '1', page_size: '1' }),
          contactsAPI.list({ page: '1', page_size: '1' }),
          dealsAPI.list({ page: '1', page_size: '1' }),
          callsAPI.history({ page: '1', page_size: '1' }),
        ]);
        setStats({
          leads: l.data?.total || 0, contacts: c.data?.total || 0,
          deals: d.data?.total || 0, calls: Array.isArray(cl.data) ? cl.data.length : 0,
        });
        const leadsRes = await leadsAPI.list({ page: '1', page_size: '5' });
        setRecentLeads(leadsRes.data?.data || []);
      } catch {}
    };
    load();
  }, []);

  const statCards = [
    { label: 'Jami Lidlar', value: stats.leads, icon: '🎯', color: 'from-blue-500/15 to-blue-600/5', border: 'border-blue-500/20', text: 'text-blue-400' },
    { label: 'Kontaktlar', value: stats.contacts, icon: '👥', color: 'from-emerald-500/15 to-emerald-600/5', border: 'border-emerald-500/20', text: 'text-emerald-400' },
    { label: 'Bitimlar', value: stats.deals, icon: '💰', color: 'from-purple-500/15 to-purple-600/5', border: 'border-purple-500/20', text: 'text-purple-400' },
    { label: 'Qo\'ng\'iroqlar', value: stats.calls, icon: '📞', color: 'from-amber-500/15 to-amber-600/5', border: 'border-amber-500/20', text: 'text-amber-400' },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-white">Salom, {user?.first_name} 👋</h1>
        <p className="text-gray-400 text-sm mt-1">Bugungi statistika va so'nggi faoliyat</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {statCards.map((card) => (
          <div key={card.label} className={`stat-card rounded-2xl p-5 bg-gradient-to-br ${card.color} border ${card.border}`}>
            <div className="flex items-center justify-between mb-3">
              <span className="text-2xl">{card.icon}</span>
              <span className={`text-xs font-medium ${card.text} bg-white/5 px-2 py-0.5 rounded-full`}>Jami</span>
            </div>
            <p className="text-3xl font-bold text-white">{card.value.toLocaleString()}</p>
            <p className="text-gray-400 text-sm mt-1">{card.label}</p>
          </div>
        ))}
      </div>

      {/* Recent Leads */}
      <div className="glass-card rounded-2xl overflow-hidden">
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800/60">
          <h2 className="text-lg font-semibold text-white">So'nggi Lidlar</h2>
          <a href="/dashboard/leads" className="text-sm text-indigo-400 hover:text-indigo-300 transition-colors">Barchasini ko'rish →</a>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full crm-table">
            <thead>
              <tr className="border-b border-gray-800/40">
                <th className="text-left px-6 py-3">Nomi</th>
                <th className="text-left px-4 py-3">Status</th>
                <th className="text-left px-4 py-3">Kontakt</th>
                <th className="text-left px-4 py-3">Manba</th>
                <th className="text-right px-6 py-3">Byudjet</th>
              </tr>
            </thead>
            <tbody>
              {recentLeads.length === 0 ? (
                <tr><td colSpan={5} className="text-center py-12 text-gray-500">Hozircha lidlar yo'q</td></tr>
              ) : recentLeads.map((lead: any) => (
                <tr key={lead.id} className="hover:bg-indigo-500/[0.03] cursor-pointer">
                  <td className="px-6 py-3">
                    <p className="font-medium text-white">{lead.title}</p>
                  </td>
                  <td className="px-4 py-3">
                    <span className={`badge ${lead.status === 'new' ? 'badge-new' : lead.status === 'active' ? 'badge-active' : 'badge-won'}`}>
                      {lead.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-gray-400">{lead.contact_name || '—'}</td>
                  <td className="px-4 py-3 text-gray-400">{lead.source || '—'}</td>
                  <td className="px-6 py-3 text-right text-gray-300 font-medium">
                    {lead.budget ? `${Number(lead.budget).toLocaleString()} ${lead.currency}` : '—'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
