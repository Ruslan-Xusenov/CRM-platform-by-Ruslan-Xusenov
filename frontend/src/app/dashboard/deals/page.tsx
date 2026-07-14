'use client';

import React, { useEffect, useState } from 'react';
import { dealsAPI } from '@/lib/api';

export default function DealsPage() {
  const [deals, setDeals] = useState<any[]>([]);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    dealsAPI.list({ page: '1', page_size: '20' }).then((r) => { setDeals(r.data?.data || []); setTotal(r.data?.total || 0); }).catch(() => {});
  }, []);

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Bitimlar</h1>
          <p className="text-gray-400 text-sm mt-0.5">Jami: {total}</p>
        </div>
        <button className="px-4 py-2.5 rounded-xl bg-gradient-to-r from-purple-500 to-pink-600 text-white text-sm font-medium hover:shadow-lg transition-all">
          + Yangi bitim
        </button>
      </div>

      <div className="glass-card rounded-2xl overflow-hidden">
        <table className="w-full crm-table">
          <thead><tr className="border-b border-gray-800/40">
            <th className="text-left px-6 py-3">Nomi</th><th className="text-left px-4 py-3">Summa</th>
            <th className="text-left px-4 py-3">Ehtimollik</th><th className="text-left px-4 py-3">Yaratilgan</th>
          </tr></thead>
          <tbody>
            {deals.length === 0 ? (
              <tr><td colSpan={4} className="text-center py-12 text-gray-500">Bitimlar topilmadi</td></tr>
            ) : deals.map((d: any) => (
              <tr key={d.id}>
                <td className="px-6 font-medium text-white">{d.title}</td>
                <td className="px-4 text-gray-300">{d.amount ? `${Number(d.amount).toLocaleString()} ${d.currency}` : '—'}</td>
                <td className="px-4"><div className="flex items-center gap-2"><div className="w-16 h-1.5 bg-gray-700 rounded-full overflow-hidden"><div className="h-full bg-purple-500 rounded-full" style={{width: `${d.probability}%`}}/></div><span className="text-gray-400 text-xs">{d.probability}%</span></div></td>
                <td className="px-4 text-gray-500 text-xs">{new Date(d.created_at).toLocaleDateString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
