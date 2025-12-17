import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import wails from "@wailsio/runtime/plugins/vite";
import path from "path";
import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
export default defineConfig({
  plugins: [
    tanstackRouter({ target: "react", autoCodeSplitting: true }),
    react(),
    wails("./bindings"),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      bindings: path.resolve(__dirname, "./bindings"),
    },
  },
});
