import { resolve } from 'path';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';

export default defineConfig({
  plugins: [react()],
  define: {
    'process.env.NODE_ENV': '"production"',
  },
  build: {
    lib: {
      entry: resolve(__dirname, 'src/main.tsx'),
      name: 'MailyBuilder',
      formats: ['umd'],
      fileName: (format) => `maily-builder.${format}.js`,
    },
    minify: 'esbuild',
    cssCodeSplit: false,
    rollupOptions: {
      output: {
        assetFileNames: 'maily-builder.[ext]',
      },
    },
  },
});
