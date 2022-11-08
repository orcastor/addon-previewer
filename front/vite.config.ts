import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { resolve } from "path";
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers';
import Icons from 'unplugin-icons/vite';
import IconsResolver from 'unplugin-icons/resolver';

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      "@": resolve(__dirname, "./src"),
      // "vue-i18n": "vue-i18n/dist/vue-i18n.cjs.js"
    }
  },
  plugins: [
    vue(),
    AutoImport({
      resolvers: [
        ElementPlusResolver(),
        IconsResolver({
          prefix: 'Icon',
        }),
      ],
    }),
    Components({
      resolvers: [
        IconsResolver({
          enabledCollections: ['ep'],
        }),
        ElementPlusResolver(),
      ],
    }),
    Icons({
      autoInstall: true,
    }),
  ],
  build: {
    outDir: "dist",
    // minify: "esbuild",
    // esbuild 打包更快，但是不能去除 console.log，terser打包慢，但能去除 console.log
    minify: "terser",
    terserOptions: {
      compress: {
     		drop_console: true,
     		drop_debugger: true
     	}
    },
    rollupOptions: {
      output: {
        // Static resource classification and packaging
        chunkFileNames: "prvw/assets/js/[name]-[hash].js",
        entryFileNames: "prvw/assets/js/[name]-[hash].js",
        assetFileNames: "prvw/assets/[ext]/[name]-[hash].[ext]",
        // key自定义 value[] 插件同步package.json名称 或 src/相对路径下的指定文件 （自己可以看manualChunks ts类型）
        manualChunks: {
          // icons: ['@element-plus/icons-vue'],
          pdf: ['vue-pdf-embed'],
        }
      }
    }
  }
})
