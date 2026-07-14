'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/stores/authStore';
import { useWebSocket } from '@/hooks/useWebSocket';
import Sidebar from '@/components/layout/Sidebar';
import DialPad from '@/components/softphone/DialPad';
import IncomingCallDialog from '@/components/softphone/IncomingCallDialog';
import { SoftPhoneProvider } from '@/components/softphone/SoftPhoneProvider';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const { isAuthenticated, loadFromStorage } = useAuthStore();

  useEffect(() => { loadFromStorage(); }, [loadFromStorage]);
  useEffect(() => {
    if (!isAuthenticated && typeof window !== 'undefined' && !localStorage.getItem('access_token')) {
      router.push('/login');
    }
  }, [isAuthenticated, router]);

  useWebSocket();

  return (
    <SoftPhoneProvider sipServer={process.env.NEXT_PUBLIC_SIP_SERVER} extension={process.env.NEXT_PUBLIC_SIP_EXTENSION} password={process.env.NEXT_PUBLIC_SIP_PASSWORD}>
      <div className="flex min-h-screen bg-gray-950">
        <Sidebar />
        <main className="flex-1 ml-64 p-6 overflow-auto">
          <div className="max-w-[1400px] mx-auto animate-fade-in">
            {children}
          </div>
        </main>
        <DialPad />
        <IncomingCallDialog />
      </div>
    </SoftPhoneProvider>
  );
}
