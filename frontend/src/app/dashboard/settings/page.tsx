'use client';

import React, { useEffect, useState } from 'react';
import { extensionsAPI } from '@/lib/api';
import { useCallStore } from '@/stores/callStore';

export default function SettingsPage() {
  const [extensions, setExtensions] = useState<any[]>([]);
  const { registered } = useCallStore();

  useEffect(() => {
    extensionsAPI.list().then((r) => { setExtensions(Array.isArray(r.data) ? r.data : []); }).catch(() => {});
  }, []);

  return (
    <div className="space-y-6">
      <div><h1 className="text-2xl font-bold text-white">Sozlamalar</h1><p className="text-gray-400 text-sm mt-0.5">PBX va tizim sozlamalari</p></div>

      {/* SoftPhone Status */}
      <div className="glass-card rounded-2xl p-6">
        <h2 className="text-lg font-semibold text-white mb-4">📞 SoftPhone holati</h2>
        <div className="flex items-center gap-3">
          <div className={`w-3 h-3 rounded-full ${registered ? 'bg-emerald-400 animate-pulse' : 'bg-red-400'}`}/>
          <span className="text-gray-300">{registered ? 'Ro\'yxatdan o\'tgan (Tayyor)' : 'Ulanmagan'}</span>
        </div>
        <p className="text-gray-500 text-sm mt-2">SIP.js orqali Asterisk WSS serveriga ulanish holati</p>
      </div>

      {/* Extensions */}
      <div className="glass-card rounded-2xl overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-800/40 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-white">Ichki raqamlar (Extensions)</h2>
          <button className="px-3 py-1.5 rounded-lg bg-indigo-500/15 text-indigo-400 text-sm border border-indigo-500/20 hover:bg-indigo-500/25 transition-all">+ Qo'shish</button>
        </div>
        <table className="w-full crm-table">
          <thead><tr className="border-b border-gray-800/40">
            <th className="text-left px-6 py-3">Raqam</th><th className="text-left px-4 py-3">Nomi</th>
            <th className="text-left px-4 py-3">Holat</th>
          </tr></thead>
          <tbody>
            {extensions.length === 0 ? (
              <tr><td colSpan={3} className="text-center py-8 text-gray-500">Extensions topilmadi</td></tr>
            ) : extensions.map((ext: any) => (
              <tr key={ext.id}>
                <td className="px-6 text-white font-mono font-medium">{ext.extension_number}</td>
                <td className="px-4 text-gray-400">{ext.display_name || '—'}</td>
                <td className="px-4"><span className={`badge ${ext.enabled ? 'badge-active' : 'badge-lost'}`}>{ext.enabled ? 'Faol' : 'O\'chiq'}</span></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* System Info */}
      <div className="glass-card rounded-2xl p-6">
        <h2 className="text-lg font-semibold text-white mb-4">ℹ️ Tizim haqida</h2>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div><span className="text-gray-500">Backend:</span> <span className="text-gray-300 ml-2">Go (Chi Router)</span></div>
          <div><span className="text-gray-500">Frontend:</span> <span className="text-gray-300 ml-2">Next.js 14</span></div>
          <div><span className="text-gray-500">PBX:</span> <span className="text-gray-300 ml-2">Asterisk 20 (ARI)</span></div>
          <div><span className="text-gray-500">Database:</span> <span className="text-gray-300 ml-2">PostgreSQL 16</span></div>
          <div><span className="text-gray-500">Cache:</span> <span className="text-gray-300 ml-2">Redis 7</span></div>
          <div><span className="text-gray-500">Storage:</span> <span className="text-gray-300 ml-2">MinIO (S3)</span></div>
        </div>
      </div>
    </div>
  );
}
