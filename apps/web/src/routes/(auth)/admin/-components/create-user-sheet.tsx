import { useMutation } from '@tanstack/react-query';
import { Plus } from 'lucide-react';
import React from 'react';
import { toast } from 'sonner';
import { z } from 'zod';

import { useAppForm } from '@/components/form/hooks';
import { Button } from '@/components/ui/button';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { orpc, queryClient } from '@/lib/orpc';

const schema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.email('Invalid email address'),
  image: z.string(),
  email_verified: z.boolean(),
});

export function CreateUserSheet() {
  const [open, setOpen] = React.useState(false);

  const mutation = useMutation(
    orpc.postAuthAdminUsers.mutationOptions({
      onSuccess: () => {
        toast.success('User created');
        void queryClient.invalidateQueries({
          queryKey: orpc.getAuthAdminUsers.key(),
        });
        setOpen(false);
      },
    }),
  );

  const form = useAppForm({
    defaultValues: {
      name: '',
      email: '',
      image: '',
      email_verified: false,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        body: {
          name: value.name,
          email: value.email,
          email_verified: value.email_verified,
          ...(value.image ? { image: value.image } : {}),
          metadata: [],
        },
      });
    },
    validators: { onSubmit: schema },
  });

  return (
    <>
      <Button onClick={() => setOpen(true)} size="sm">
        <Plus className="size-4 mr-2" />
        Create user
      </Button>
      <Sheet open={open} onOpenChange={setOpen}>
        <SheetContent className="sm:max-w-lg overflow-y-auto">
          <SheetHeader>
            <SheetTitle>Create user</SheetTitle>
            <SheetDescription>
              Add a new user. They won't have a password until you also create
              an account for them.
            </SheetDescription>
          </SheetHeader>
          <form
            onSubmit={async (e) => {
              e.preventDefault();
              e.stopPropagation();
              await form.handleSubmit();
            }}
            className="grid flex-1 auto-rows-min gap-6 px-4"
          >
            <form.AppField name="name">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Name</field.Label>
                  <field.Input placeholder="Jane Doe" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="email">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Email</field.Label>
                  <field.Input type="email" placeholder="jane@example.com" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="image">
              {(field) => (
                <div className="space-y-2">
                  <field.Label>Avatar URL (optional)</field.Label>
                  <field.Input placeholder="https://…" />
                  <field.Error />
                </div>
              )}
            </form.AppField>
            <form.AppField name="email_verified">
              {(field) => (
                <div className="flex items-center gap-2">
                  <field.Checkbox />
                  <field.Label>Mark email as verified</field.Label>
                </div>
              )}
            </form.AppField>
            <form.AppForm>
              <form.Button
                className="w-full"
                loading={mutation.isPending}
                loadingText="Creating..."
              >
                Create user
              </form.Button>
            </form.AppForm>
          </form>
        </SheetContent>
      </Sheet>
    </>
  );
}
