import { useAuth } from "@/auth";
import {
  Outlet,
  createRootRouteWithContext,
  useNavigate,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { useEffect } from "react";
// Create root route with context that matches what we pass to RouterProvider
export const Route = createRootRouteWithContext<{
  auth: {
    isAuthenticated: boolean;
    user: string | null;
    login: (username: string) => Promise<void>;
    logout: () => Promise<void>;
  };
}>()({
  component: RootComponent,
});

function RootComponent() {
  const auth = useAuth();
  const navigate = useNavigate();
  useEffect(() => {
    if (!auth.isAuthenticated) {
      if (window.location.pathname != "/google-auth") {
        navigate({ to: "/google-auth" });
      }
    }
  });

  return (
    <>
      <Outlet />
      {process.env.NODE_ENV === "development" && (
        <TanStackRouterDevtools position="bottom-right" />
      )}
    </>
  );
}
