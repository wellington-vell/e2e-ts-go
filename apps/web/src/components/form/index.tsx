import { Loader2 } from 'lucide-react';
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';

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

export function FormCalendar({
  mode = 'single',
  ...props
}: {
  mode?: 'single' | 'multiple' | 'range';
} & Omit<
  React.ComponentProps<typeof Calendar>,
  'mode' | 'selected' | 'onSelect'
>) {
  const field = useFieldContext<FormCalendarValue>();
  const value = field.state.value;
  const onSelect = (next: FormCalendarValue) => {
    field.handleChange(next);
    field.handleBlur();
  };

  switch (mode) {
    case 'multiple':
      return (
        <Calendar
          {...props}
          mode="multiple"
          selected={Array.isArray(value) ? value : undefined}
          onSelect={onSelect}
        />
      );
    case 'range':
      return (
        <Calendar
          {...props}
          mode="range"
          selected={
            value && !Array.isArray(value) && !(value instanceof Date)
              ? value
              : undefined
          }
          onSelect={onSelect}
        />
      );
    default:
      return (
        <Calendar
          {...props}
          mode="single"
          selected={value instanceof Date ? value : undefined}
          onSelect={onSelect}
        />
      );
  }
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
