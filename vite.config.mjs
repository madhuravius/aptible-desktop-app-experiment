import react from "@vitejs/plugin-react";
import electron from "vite-plugin-electron";
import tsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vitest/config";

// https://vitejs.dev/config/
export default defineConfig(() => {
  return {
    base: "",
    plugins: [
      electron([
        {
          entry: "electron/main.ts",
        },
        {
          entry: "electron/preload.ts",
          onstart(options) {
            // Notify the Renderer-Process to reload the page when the Preload-Scripts build is complete,
            // instead of restarting the entire Electron App.
            options.reload();
          },
        },
      ]),
      react(),
      tsconfigPaths(),
    ],
  };
});
