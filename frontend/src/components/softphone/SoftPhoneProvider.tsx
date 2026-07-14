'use client';

import React, { createContext, useContext, useEffect, useRef, useState, ReactNode } from 'react';
import { UserAgent, Registerer, Inviter, SessionState, Web } from 'sip.js';
import { useCallStore } from '@/stores/callStore';

interface SIPContextType {
  makeCall: (target: string) => void;
  hangup: () => void;
  answer: () => void;
  toggleMute: () => void;
  isMuted: boolean;
  callDuration: number;
  callStatus: string;
}

const SIPContext = createContext<SIPContextType | null>(null);

export const useSIP = () => {
  const ctx = useContext(SIPContext);
  if (!ctx) throw new Error('useSIP must be used within SoftPhoneProvider');
  return ctx;
};

interface Props { children: ReactNode; sipServer?: string; extension?: string; password?: string; }

export function SoftPhoneProvider({ children, sipServer, extension, password }: Props) {
  const uaRef = useRef<UserAgent | null>(null);
  const sessionRef = useRef<any>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const [isMuted, setIsMuted] = useState(false);
  const [callDuration, setCallDuration] = useState(0);
  const [callStatus, setCallStatus] = useState('idle');
  const setRegistered = useCallStore((s) => s.setRegistered);
  const setIncomingCall = useCallStore((s) => s.setIncomingCall);
  const setCurrentCall = useCallStore((s) => s.setCurrentCall);

  useEffect(() => {
    if (!sipServer || !extension || !password) return;

    const server = sipServer.startsWith('wss://') ? sipServer : `wss://${sipServer}`;
    const uri = UserAgent.makeURI(`sip:${extension}@${sipServer.replace('wss://', '').replace(':8089', '')}`);
    if (!uri) return;

    const ua = new UserAgent({
      uri,
      transportOptions: { server },
      authorizationUsername: extension,
      authorizationPassword: password,
      delegate: {
        onInvite: (invitation) => {
          setCallStatus('ringing');
          setIncomingCall({ id: invitation.id, channel_id: '', caller: invitation.remoteIdentity.uri.user || 'Unknown', callee: extension, direction: 'inbound', status: 'ringing', started_at: new Date().toISOString() });
          sessionRef.current = invitation;
          invitation.stateChange.addListener((state: SessionState) => {
            if (state === SessionState.Terminated) { cleanup(); }
          });
        },
      },
    });

    const registerer = new Registerer(ua);
    ua.start().then(() => {
      registerer.register().then(() => setRegistered(true)).catch(() => setRegistered(false));
    }).catch(console.error);

    uaRef.current = ua;
    return () => { ua.stop(); };
  }, [sipServer, extension, password, setRegistered, setIncomingCall]);

  const startTimer = () => {
    setCallDuration(0);
    timerRef.current = setInterval(() => setCallDuration((d) => d + 1), 1000);
  };

  const cleanup = () => {
    if (timerRef.current) clearInterval(timerRef.current);
    setCallDuration(0);
    setCallStatus('idle');
    setIsMuted(false);
    setCurrentCall(null);
    setIncomingCall(null);
    sessionRef.current = null;
  };

  const setupAudio = (session: any) => {
    const pc = session.sessionDescriptionHandler?.peerConnection;
    if (!pc) return;
    pc.ontrack = (e: RTCTrackEvent) => {
      if (!audioRef.current) { audioRef.current = new Audio(); }
      audioRef.current.srcObject = e.streams[0];
      audioRef.current.play().catch(console.error);
    };
  };

  const makeCall = (target: string) => {
    if (!uaRef.current) return;
    const uri = UserAgent.makeURI(`sip:${target}@${sipServer?.replace('wss://', '').replace(':8089', '')}`);
    if (!uri) return;

    const inviter = new Inviter(uaRef.current, uri);
    setCallStatus('calling');
    setCurrentCall({ id: inviter.id, channel_id: '', caller: extension || '', callee: target, direction: 'outbound', status: 'calling', started_at: new Date().toISOString() });

    inviter.invite().then(() => {
      sessionRef.current = inviter;
      setupAudio(inviter);
      inviter.stateChange.addListener((state: SessionState) => {
        if (state === SessionState.Established) { setCallStatus('active'); startTimer(); }
        if (state === SessionState.Terminated) { cleanup(); }
      });
    }).catch(console.error);
  };

  const hangup = () => {
    if (!sessionRef.current) return;
    const state = sessionRef.current.state;
    if (state === SessionState.Established) sessionRef.current.bye();
    else if (state === SessionState.Establishing) sessionRef.current.cancel();
    cleanup();
  };

  const answer = () => {
    if (!sessionRef.current?.accept) return;
    sessionRef.current.accept().then(() => {
      setCallStatus('active');
      startTimer();
      setupAudio(sessionRef.current);
    }).catch(console.error);
  };

  const toggleMute = () => {
    const pc = sessionRef.current?.sessionDescriptionHandler?.peerConnection;
    if (!pc) return;
    pc.getSenders().forEach((s: RTCRtpSender) => { if (s.track?.kind === 'audio') s.track.enabled = isMuted; });
    setIsMuted(!isMuted);
  };

  return (
    <SIPContext.Provider value={{ makeCall, hangup, answer, toggleMute, isMuted, callDuration, callStatus }}>
      {children}
    </SIPContext.Provider>
  );
}
