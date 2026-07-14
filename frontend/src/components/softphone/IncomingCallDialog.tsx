'use client';

import React from 'react';
import { useCallStore } from '@/stores/callStore';
import { useSIP } from './SoftPhoneProvider';

export default function IncomingCallDialog() {
  const { incomingCall } = useCallStore();
  const { answer, hangup } = useSIP();

  if (!incomingCall) return null;

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm animate-in fade-in">
      <div className="bg-gray-900 rounded-3xl shadow-2xl border border-gray-700/50 p-8 w-80 text-center animate-in zoom-in-95">
        {/* Pulsing ring animation */}
        <div className="relative mx-auto w-20 h-20 mb-6">
          <div className="absolute inset-0 rounded-full bg-emerald-500/20 animate-ping"/>
          <div className="absolute inset-2 rounded-full bg-emerald-500/30 animate-pulse"/>
          <div className="relative w-full h-full rounded-full bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center">
            <svg width="32" height="32" fill="none" stroke="white" strokeWidth="2" viewBox="0 0 24 24">
              <path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07 19.5 19.5 0 01-6-6 19.79 19.79 0 01-3.07-8.67A2 2 0 014.11 2h3a2 2 0 012 1.72c.127.96.361 1.903.7 2.81a2 2 0 01-.45 2.11L8.09 9.91a16 16 0 006 6l1.27-1.27a2 2 0 012.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0122 16.92z"/>
            </svg>
          </div>
        </div>

        <p className="text-emerald-400 text-xs uppercase tracking-wider font-semibold mb-2">Incoming Call</p>
        <p className="text-white text-2xl font-bold mb-1">{incomingCall.caller}</p>
        <p className="text-gray-400 text-sm mb-8">{incomingCall.direction === 'inbound' ? 'External Call' : 'Internal Call'}</p>

        <div className="flex items-center justify-center gap-6">
          <button onClick={hangup} id="reject-call"
            className="w-16 h-16 rounded-full bg-gradient-to-br from-red-500 to-rose-600 text-white flex items-center justify-center shadow-lg hover:shadow-red-500/40 hover:scale-110 transition-all">
            <svg width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24">
              <path d="M23 1L1 23M16.92 22A19.79 19.79 0 018.29 18.93"/>
            </svg>
          </button>
          <button onClick={answer} id="accept-call"
            className="w-16 h-16 rounded-full bg-gradient-to-br from-emerald-500 to-teal-600 text-white flex items-center justify-center shadow-lg hover:shadow-emerald-500/40 hover:scale-110 transition-all animate-bounce">
            <svg width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07 19.5 19.5 0 01-6-6 19.79 19.79 0 01-3.07-8.67A2 2 0 014.11 2h3a2 2 0 012 1.72c.127.96.361 1.903.7 2.81a2 2 0 01-.45 2.11L8.09 9.91a16 16 0 006 6l1.27-1.27a2 2 0 012.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0122 16.92z"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
