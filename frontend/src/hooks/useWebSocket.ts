'use client';

import { useEffect, useRef, useCallback } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { useCRMStore } from '@/stores/crmStore';
import { useCallStore } from '@/stores/callStore';

export function useWebSocket() {
  const wsRef = useRef<WebSocket | null>(null);
  const user = useAuthStore((s) => s.user);
  const addLead = useCRMStore((s) => s.addLead);
  const updateLead = useCRMStore((s) => s.updateLead);
  const removeLead = useCRMStore((s) => s.removeLead);
  const setIncomingCall = useCallStore((s) => s.setIncomingCall);

  const connect = useCallback(() => {
    if (!user) return;
    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';
    const ws = new WebSocket(`${wsUrl}?user_id=${user.id}&tenant_id=${user.tenant_id}`);

    ws.onopen = () => console.log('[WS] Connected');
    ws.onclose = () => {
      console.log('[WS] Disconnected, reconnecting in 3s...');
      setTimeout(connect, 3000);
    };
    ws.onerror = (e) => console.error('[WS] Error', e);

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        switch (msg.type) {
          case 'lead.created': addLead(msg.payload); break;
          case 'lead.updated': updateLead(msg.payload); break;
          case 'lead.deleted': removeLead(msg.payload.id); break;
          case 'call.incoming': setIncomingCall(msg.payload); break;
          case 'call.ended': setIncomingCall(null); break;
        }
      } catch (e) { console.error('[WS] Parse error', e); }
    };

    wsRef.current = ws;
  }, [user, addLead, updateLead, removeLead, setIncomingCall]);

  useEffect(() => {
    connect();
    return () => { wsRef.current?.close(); };
  }, [connect]);

  return wsRef;
}
