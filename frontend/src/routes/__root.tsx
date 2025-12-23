import { useAuth } from "@/auth";
import {
  Outlet,
  createRootRouteWithContext,
  useNavigate,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { useEffect } from "react";

export const Route = createRootRouteWithContext<{
  auth: {
    isAuthenticated: boolean;
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
    } else {
      navigate({ to: "/" });
    }
  }, [auth]);

  return (
    <>
      <Outlet />
      {process.env.NODE_ENV === "development" && (
        <TanStackRouterDevtools position="bottom-right" />
      )}
    </>
  );
}
