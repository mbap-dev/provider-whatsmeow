import { defineConfig } from 'vite'

export default defineConfig({
  base: '/assets/',
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
