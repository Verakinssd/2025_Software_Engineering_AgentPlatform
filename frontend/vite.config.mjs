// Plugins
import Components from 'unplugin-vue-components/vite'
import Vue from '@vitejs/plugin-vue'
import Vuetify, { transformAssetUrls } from 'vite-plugin-vuetify'
import ViteFonts from 'unplugin-fonts/vite'
import VueRouter from 'unplugin-vue-router/vite'
import monacoNS from 'vite-plugin-monaco-editor'

// Utilities
import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'

import Markdown from 'vite-plugin-md'
import markdownItPrism from 'markdown-it-prism'
import markdownItTaskLists from 'markdown-it-task-lists'
import markdownItAnchor from 'markdown-it-anchor'

const monacoPlugin = monacoNS.default ?? monacoNS

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    VueRouter({
      extensions: ['.vue', '.md']
    }),
    Vue({
      include: [/\.vue$/, /\.md$/],
      template: { transformAssetUrls }
    }),
    Markdown({
      wrapperClasses: 'markdown-page',
      markdownItSetup(md) {
        md.use(markdownItPrism)
        md.use(markdownItTaskLists)
        md.use(markdownItAnchor, {
          slugify: s => s.trim().toLowerCase().replace(/\s+/g, '-'),
          permalink: markdownItAnchor.permalink.ariaHidden({ placement: 'before' })
        })
      }
    }),
    Vuetify({
      autoImport: true,
      styles: {
        configFile: 'src/styles/settings.scss',
      },
    }),
    Components(),
    ViteFonts({
      google: {
        families: [{
          name: 'Roboto',
          styles: 'wght@100;300;400;500;700;900',
        }],
      },
    }),
    monacoPlugin({
      languages: [
        'javascript',
        'typescript',
        'json',
        'html',
        'css',
      ]
    })
  ],
  define: { 'process.env': {} },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
    extensions: [
      '.js',
      '.json',
      '.jsx',
      '.mjs',
      '.ts',
      '.tsx',
      '.vue',
      '.md',
    ],
  },
  optimizeDeps: {
    include: [
      '@vue-flow/core',
      '@vue-flow/controls',
      '@vue-flow/background',
      '@vue-flow/minimap'
    ],
    exclude: ['@vue-flow/devtools']
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://114.116.204.215',
        // target: 'http://localhost:8000',
        changeOrigin: true,
        secure: false,
        ws: true,
      },
    }
  },
  css: {
    preprocessorOptions: {
      sass: {
        api: 'modern-compiler',
      },
    },
  },
})