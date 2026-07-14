'use client';

import React, { useEffect, useState } from 'react';
import { callsAPI } from '@/lib/api';

export default function CallsPage() {
  const [calls, setCalls] = useState<any[]>([]);

  useEffect(() => {
    callsAPI.history({ page: '1', page_size: '50' }).then((r) => { setCalls(Array.isArray(r.data) ? r.data : r.data?.data || []); }).catch(() => {});
  }, []);

  const formatDuration = (s: number) => { const m = Math.floor(s / 60); return `${m}:${(s % 60).toString().padStart(2, '0')}`; };

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-2xl font-bold text-white">Qo'ng'iroqlar tarixi</h1>
        <p className="text-gray-400 text-sm mt-0.5">CDR — barcha qo'ng'iroqlar jurnali</p>
      </div>

      <div className="glass-card rounded-2xl overflow-hidden">
        <table className="w-full crm-table">
          <thead><tr className="border-b border-gray-800/40">
            <th className="text-left px-6 py-3">Yo'nalish</th>
            <th className="text-left px-4 py-3">Qo'ng'iroqchi</th>
            <th className="text-left px-4 py-3">Qabul qiluvchi</th>
            <th className="text-left px-4 py-3">Status</th>
            <th className="text-left px-4 py-3">Davomiylik</th>
            <th className="text-left px-4 py-3">Yozuv</th>
            <th className="text-left px-6 py-3">Sana</th>
          </tr></thead>
          <tbody>
            {calls.length === 0 ? (
              <tr><td colSpan={7} className="text-center py-12 text-gray-500">Qo'ng'iroqlar topilmadi</td></tr>
            ) : calls.map((c: any) => (
              <tr key={c.id}>
                <td className="px-6"><span className={`text-lg ${c.direction === 'inbound' ? '' : ''}`}>{c.direction === 'inbound' ? '📥' : '📤'}</span></td>
                <td className="px-4 text-white font-medium">{c.caller}</td>
                <td className="px-4 text-gray-400">{c.callee}</td>
                <td className="px-4"><span className={`badge ${c.status === 'completed' || c.status === 'answered' ? 'badge-active' : 'badge-lost'}`}>{c.status}</span></td>
                <td className="px-4 text-gray-300 font-mono text-sm">{formatDuration(c.duration_seconds || 0)}</td>
                <td className="px-4">{c.recording_url ? <a href="#" className="text-indigo-400 text-xs hover:underline">▶ Tinglash</a> : <span className="text-gray-600 text-xs">—</span>}</td>
                <td className="px-6 text-gray-500 text-xs">{new Date(c.started_at).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
