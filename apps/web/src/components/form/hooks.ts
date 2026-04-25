import { createFormHook } from '@tanstack/react-form';

import {
  FormButton,
  FormCheckbox,
  FormError,
  FormField,
  FormInput,
  FormInputGroupInput,
  FormInputOTP,
  FormLabel,
  FormSelect,
  FormSelectContent,
  FormSelectItem,
  FormSelectTrigger,
  FormTextarea,
} from '@/components/form';
import { fieldContext, formContext } from '@/components/form/context';

export const { useAppForm } = createFormHook({
  fieldComponents: {
    Input: FormInput,
    Label: FormLabel,
    Field: FormField,
    Error: FormError,
    Textarea: FormTextarea,
    Checkbox: FormCheckbox,
    InputOTP: FormInputOTP,
    InputGroupInput: FormInputGroupInput,
    Select: FormSelect,
    SelectTrigger: FormSelectTrigger,
    SelectContent: FormSelectContent,
    SelectItem: FormSelectItem,
  },
  formComponents: {
    Button: FormButton,
  },
  fieldContext,
  formContext,
});
