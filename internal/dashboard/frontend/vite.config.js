import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// The Go binary embeds internal/dashboard/public/* and serves index.html at "/"
// with hashed-free asset names under /assets/. Keep the output names stable
// (app.js / styles.css) and base "/assets/" so static.go and index.html match.
export default defineConfig({
  base: '/assets/',
  plugins: [react(), tailwindcss()],
  build: {
    outDir: '../public',
    emptyOutDir: true,
    assetsDir: '.',
    cssCodeSplit: false,
    rollupOptions: {
      output: {
        entryFileNames: 'app.js',
        chunkFileNames: '[name].js',
        assetFileNames: (assetInfo) =>
          assetInfo.name && assetInfo.name.endsWith('.css')
            ? 'styles.css'
            : '[name][extname]',
      },
    },
  },
})
