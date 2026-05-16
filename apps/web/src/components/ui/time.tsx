import { formatInTimeZone } from 'date-fns-tz';

import { useTimezone } from '@/context/timezone';
import { cn } from '@/lib/utils';

/** Allowed date-fns format strings for Time. See https://date-fns.org/v4.1.0/docs/format */
const TIME_FORMATS = [
  'dd-MM-yyyy HH:mm:ss', // 14-02-2025 15:30:00
  'yyyy-MM-dd', // 2025-02-14
  "yyyy-MM-dd'T'HH:mm:ss", // 2025-02-14T15:30:00
  'MMM d, yyyy', // Feb 14, 2025
  'MMM d', // Feb 14
  'EEE', // Fri
  'EEEE', // Friday
  'HH:mm', // 15:30
  'P', // 02/14/2025
  'PP', // Feb 14, 2025
  'PPP', // February 14th, 2025
  'PPPP', // Friday, February 14th, 2025
  'p', // 3:30 PM
  'pp', // 3:30:00 PM
  'ppp', // 3:30:00 PM GMT+1
  'pppp', // 3:30:00 PM GMT+01:00
  'Pp', // 02/14/2025, 3:30 PM
  'PPpp', // Feb 14, 2025, 3:30:00 PM
  'PPPppp', // February 14th, 2025 at 3:30:00 PM GMT+1
  'PPPPpppp', // Friday, February 14th, 2025 at 3:30:00 PM GMT+01:00
] as const;

type TimeFormatStr = (typeof TIME_FORMATS)[number];

interface TimeProps {
  date: Date | null | undefined;
  formatStr?: TimeFormatStr;
}

export function Time({
  className,
  date,
  formatStr = 'dd-MM-yyyy HH:mm:ss',
  ...props
}: React.ComponentProps<'time'> & TimeProps) {
  const { timezone } = useTimezone();
  if (date == null || Number.isNaN(date.getTime())) return null;
  const formatted = formatInTimeZone(date, timezone, formatStr);
  const utc = date.toISOString();
  return (
    <time
      dateTime={utc}
      title={`${utc} (${timezone})`}
      className={cn('text-muted-foreground', className)}
      {...props}
    >
      {formatted}
    </time>
  );
}
