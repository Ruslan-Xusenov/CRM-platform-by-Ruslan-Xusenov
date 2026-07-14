'use client';

import React, { useEffect, useState } from 'react';
import { contactsAPI } from '@/lib/api';

export default function ContactsPage() {
  const [contacts, setContacts] = useState<any[]>([]);
  const [total, setTotal] = useState(0);
  const [showCreate, setShowCreate] = useState(false);
  const [form, setForm] = useState({ first_name: '', last_name: '', email: '', phone: '', position: '' });

  useEffect(() => {
    contactsAPI.list({ page: '1', page_size: '20' }).then((r) => { setContacts(r.data?.data || []); setTotal(r.data?.total || 0); }).catch(() => {});
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    await contactsAPI.create(form);
    setShowCreate(false);
    setForm({ first_name: '', last_name: '', email: '', phone: '', position: '' });
    const r = await contactsAPI.list({ page: '1', page_size: '20' });
    setContacts(r.data?.data || []); setTotal(r.data?.total || 0);
  };

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Kontaktlar</h1>
          <p className="text-gray-400 text-sm mt-0.5">Jami: {total}</p>
        </div>
        <button onClick={() => setShowCreate(true)}
          className="px-4 py-2.5 rounded-xl bg-gradient-to-r from-emerald-500 to-teal-600 text-white text-sm font-medium hover:shadow-lg transition-all">
          + Yangi kontakt
        </button>
      </div>

      {showCreate && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm" onClick={() => setShowCreate(false)}>
          <div className="bg-gray-900 rounded-2xl border border-gray-700/50 p-6 w-full max-w-md animate-scale-in" onClick={(e) => e.stopPropagation()}>
            <h2 className="text-lg font-semibold text-white mb-4">Yangi kontakt</h2>
            <form onSubmit={handleCreate} className="space-y-3">
              <div className="grid grid-cols-2 gap-3">
                <input placeholder="Ism *" value={form.first_name} onChange={(e) => setForm({...form, first_name: e.target.value})} required
                  className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-emerald-500/50 focus:outline-none text-sm"/>
                <input placeholder="Familiya *" value={form.last_name} onChange={(e) => setForm({...form, last_name: e.target.value})} required
                  className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-emerald-500/50 focus:outline-none text-sm"/>
              </div>
              <input placeholder="Email" value={form.email} onChange={(e) => setForm({...form, email: e.target.value})}
                className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-emerald-500/50 focus:outline-none text-sm"/>
              <input placeholder="Telefon" value={form.phone} onChange={(e) => setForm({...form, phone: e.target.value})}
                className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-emerald-500/50 focus:outline-none text-sm"/>
              <input placeholder="Lavozim" value={form.position} onChange={(e) => setForm({...form, position: e.target.value})}
                className="w-full px-4 py-2.5 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-emerald-500/50 focus:outline-none text-sm"/>
              <div className="flex gap-3 pt-2">
                <button type="button" onClick={() => setShowCreate(false)} className="flex-1 py-2.5 rounded-xl border border-gray-700/50 text-gray-400 text-sm">Bekor</button>
                <button type="submit" className="flex-1 py-2.5 rounded-xl bg-gradient-to-r from-emerald-500 to-teal-600 text-white text-sm font-medium">Yaratish</button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="glass-card rounded-2xl overflow-hidden">
        <table className="w-full crm-table">
          <thead><tr className="border-b border-gray-800/40">
            <th className="text-left px-6 py-3">Ism</th><th className="text-left px-4 py-3">Email</th>
            <th className="text-left px-4 py-3">Telefon</th><th className="text-left px-4 py-3">Lavozim</th>
          </tr></thead>
          <tbody>
            {contacts.length === 0 ? (
              <tr><td colSpan={4} className="text-center py-12 text-gray-500">Kontaktlar topilmadi</td></tr>
            ) : contacts.map((c: any) => (
              <tr key={c.id}>
                <td className="px-6 font-medium text-white">{c.first_name} {c.last_name}</td>
                <td className="px-4 text-gray-400">{c.email || '—'}</td>
                <td className="px-4 text-gray-400">{c.phone || '—'}</td>
                <td className="px-4 text-gray-400">{c.position || '—'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
