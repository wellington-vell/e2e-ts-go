import { createContext, useCallback, useMemo, useState } from 'react';
import React from 'react';

const CLIENT_TIMEZONE_COOKIE_NAME = 'client-timezone';

const COOKIE_MAX_AGE = 31_536_000; // 1 year

function getClientTimezone(): string {
  if (typeof window === 'undefined') {
    return Intl.DateTimeFormat().resolvedOptions().timeZone;
  }
  const match = document.cookie
    .split('; ')
    .find((row) => row.startsWith(`${CLIENT_TIMEZONE_COOKIE_NAME}=`));
  const raw = match?.split('=').slice(1).join('=');
  const value = raw ? decodeURIComponent(raw) : null;
  if (value === null || value === '' || value === 'system') {
    return Intl.DateTimeFormat().resolvedOptions().timeZone;
  }
  return value;
}

function setClientTimezone(ianaValue: string): void {
  const sameSite = 'Lax';
  const secure =
    typeof window !== 'undefined' && window.location.protocol === 'https:'
      ? '; Secure'
      : '';
  if (ianaValue === 'system') {
    document.cookie = `${CLIENT_TIMEZONE_COOKIE_NAME}=; path=/; max-age=0; SameSite=${sameSite}${secure}`;
    return;
  }
  document.cookie = `${CLIENT_TIMEZONE_COOKIE_NAME}=${encodeURIComponent(ianaValue)}; path=/; max-age=${COOKIE_MAX_AGE}; SameSite=${sameSite}${secure}`;
}

export const TIMEZONE_OPTIONS: { value: string; label: string }[] = [
  { value: 'system', label: 'System' },
  { value: 'UTC', label: 'UTC' },
  { value: 'America/Sao_Paulo', label: 'Brazil' },
  { value: 'America/New_York', label: 'New York' },
  { value: 'Europe/London', label: 'London' },
  { value: 'Europe/Paris', label: 'Paris' },
  { value: 'Asia/Tokyo', label: 'Tokyo' },
];

function getTimezoneLabel(iana: string): string {
  if (iana === 'system' || !iana) return 'System';
  const found = TIMEZONE_OPTIONS.find((o) => o.value === iana);
  return found?.label ?? iana;
}

type TimezoneContextValue = {
  timezone: string;
  timezoneLabel: string;
  selectedValue: string;
  setTimezone: (ianaValue: string) => void;
};

const TimezoneContext = createContext<TimezoneContextValue | null>(null);

function getDefaultTimezoneValue(): TimezoneContextValue {
  const timezone = getClientTimezone();
  const systemTz = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const selectedValue = timezone === systemTz ? 'system' : timezone;
  const timezoneLabel = getTimezoneLabel(selectedValue);
  return {
    timezone,
    timezoneLabel,
    selectedValue,
    setTimezone: () => {},
  };
}

const defaultTimezoneValue = getDefaultTimezoneValue();

export function useTimezone(): TimezoneContextValue {
  const ctx = React.use(TimezoneContext);
  if (!ctx) return defaultTimezoneValue;
  return ctx;
}

export function TimezoneProvider({ children }: { children: React.ReactNode }) {
  const systemTz = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const [storedValue, setStoredValue] = useState<string>(() => {
    if (typeof window === 'undefined') return 'system';
    const tz = getClientTimezone();
    return tz === systemTz ? 'system' : tz;
  });

  const timezone = storedValue === 'system' ? getClientTimezone() : storedValue;
  const timezoneLabel = getTimezoneLabel(storedValue);

  const setTimezone = useCallback((ianaValue: string) => {
    if (ianaValue === 'system') {
      setStoredValue('system');
    } else {
      setStoredValue(ianaValue);
    }
    setClientTimezone(ianaValue);
  }, []);

  const value = useMemo<TimezoneContextValue>(
    () => ({
      timezone,
      timezoneLabel,
      selectedValue: storedValue,
      setTimezone,
    }),
    [timezone, timezoneLabel, storedValue, setTimezone],
  );

  return (
    <TimezoneContext.Provider value={value}>
      {children}
    </TimezoneContext.Provider>
  );
}
