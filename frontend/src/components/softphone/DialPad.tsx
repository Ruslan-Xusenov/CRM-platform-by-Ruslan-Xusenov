'use client';

import React, { useState } from 'react';
import { useSIP } from './SoftPhoneProvider';
import { useCallStore } from '@/stores/callStore';

const dialButtons = ['1','2','3','4','5','6','7','8','9','*','0','#'];

export default function DialPad() {
  const { makeCall, hangup, toggleMute, isMuted, callDuration, callStatus } = useSIP();
  const { softphoneOpen, toggleSoftphone, registered, currentCall } = useCallStore();
  const [number, setNumber] = useState('');

  const formatDuration = (s: number) => {
    const m = Math.floor(s / 60);
    const sec = s % 60;
    return `${m.toString().padStart(2, '0')}:${sec.toString().padStart(2, '0')}`;
  };

  if (!softphoneOpen) {
    return (
      <button onClick={toggleSoftphone} id="softphone-toggle"
        className="fixed bottom-6 right-6 z-50 w-14 h-14 rounded-full bg-gradient-to-br from-emerald-500 to-teal-600 text-white shadow-2xl flex items-center justify-center hover:scale-110 transition-all duration-200 hover:shadow-emerald-500/40">
        <svg width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07 19.5 19.5 0 01-6-6 19.79 19.79 0 01-3.07-8.67A2 2 0 014.11 2h3a2 2 0 012 1.72c.127.96.361 1.903.7 2.81a2 2 0 01-.45 2.11L8.09 9.91a16 16 0 006 6l1.27-1.27a2 2 0 012.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0122 16.92z"/>
        </svg>
        {!registered && <span className="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full animate-pulse"/>}
      </button>
    );
  }

  return (
    <div className="fixed bottom-6 right-6 z-50 w-80 bg-gray-900/95 backdrop-blur-xl rounded-2xl shadow-2xl border border-gray-700/50 overflow-hidden" id="softphone-dialpad">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 bg-gradient-to-r from-emerald-600/20 to-teal-600/20 border-b border-gray-700/50">
        <div className="flex items-center gap-2">
          <div className={`w-2 h-2 rounded-full ${registered ? 'bg-emerald-400 animate-pulse' : 'bg-red-400'}`}/>
          <span className="text-xs text-gray-300 font-medium">{registered ? 'Registered' : 'Offline'}</span>
        </div>
        <button onClick={toggleSoftphone} className="text-gray-400 hover:text-white transition-colors">
          <svg width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"/></svg>
        </button>
      </div>

      {/* Active Call Display */}
      {callStatus !== 'idle' && (
        <div className="px-4 py-4 text-center border-b border-gray-700/50">
          <p className="text-xs text-emerald-400 uppercase tracking-wider font-semibold mb-1">
            {callStatus === 'ringing' ? '📞 Incoming Call' : callStatus === 'calling' ? '📲 Calling...' : '🔊 In Call'}
          </p>
          <p className="text-white text-lg font-bold">{currentCall?.caller || currentCall?.callee || number}</p>
          {callStatus === 'active' && (
            <p className="text-emerald-400 text-2xl font-mono mt-1">{formatDuration(callDuration)}</p>
          )}
        </div>
      )}

      {/* Number Input */}
      {callStatus === 'idle' && (
        <div className="px-4 pt-4">
          <input type="text" value={number} onChange={(e) => setNumber(e.target.value)} placeholder="Enter number..."
            className="w-full text-center text-xl font-mono text-white bg-gray-800/60 rounded-xl px-4 py-3 border border-gray-700/50 focus:border-emerald-500/50 focus:outline-none placeholder-gray-500 tracking-widest"/>
        </div>
      )}

      {/* Dial Pad Grid */}
      {callStatus === 'idle' && (
        <div className="grid grid-cols-3 gap-2 px-4 py-3">
          {dialButtons.map((btn) => (
            <button key={btn} onClick={() => setNumber((n) => n + btn)}
              className="h-12 rounded-xl bg-gray-800/60 hover:bg-gray-700/80 text-white text-lg font-semibold transition-all duration-150 active:scale-95 border border-gray-700/30 hover:border-gray-600">
              {btn}
            </button>
          ))}
        </div>
      )}

      {/* Action Buttons */}
      <div className="flex items-center justify-center gap-3 px-4 pb-4">
        {callStatus === 'idle' ? (
          <>
            <button onClick={() => setNumber((n) => n.slice(0, -1))}
              className="w-12 h-12 rounded-full bg-gray-700/60 hover:bg-gray-600 text-white flex items-center justify-center transition-all">
              <svg width="18" height="18" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path d="M21 4H8l-7 8 7 8h13a2 2 0 002-2V6a2 2 0 00-2-2z"/><line x1="18" y1="9" x2="12" y2="15"/><line x1="12" y1="9" x2="18" y2="15"/></svg>
            </button>
            <button onClick={() => { if (number) makeCall(number); }} disabled={!number || !registered}
              className="w-14 h-14 rounded-full bg-gradient-to-br from-emerald-500 to-teal-600 text-white flex items-center justify-center shadow-lg hover:shadow-emerald-500/30 hover:scale-105 transition-all disabled:opacity-40 disabled:cursor-not-allowed">
              <svg width="22" height="22" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07 19.5 19.5 0 01-6-6 19.79 19.79 0 01-3.07-8.67A2 2 0 014.11 2h3a2 2 0 012 1.72c.127.96.361 1.903.7 2.81a2 2 0 01-.45 2.11L8.09 9.91a16 16 0 006 6l1.27-1.27a2 2 0 012.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0122 16.92z"/>
              </svg>
            </button>
          </>
        ) : (
          <>
            <button onClick={toggleMute}
              className={`w-12 h-12 rounded-full flex items-center justify-center transition-all ${isMuted ? 'bg-red-500/20 text-red-400 border border-red-500/40' : 'bg-gray-700/60 text-white hover:bg-gray-600'}`}>
              {isMuted ? '🔇' : '🎤'}
            </button>
            <button onClick={hangup}
              className="w-14 h-14 rounded-full bg-gradient-to-br from-red-500 to-rose-600 text-white flex items-center justify-center shadow-lg hover:shadow-red-500/30 hover:scale-105 transition-all">
              <svg width="22" height="22" fill="none" stroke="currentColor" strokeWidth="2.5" viewBox="0 0 24 24">
                <path d="M10.68 13.31a16 16 0 006 6l1.27-1.27" /><line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
            </button>
          </>
        )}
      </div>
    </div>
  );
}
