import { defineConfig } from "vite";

export default defineConfig({
  assetsInclude: ["**/*.wasm"],
  optimizeDeps: {
    exclude: [
      "@amd-cmaxwell/uniffi-testing-core-ffi",
      "@amd-cmaxwell/uniffi-testing-complex-ffi",
      "@amd-cmaxwell/uniffi-testing-polynomials-ffi",
      "@amd-cmaxwell/uniffi-testing-linear-algebra-ffi",
      "@amd-cmaxwell/uniffi-testing-trig-ffi",
      "@amd-cmaxwell/uniffi-testing-sequences-ffi",
      "@amd-cmaxwell/uniffi-testing-types-ffi",
    ],
  },
});
