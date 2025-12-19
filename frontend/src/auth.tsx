import * as React from "react";
import { Events, WML } from "@wailsio/runtime";
import { Navigate } from "@tanstack/react-router";

export interface AuthContext {
  isAuthenticated: boolean;
  isLoading: boolean;
}

const AuthContext = React.createContext<AuthContext | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = React.useState(false);
  const [isLoading, setIsLoading] = React.useState(true);
  React.useEffect(() => {
    // Event handler for token-needed event
    const handleTokenNeeded = (event: any) => {
      console.log("Token event received:", event);
      setIsAuthenticated(event.data.TokenNeeded);
      setIsLoading(false);
    };

    Events.On("vcalendar-v2:token-needed", handleTokenNeeded);
    console.log("AuthProvider mounted, event listeners set up");
    Events.Emit("vcalendar-v2:auth-needed");
    return () => {
      Events.Off("vcalendar-v2:token-needed", handleTokenNeeded as any);
    };
  }, []); // Empty dependency array - run once on mount

  // Add this to debug
  React.useEffect(() => {
    console.log("isAuthenticated updated:", isAuthenticated);
  }, [isAuthenticated, isLoading]);

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = React.useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
