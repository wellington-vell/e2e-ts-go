import { format } from 'date-fns';
import { CalendarIcon, Loader2 } from 'lucide-react';
import React from 'react';
import type { DateRange } from 'react-day-picker';

import { useFieldContext, useFormContext } from '@/components/form/context';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import { Checkbox } from '@/components/ui/checkbox';
import { Field, FieldError, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { InputGroupInput } from '@/components/ui/input-group';
import { InputOTP } from '@/components/ui/input-otp';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';

export function FormInput({
  type,
  onKeyDown,
  ...props
}: React.ComponentProps<typeof Input>) {
  const field = useFieldContext<string | number>();
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
  const isNumber = type === 'number';
  const value = field.state.value;

  return (
    <Input
      aria-invalid={isInvalid}
      id={field.name}
      name={field.name}
      onBlur={field.handleBlur}
      onChange={(e) =>
        field.handleChange(isNumber ? e.target.valueAsNumber : e.target.value)
      }
      onKeyDown={(e) => {
        if (isNumber && e.key.length === 1 && /[^\d.\-+]/.test(e.key)) {
          e.preventDefault();
        }
        onKeyDown?.(e);
      }}
      type={type}
      value={typeof value === 'number' && Number.isNaN(value) ? '' : value}
      {...props}
    />
  );
}

export function FormLabel({
  ...props
}: React.ComponentProps<typeof FieldLabel>) {
  const field = useFieldContext();
  return <FieldLabel htmlFor={field.name} {...props} />;
}

export function FormField({ ...props }: React.ComponentProps<typeof Field>) {
  const field = useFieldContext();
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;

  return (
    <Field
      aria-invalid={isInvalid}
      data-invalid={isInvalid ? 'true' : undefined}
      {...props}
    />
  );
}

export function FormError({
  ...props
}: React.ComponentProps<typeof FieldError>) {
  const field = useFieldContext();
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;

  return isInvalid ? (
    <FieldError errors={field.state.meta.errors} {...props} />
  ) : null;
}

export function FormButton({
  type = 'submit',
  spinner,
  spinnerPosition = 'left',
  loading = false,
  loadingText = 'Loading...',
  ...props
}: React.ComponentProps<typeof Button> & {
  spinner?: React.ReactNode;
  loadingText?: React.ReactNode;
  spinnerPosition?: 'left' | 'right';
  loading?: boolean;
}) {
  const form = useFormContext();

  const spinnerprop = spinner ?? (
    <Loader2 className="size-4 animate-spin" data-slot="spinner" />
  );

  const loadingTextElement =
    typeof loadingText === 'string' ? (
      <span className="hidden md:block" data-slot="loading-text">
        {loadingText}
      </span>
    ) : (
      loadingText
    );

  return (
    <form.Subscribe>
      {({ isSubmitting, canSubmit }) => (
        <Button
          aria-disabled={!canSubmit}
          disabled={isSubmitting || props.disabled}
          type={type}
          {...props}
        >
          {isSubmitting || loading ? (
            <>
              {spinnerPosition === 'left' && spinnerprop}
              {loadingTextElement}
              {spinnerPosition === 'right' && spinnerprop}
            </>
          ) : (
            props.children
          )}
        </Button>
      )}
    </form.Subscribe>
  );
}

export function FormTextarea({
  ...props
}: React.ComponentProps<typeof Textarea>) {
  const field = useFieldContext<string>();
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
  return (
    <Textarea
      aria-invalid={isInvalid}
      id={field.name}
      name={field.name}
      onBlur={field.handleBlur}
      onChange={(e) => field.handleChange(e.target.value)}
      value={field.state.value}
      {...props}
    />
  );
}

type FormCalendarValue = Date | Date[] | DateRange | undefined;
function formatCalendarValue(
  mode: React.ComponentProps<typeof Calendar>['mode'],
  value: FormCalendarValue,
): string {
  if (!value) return '';
  if (mode === 'single' && value instanceof Date) {
    return format(value, 'PPP');
  }
  if (mode === 'multiple' && Array.isArray(value) && value.length > 0) {
    return `${value.length} date${value.length === 1 ? '' : 's'} selected`;
  }
  if (mode === 'range' && !Array.isArray(value) && !(value instanceof Date)) {
    if (value.from && value.to) {
      return `${format(value.from, 'PPP')} – ${format(value.to, 'PPP')}`;
    }
    if (value.from) {
      return format(value.from, 'PPP');
    }
  }
  return '';
}

export function FormCalendar({
  mode = 'single',
  placeholder = 'Pick a date',
  ...props
}: {
  placeholder?: string;
} & Omit<React.ComponentProps<typeof Calendar>, 'selected' | 'onSelect'>) {
  const field = useFieldContext<FormCalendarValue>();
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
  const value = field.state.value;
  const [open, setOpen] = React.useState(false);

  const handleOpenChange = (next: boolean) => {
    setOpen(next);
    if (!next) {
      field.handleBlur();
    }
  };

  const displayValue = formatCalendarValue(mode, value);

  let calendar: React.ReactNode;
  switch (mode) {
    case 'multiple':
      calendar = (
        <Calendar
          {...props}
          mode="multiple"
          selected={Array.isArray(value) ? value : undefined}
          onSelect={(next) => field.handleChange(next)}
        />
      );
      break;
    case 'range':
      calendar = (
        <Calendar
          {...props}
          mode="range"
          selected={
            value && !Array.isArray(value) && !(value instanceof Date)
              ? value
              : undefined
          }
          onSelect={(next) => field.handleChange(next)}
        />
      );
      break;
    default:
      calendar = (
        <Calendar
          {...props}
          mode="single"
          selected={value instanceof Date ? value : undefined}
          onSelect={(next) => {
            field.handleChange(next);
            setOpen(false);
          }}
        />
      );
  }

  return (
    <Popover open={open} onOpenChange={handleOpenChange}>
      <PopoverTrigger
        render={
          <Button
            aria-invalid={isInvalid}
            id={field.name}
            variant="outline"
            className={cn(
              'w-full justify-start text-left font-normal',
              !displayValue && 'text-muted-foreground',
            )}
          >
            <CalendarIcon className="size-4" />
            {displayValue || placeholder}
          </Button>
        }
      />
      <PopoverContent align="start" className="w-auto p-0">
        {calendar}
      </PopoverContent>
    </Popover>
  );
}

export function FormCheckbox({
  ...props
}: React.ComponentProps<typeof Checkbox>) {
  const field = useFieldContext<boolean>();
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;

  return (
    <Checkbox
      aria-invalid={isInvalid}
      checked={field.state.value}
      id={field.name}
      name={field.name}
      onBlur={field.handleBlur}
      onCheckedChange={(e) => field.handleChange(e)}
      {...props}
    />
  );
}

export function FormInputOTP({
  ...props
}: React.ComponentProps<typeof InputOTP>) {
  const field = useFieldContext<string>();
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;

  return (
    <InputOTP
      aria-invalid={isInvalid}
      id={field.name}
      name={field.name}
      onBlur={field.handleBlur}
      onChange={field.handleChange}
      value={field.state.value}
      {...props}
    />
  );
}

export function FormInputGroupInput({
  ...props
}: React.ComponentProps<typeof InputGroupInput>) {
  const field = useFieldContext<string>();
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;

  return (
    <InputGroupInput
      aria-invalid={isInvalid}
      id={field.name}
      name={field.name}
      onBlur={field.handleBlur}
      onChange={(e) => field.handleChange(e.target.value)}
      value={field.state.value}
      {...props}
    />
  );
}

export function FormSelect({
  children,
  ...props
}: React.ComponentProps<typeof Select>) {
  const field = useFieldContext();

  return (
    <Select
      value={field.state.value}
      onValueChange={(value) => {
        field.handleChange(value);
        field.handleBlur();
      }}
      {...props}
    >
      {children}
    </Select>
  );
}

export function FormSelectTrigger({
  ...props
}: React.ComponentProps<typeof SelectTrigger>) {
  const field = useFieldContext<string>();
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;

  return (
    <SelectTrigger aria-invalid={isInvalid} id={field.name} {...props}>
      <SelectValue />
    </SelectTrigger>
  );
}

export function FormSelectContent({
  ...props
}: React.ComponentProps<typeof SelectContent>) {
  return <SelectContent {...props} />;
}

export function FormSelectItem({
  ...props
}: React.ComponentProps<typeof SelectItem>) {
  return <SelectItem {...props} />;
}
