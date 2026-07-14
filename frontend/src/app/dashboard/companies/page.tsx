'use client';

import React, { useEffect, useState } from 'react';
import api from '@/lib/api';

export default function CompaniesPage() {
  const [companies, setCompanies] = useState<any[]>([]);
  const [total, setTotal] = useState(0);
  const [showCreate, setShowCreate] = useState(false);
  const [form, setForm] = useState({ name: '', industry: '', phone: '', email: '', website: '' });

  useEffect(() => {
    api.get('/companies', { params: { page: '1', page_size: '20' } }).then((r) => { setCompanies(r.data?.data || []); setTotal(r.data?.total || 0); }).catch(() => {});
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    await api.post('/companies', form);
    setShowCreate(false);
    setForm({ name: '', industry: '', phone: '', email: '', website: '' });
    const r = await api.get('/companies', { params: { page: '1', page_size: '20' } });
    setCompanies(r.data?.data || []); setTotal(r.data?.total || 0);
  };

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div><h1 className="text-2xl font-bold text-white">Kompaniyalar</h1><p className="text-gray-400 text-sm mt-0.5">Jami: {total}</p></div>
        <button onClick={() => setShowCreate(true)} className="px-4 py-2.5 rounded-xl bg-gradient-to-r from-amber-500 to-orange-600 text-white text-sm font-medium hover:shadow-lg transition-all">+ Yangi kompaniya</button>
      </div>

      {showCreate && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm" onClick={() => setShowCreate(false)}>
          <div className="bg-gray-900 rounded-2xl border border-gray-700/50 p-6 w-full max-w-md animate-scale-in" onClick={(e) => e.stopPropagation()}>
            <h2 className="text-lg font-semibold text-white mb-4">Yangi kompaniya</h2>
            <form onSubmit={handleCreate} className="space-y-3">
              <input placeholder="Nomi *" value={form.name} onChange={(e) => setForm({...form, name: e.target.value})} required className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-amber-500/50 focus:outline-none text-sm"/>
              <input placeholder="Soha" value={form.industry} onChange={(e) => setForm({...form, industry: e.target.value})} className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-amber-500/50 focus:outline-none text-sm"/>
              <input placeholder="Telefon" value={form.phone} onChange={(e) => setForm({...form, phone: e.target.value})} className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-amber-500/50 focus:outline-none text-sm"/>
              <input placeholder="Email" value={form.email} onChange={(e) => setForm({...form, email: e.target.value})} className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-amber-500/50 focus:outline-none text-sm"/>
              <input placeholder="Veb-sayt" value={form.website} onChange={(e) => setForm({...form, website: e.target.value})} className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-amber-500/50 focus:outline-none text-sm"/>
              <div className="flex gap-3 pt-2">
                <button type="button" onClick={() => setShowCreate(false)} className="flex-1 py-2.5 rounded-xl border border-gray-700/50 text-gray-400 text-sm">Bekor</button>
                <button type="submit" className="flex-1 py-2.5 rounded-xl bg-gradient-to-r from-amber-500 to-orange-600 text-white text-sm font-medium">Yaratish</button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="glass-card rounded-2xl overflow-hidden">
        <table className="w-full crm-table">
          <thead><tr className="border-b border-gray-800/40">
            <th className="text-left px-6 py-3">Nomi</th><th className="text-left px-4 py-3">Soha</th>
            <th className="text-left px-4 py-3">Telefon</th><th className="text-left px-4 py-3">Email</th>
          </tr></thead>
          <tbody>
            {companies.length === 0 ? (
              <tr><td colSpan={4} className="text-center py-12 text-gray-500">Kompaniyalar topilmadi</td></tr>
            ) : companies.map((c: any) => (
              <tr key={c.id}><td className="px-6 font-medium text-white">{c.name}</td><td className="px-4 text-gray-400">{c.industry || '—'}</td><td className="px-4 text-gray-400">{c.phone || '—'}</td><td className="px-4 text-gray-400">{c.email || '—'}</td></tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
