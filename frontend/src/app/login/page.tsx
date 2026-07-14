'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { authAPI } from '@/lib/api';
import { useAuthStore } from '@/stores/authStore';

export default function LoginPage() {
  const router = useRouter();
  const login = useAuthStore((s) => s.login);
  const [isRegister, setIsRegister] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [form, setForm] = useState({ email: '', password: '', first_name: '', last_name: '', tenant_name: '' });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const res = isRegister
        ? await authAPI.register(form)
        : await authAPI.login({ email: form.email, password: form.password });
      login(res.data.user, res.data.access_token, res.data.refresh_token);
      router.push('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Xatolik yuz berdi');
    } finally { setLoading(false); }
  };

  const set = (key: string, val: string) => setForm((f) => ({ ...f, [key]: val }));

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-950 relative overflow-hidden">
      {/* Background decoration */}
      <div className="absolute inset-0">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-indigo-600/10 rounded-full blur-3xl"/>
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-purple-600/10 rounded-full blur-3xl"/>
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-teal-600/5 rounded-full blur-3xl"/>
      </div>

      <div className="relative w-full max-w-md px-4">
        {/* Logo */}
        <div className="text-center mb-8 animate-fade-in">
          <div className="inline-flex w-16 h-16 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 items-center justify-center text-white text-xl font-bold shadow-2xl shadow-indigo-500/20 mb-4">
            CRM
          </div>
          <h1 className="text-3xl font-bold text-white mb-2">Omni Platform</h1>
          <p className="text-gray-400 text-sm">CRM & PBX tizimiga kirish</p>
        </div>

        {/* Form Card */}
        <div className="glass-card rounded-2xl p-8 animate-scale-in">
          {error && (
            <div className="mb-4 p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-sm">{error}</div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            {isRegister && (
              <>
                <div className="grid grid-cols-2 gap-3">
                  <input type="text" placeholder="Ism" value={form.first_name} onChange={(e) => set('first_name', e.target.value)} required
                    className="w-full px-4 py-3 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>
                  <input type="text" placeholder="Familiya" value={form.last_name} onChange={(e) => set('last_name', e.target.value)} required
                    className="w-full px-4 py-3 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>
                </div>
                <input type="text" placeholder="Kompaniya nomi" value={form.tenant_name} onChange={(e) => set('tenant_name', e.target.value)} required
                  className="w-full px-4 py-3 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>
              </>
            )}

            <input type="email" placeholder="Email" value={form.email} onChange={(e) => set('email', e.target.value)} required id="login-email"
              className="w-full px-4 py-3 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>

            <input type="password" placeholder="Parol" value={form.password} onChange={(e) => set('password', e.target.value)} required id="login-password"
              className="w-full px-4 py-3 rounded-xl bg-gray-800/60 border border-gray-700/50 text-white placeholder-gray-500 focus:border-indigo-500/50 focus:outline-none text-sm"/>

            <button type="submit" disabled={loading} id="login-submit"
              className="w-full py-3 rounded-xl bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold text-sm hover:shadow-lg hover:shadow-indigo-500/20 transition-all disabled:opacity-50 active:scale-[0.98]">
              {loading ? '...' : isRegister ? 'Ro\'yxatdan o\'tish' : 'Kirish'}
            </button>
          </form>

          <div className="mt-5 text-center">
            <button onClick={() => { setIsRegister(!isRegister); setError(''); }}
              className="text-sm text-gray-400 hover:text-indigo-400 transition-colors">
              {isRegister ? 'Akkauntingiz bormi? Kirish' : 'Yangi akkaunt yaratish'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
