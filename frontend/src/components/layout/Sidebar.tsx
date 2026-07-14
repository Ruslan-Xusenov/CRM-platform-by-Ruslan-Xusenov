'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuthStore } from '@/stores/authStore';

const navItems = [
  { href: '/dashboard', label: 'Dashboard', icon: '📊' },
  { href: '/dashboard/leads', label: 'Lidlar', icon: '🎯' },
  { href: '/dashboard/contacts', label: 'Kontaktlar', icon: '👥' },
  { href: '/dashboard/companies', label: 'Kompaniyalar', icon: '🏢' },
  { href: '/dashboard/deals', label: 'Bitimlar', icon: '💰' },
  { href: '/dashboard/calls', label: 'Qo\'ng\'iroqlar', icon: '📞' },
  { href: '/dashboard/settings', label: 'Sozlamalar', icon: '⚙️' },
];

export default function Sidebar() {
  const pathname = usePathname();
  const user = useAuthStore((s) => s.user);
  const logout = useAuthStore((s) => s.logout);
  const [collapsed, setCollapsed] = useState(false);

  return (
    <aside className={`fixed left-0 top-0 h-screen bg-gray-950 border-r border-gray-800/60 flex flex-col transition-all duration-300 z-40 ${collapsed ? 'w-[72px]' : 'w-64'}`}>
      {/* Logo */}
      <div className="flex items-center gap-3 px-4 h-16 border-b border-gray-800/60 shrink-0">
        <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-sm shadow-lg shadow-indigo-500/20">
          CRM
        </div>
        {!collapsed && <span className="text-white font-semibold text-sm tracking-tight">Omni Platform</span>}
      </div>

      {/* Navigation */}
      <nav className="flex-1 py-4 px-3 space-y-1 overflow-y-auto">
        {navItems.map((item) => {
          const isActive = pathname === item.href || (item.href !== '/dashboard' && pathname?.startsWith(item.href));
          return (
            <Link key={item.href} href={item.href}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group
                ${isActive
                  ? 'bg-gradient-to-r from-indigo-500/15 to-purple-500/10 text-indigo-400 border border-indigo-500/20 shadow-sm shadow-indigo-500/5'
                  : 'text-gray-400 hover:text-white hover:bg-gray-800/50'}`}>
              <span className="text-lg shrink-0">{item.icon}</span>
              {!collapsed && <span>{item.label}</span>}
              {isActive && !collapsed && <div className="ml-auto w-1.5 h-1.5 rounded-full bg-indigo-400"/>}
            </Link>
          );
        })}
      </nav>

      {/* Collapse Toggle */}
      <button onClick={() => setCollapsed(!collapsed)}
        className="mx-3 mb-2 p-2 rounded-lg text-gray-500 hover:text-gray-300 hover:bg-gray-800/50 transition-all text-center">
        {collapsed ? '→' : '←'}
      </button>

      {/* User Profile */}
      <div className="border-t border-gray-800/60 px-3 py-3 shrink-0">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-full bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center text-white text-xs font-bold shrink-0">
            {user?.first_name?.[0]}{user?.last_name?.[0]}
          </div>
          {!collapsed && (
            <div className="flex-1 min-w-0">
              <p className="text-white text-sm font-medium truncate">{user?.first_name} {user?.last_name}</p>
              <p className="text-gray-500 text-xs truncate">{user?.role}</p>
            </div>
          )}
          {!collapsed && (
            <button onClick={() => { logout(); window.location.href = '/login'; }}
              className="text-gray-500 hover:text-red-400 transition-colors p-1" title="Chiqish">
              <svg width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4M16 17l5-5-5-5M21 12H9"/>
              </svg>
            </button>
          )}
        </div>
      </div>
    </aside>
  );
}
