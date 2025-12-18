import React, { useEffect } from "react";
import ReactDOM from "react-dom/client";
import {
  RouterProvider,
  createRouter,
  useNavigate,
} from "@tanstack/react-router";
import "./styles.css";
import { routeTree } from "./routeTree.gen";
import { AuthProvider, useAuth } from "./auth";

// Define the router context type matching your AuthContext
export interface RouterContext {
  auth: {
    isAuthenticated: boolean;
  };
}

// Create the router with the proper context structure
const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  scrollRestoration: true,
  context: {
    // This will be overridden when we render, but needs the right type
    auth: undefined!,
  },
});

// Register for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

function InnerApp() {
  const auth = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
}

function App() {
  return (
    <AuthProvider>
      <InnerApp />
    </AuthProvider>
  );
}

const rootElement = document.getElementById("root");
if (rootElement && !rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  );
}
