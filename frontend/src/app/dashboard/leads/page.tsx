'use client';

import React, { useEffect, useState } from 'react';
import { leadsAPI } from '@/lib/api';
import { useCRMStore, Lead } from '@/stores/crmStore';

export default function LeadsPage() {
  const { leads, totalLeads, setLeads, loading, setLoading } = useCRMStore();
  const [page, setPage] = useState(1);
  const [showCreate, setShowCreate] = useState(false);
  const [form, setForm] = useState({ title: '', contact_name: '', contact_phone: '', source: '', budget: '' });

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      try {
        const res = await leadsAPI.list({ page: String(page), page_size: '20' });
        setLeads(res.data?.data || [], res.data?.total || 0);
      } catch {} finally { setLoading(false); }
    };
    load();
  }, [page, setLeads, setLoading]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await leadsAPI.create({ ...form, budget: form.budget ? Number(form.budget) : undefined });
      setShowCreate(false);
      setForm({ title: '', contact_name: '', contact_phone: '', source: '', budget: '' });
      const res = await leadsAPI.list({ page: '1', page_size: '20' });
      setLeads(res.data?.data || [], res.data?.total || 0);
    } catch {}
  };

  const handleDelete = async (id: string) => {
    if (!confirm('O\'chirmoqchimisiz?')) return;
    await leadsAPI.delete(id);
    const res = await leadsAPI.list({ page: String(page), page_size: '20' });
    setLeads(res.data?.data || [], res.data?.total || 0);
  };

  const set = (k: string, v: string) => setForm((f) => ({ ...f, [k]: v }));

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Lidlar</h1>
          <p className="text-gray-400 text-sm mt-0.5">Jami: {totalLeads}</p>
        </div>
        <button onClick={() => setShowCreate(true)} id="create-lead-btn"
          className="px-4 py-2.5 rounded-xl bg-gradient-to-r from-indigo-500 to-purple-600 text-white text-sm font-medium hover:shadow-lg hover:shadow-indigo-500/20 transition-all active:scale-[0.97]">
          + Yangi lid
        </button>
      </div>

      {/* Create Modal */}
      {showCreate && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm" onClick={() => setShowCreate(false)}>
          <div className="bg-gray-900 rounded-2xl border border-gray-700/50 p-6 w-full max-w-md animate-scale-in" onClick={(e) => e.stopPropagation()}>
            <h2 className="text-lg font-semibold text-white mb-4">Yangi lid yaratish</h2>
            <form onSubmit={handleCreate} className="space-y-3">
              <input placeholder="Nomi *" value={form.title} onChange={(e) => set('title', e.target.value)} required
                className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>
              <input placeholder="Kontakt ismi" value={form.contact_name} onChange={(e) => set('contact_name', e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>
              <input placeholder="Telefon" value={form.contact_phone} onChange={(e) => set('contact_phone', e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>
              <div className="grid grid-cols-2 gap-3">
                <input placeholder="Manba" value={form.source} onChange={(e) => set('source', e.target.value)}
                  className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>
                <input placeholder="Byudjet" type="number" value={form.budget} onChange={(e) => set('budget', e.target.value)}
                  className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>
              </div>
              <div className="flex gap-3 pt-2">
                <button type="button" onClick={() => setShowCreate(false)} className="flex-1 py-2.5 rounded-xl border border-gray-700/50 text-gray-400 text-sm hover:bg-gray-800/50">Bekor</button>
                <button type="submit" className="flex-1 py-2.5 rounded-xl bg-gradient-to-r from-indigo-500 to-purple-600 text-white text-sm font-medium">Yaratish</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Table */}
      <div className="glass-card rounded-2xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full crm-table">
            <thead>
              <tr className="border-b border-gray-800/40">
                <th className="text-left px-6 py-3">Nomi</th>
                <th className="text-left px-4 py-3">Status</th>
                <th className="text-left px-4 py-3">Kontakt</th>
                <th className="text-left px-4 py-3">Telefon</th>
                <th className="text-left px-4 py-3">Manba</th>
                <th className="text-right px-4 py-3">Byudjet</th>
                <th className="text-center px-4 py-3">Amallar</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr><td colSpan={7} className="text-center py-12 text-gray-500">Yuklanmoqda...</td></tr>
              ) : leads.length === 0 ? (
                <tr><td colSpan={7} className="text-center py-12 text-gray-500">Lidlar topilmadi</td></tr>
              ) : leads.map((lead: Lead) => (
                <tr key={lead.id}>
                  <td className="px-6"><span className="font-medium text-white">{lead.title}</span></td>
                  <td className="px-4"><span className={`badge ${lead.status === 'new' ? 'badge-new' : lead.status === 'converted' ? 'badge-won' : 'badge-active'}`}>{lead.status}</span></td>
                  <td className="px-4 text-gray-400">{lead.contact_name || '—'}</td>
                  <td className="px-4 text-gray-400">{lead.contact_phone || '—'}</td>
                  <td className="px-4 text-gray-400">{lead.source || '—'}</td>
                  <td className="px-4 text-right text-gray-300">{lead.budget ? `${Number(lead.budget).toLocaleString()} ${lead.currency}` : '—'}</td>
                  <td className="px-4 text-center">
                    <button onClick={() => handleDelete(lead.id)} className="text-gray-500 hover:text-red-400 transition-colors text-xs">🗑</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {totalLeads > 20 && (
          <div className="flex items-center justify-center gap-2 px-6 py-4 border-t border-gray-800/40">
            <button disabled={page <= 1} onClick={() => setPage(page - 1)} className="px-3 py-1.5 rounded-lg bg-gray-800/60 text-sm text-gray-400 disabled:opacity-30">←</button>
            <span className="text-sm text-gray-400">Sahifa {page}</span>
            <button onClick={() => setPage(page + 1)} className="px-3 py-1.5 rounded-lg bg-gray-800/60 text-sm text-gray-400">→</button>
          </div>
        )}
      </div>
    </div>
  );
}
