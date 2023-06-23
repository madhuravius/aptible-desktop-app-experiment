import react from "@vitejs/plugin-react";
import electron from "vite-plugin-electron";
import tsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vitest/config";

// https://vitejs.dev/config/
export default defineConfig(() => {
  return {
    plugins: [
      electron({
        entry: "electron/main.ts",
      }),
      react(),
      tsconfigPaths(),
    ],
  };
});
