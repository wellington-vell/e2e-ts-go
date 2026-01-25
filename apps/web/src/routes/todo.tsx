import { createFileRoute } from "@tanstack/react-router";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

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
          {/* <TodoForm />

          <React.Suspense fallback={<>...loading</>}>
            <TodoList />
          </React.Suspense> */}
        </CardContent>
      </Card>
    </div>
  );
}

// function TodoForm() {
//   const todos = useSuspenseQuery(rpc.v1.todo.getAllTodos);
//   const schema = z.object({
//     text: z
//       .string()
//       .min(2, { message: "Todo must be at least 5 characters long" }),
//   });

//   const mutation = useMutation(rpc.v1.todo.createTodo, {
//     onSuccess: () => {
//       todos.refetch();
//       toast.success("Todo created successfully");
//       form.reset();
//     },
//   });

//   const form = useForm({
//     defaultValues: {
//       text: "",
//     },
//     validators: {
//       onChange: schema,
//     },
//     onSubmit: (values) => {
//       const request = create(CreateTodoRequestSchema, {
//         text: values.value.text,
//       });
//       mutation.mutate(request);
//     },
//   });

//   return (
//     <form
//       onSubmit={(e) => {
//         e.preventDefault();
//         form.handleSubmit();
//       }}
//       className="mb-6 flex items-center space-x-2"
//     >
//       <form.Field name="text">
//         {(field) => (
//           <div className="w-full flex flex-col gap-2">
//             <Input
//               value={field.state.value}
//               onBlur={field.handleBlur}
//               onChange={(e) => field.handleChange(e.target.value)}
//               placeholder="Add a new task..."
//               disabled={mutation.isPending}
//             />

//             {field.state.meta.errors.length > 0 && (
//               <p className="text-red-500">
//                 {field.state.meta.errors
//                   .map((error) => error?.message)
//                   .join(", ")}
//               </p>
//             )}
//           </div>
//         )}
//       </form.Field>

//       <form.Subscribe>
//         {({ canSubmit, isSubmitting }) => (
//           <Button type="submit" disabled={!canSubmit || isSubmitting}>
//             {isSubmitting ? (
//               <Loader2 className="h-4 w-4 animate-spin" />
//             ) : (
//               "Add"
//             )}
//           </Button>
//         )}
//       </form.Subscribe>
//     </form>
//   );
// }

// function TodoList() {
//   const todos = useSuspenseQuery(rpc.v1.todo.getAllTodos);
//   const toggleMutation = useMutation(rpc.v1.todo.toggleTodo, {
//     onSuccess: () => {
//       todos.refetch();
//     },
//   });
//   const deleteMutation = useMutation(rpc.v1.todo.deleteTodo, {
//     onSuccess: () => {
//       todos.refetch();
//     },
//   });

//   const handleToggleTodo = (id: bigint) => {
//     const request = create(ToggleTodoRequestSchema, {
//       id,
//     });
//     toggleMutation.mutate(request);
//   };

//   const handleDeleteTodo = (id: bigint) => {
//     const request = create(DeleteTodoRequestSchema, {
//       id,
//     });
//     deleteMutation.mutate(request);
//   };

//   return (
//     <>
//       {todos.data.todos.length === 0 ? (
//         <p className="py-4 text-center">No todos yet. Add one above!</p>
//       ) : (
//         <ul className="space-y-2">
//           {todos.data.todos.map((todo) => (
//             <li
//               key={Number(todo.id)}
//               className="flex items-center justify-between rounded-md border p-2"
//             >
//               <div className="flex items-center space-x-2">
//                 <Checkbox
//                   checked={todo.completed}
//                   onCheckedChange={() => handleToggleTodo(todo.id)}
//                   id={`todo-${Number(todo.id)}`}
//                 />
//                 <label
//                   htmlFor={`todo-${Number(todo.id)}`}
//                   className={`${todo.completed ? "line-through" : ""}`}
//                 >
//                   {todo.text}
//                 </label>
//               </div>
//               <Button
//                 variant="ghost"
//                 size="icon"
//                 onClick={() => handleDeleteTodo(todo.id)}
//                 aria-label="Delete todo"
//               >
//                 <Trash2 className="h-4 w-4" />
//               </Button>
//             </li>
//           ))}
//         </ul>
//       )}
//     </>
//   );
// }
