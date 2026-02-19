import { useForm } from "@tanstack/react-form";
import { useMutation, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Loader2, Trash2 } from "lucide-react";
import React from "react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { zPostApiTodosData } from "@/lib/api/zod.gen";
import { orpc } from "@/lib/orpc";

export const Route = createFileRoute("/todo")({
  component: TodosRoute,
});

function TodosRoute() {
  return (
    <div className="mx-auto w-full max-w-md py-10">
      <Card>
        <CardHeader>
          <CardTitle>Todo List</CardTitle>
          <CardDescription>Manage your tasks efficiently</CardDescription>
        </CardHeader>
        <CardContent>
          <TodoForm />

          <React.Suspense fallback={<>...loading</>}>
            <TodoList />
          </React.Suspense>
        </CardContent>
      </Card>
    </div>
  );
}

function TodoForm() {
  const todos = useSuspenseQuery(orpc.getApiTodos.queryOptions());

  const mutation = useMutation(
    orpc.postApiTodos.mutationOptions({
      onSuccess: () => {
        todos.refetch();
        toast.success("Todo created successfully");
        form.reset();
      },
    }),
  );

  const form = useForm({
    defaultValues: {
      text: "",
    },
    validators: {
      onChange: zPostApiTodosData.shape.body,
    },
    onSubmit: (values) => {
      mutation.mutate({
        body: {
          text: values.value.text,
        },
      });
    },
  });

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
      className="mb-6 flex items-center space-x-2"
    >
      <form.Field name="text">
        {(field) => (
          <div className="w-full flex flex-col gap-2">
            <Input
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder="Add a new task..."
              disabled={mutation.isPending}
            />

            {field.state.meta.errors.length > 0 && (
              <p className="text-red-500">
                {field.state.meta.errors
                  .map((error) => error?.message)
                  .join(", ")}
              </p>
            )}
          </div>
        )}
      </form.Field>

      <form.Subscribe>
        {({ canSubmit, isSubmitting }) => (
          <Button type="submit" disabled={!canSubmit || isSubmitting}>
            {isSubmitting ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              "Add"
            )}
          </Button>
        )}
      </form.Subscribe>
    </form>
  );
}

function TodoList() {
  const todos = useSuspenseQuery(orpc.getApiTodos.queryOptions());
  const toggleMutation = useMutation(
    orpc.putApiTodosById.mutationOptions({
      onSuccess: () => {
        toast.success("Todo toggled successfully");
        todos.refetch();
      },
    }),
  );

  const deleteMutation = useMutation(
    orpc.deleteApiTodosById.mutationOptions({
      onSuccess: () => {
        toast.success("Todo deleted successfully");
        todos.refetch();
      },
    }),
  );

  const handleToggleTodo = (id: number) => {
    toggleMutation.mutate({
      body: {
        completed: !todos.data?.body.find((todo) => todo.id === id)?.completed,
      },
      path: {
        id,
      },
    });
  };

  const handleDeleteTodo = (id: number) => {
    deleteMutation.mutate({
      path: {
        id,
      },
    });
  };

  return (
    <>
      {todos.data?.body?.length === 0 ? (
        <p className="py-4 text-center">No todos yet. Add one above!</p>
      ) : (
        <ul className="space-y-2">
          {todos.data?.body?.map((todo) => (
            <li
              key={Number(todo.id)}
              className="flex items-center justify-between rounded-md border p-2"
            >
              <div className="flex items-center space-x-2">
                <Checkbox
                  checked={todo.completed}
                  onCheckedChange={() => handleToggleTodo(todo.id)}
                  id={`todo-${Number(todo.id)}`}
                />
                <label
                  htmlFor={`todo-${Number(todo.id)}`}
                  className={`${todo.completed ? "line-through" : ""}`}
                >
                  {todo.text}
                </label>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => handleDeleteTodo(todo.id)}
                aria-label="Delete todo"
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </li>
          ))}
        </ul>
      )}
    </>
  );
}
