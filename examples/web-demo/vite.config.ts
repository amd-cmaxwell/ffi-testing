import { defineConfig } from "vite";

export default defineConfig({
  assetsInclude: ["**/*.wasm"],
  optimizeDeps: {
    exclude: [
      "@amd-cmaxwell/uniffi-math-demo-core-ffi",
      "@amd-cmaxwell/uniffi-math-demo-complex-ffi",
      "@amd-cmaxwell/uniffi-math-demo-polynomials-ffi",
      "@amd-cmaxwell/uniffi-math-demo-linear-algebra-ffi",
      "@amd-cmaxwell/uniffi-math-demo-trig-ffi",
      "@amd-cmaxwell/uniffi-math-demo-sequences-ffi",
      "@amd-cmaxwell/uniffi-math-demo-types-ffi",
    ],
  },
});
