import * as React from "react";
import { Events } from "@wailsio/runtime";

export interface AuthContext {
  isAuthenticated: boolean;
}

const AuthContext = React.createContext<AuthContext | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = React.useState(false);
  React.useEffect(() => {
    Events.On("vcalendar-v2:token-needed", (event) => {
      setIsAuthenticated(event.data.TokenNeeded);
    });
  }, [isAuthenticated]);
  return (
    <AuthContext.Provider value={{ isAuthenticated }}>
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
