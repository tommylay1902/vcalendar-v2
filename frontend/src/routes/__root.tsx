import { Outlet, createRootRouteWithContext } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
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
  return (
    <>
      <Outlet />
      {process.env.NODE_ENV === "development" && (
        <TanStackRouterDevtools position="bottom-right" />
      )}
    </>
  );
}
