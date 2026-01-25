import { RouterProvider, createRouter } from "@tanstack/react-router";
import React from "react";
import ReactDOM from "react-dom/client";

import { Loader } from "@/components/loader";
import { routeTree } from "@/routeTree.gen";

const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  defaultPendingComponent: () => <Loader />,
  // context: {
  //   transport,
  //   queryClient,
  // },
  // Wrap: function WrapComponent({ children }: { children: React.ReactNode }) {
  //   return (
  //     <TransportProvider transport={transport}>
  //       <QueryClientProvider client={queryClient}>
  //         {children}
  //       </QueryClientProvider>
  //     </TransportProvider>
  //   );
  // },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const rootElement = document.getElementById("app");

if (!rootElement) {
  throw new Error("Root element not found");
}

if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <React.StrictMode>
      <RouterProvider router={router} />
    </React.StrictMode>,
  );
}
