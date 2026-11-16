import { defineConfig } from "vite";

export default defineConfig({
  assetsInclude: ["**/*.wasm"],
  optimizeDeps: {
    exclude: [
      "@%NPM_SCOPE%/uniffi-math-demo-core-ffi",
      "@%NPM_SCOPE%/uniffi-math-demo-complex-ffi",
      "@%NPM_SCOPE%/uniffi-math-demo-polynomials-ffi",
      "@%NPM_SCOPE%/uniffi-math-demo-linear-algebra-ffi",
      "@%NPM_SCOPE%/uniffi-math-demo-trig-ffi",
      "@%NPM_SCOPE%/uniffi-math-demo-sequences-ffi",
      "@%NPM_SCOPE%/uniffi-math-demo-types-ffi",
    ],
  },
});
