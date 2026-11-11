export * from './uniffi_testing_polynomials_ffi.ts';
export * from './uniffi_testing_types_ffi.ts';
import * as uniffi_testing_polynomials_ffi from './uniffi_testing_polynomials_ffi.ts';
import * as uniffi_testing_types_ffi from './uniffi_testing_types_ffi.ts';
import { type WasmSource } from '@ubjs/wasm';
export declare function uniffiInitAsync(source: WasmSource): Promise<void>;
declare const _default: {
    uniffi_testing_polynomials_ffi: typeof uniffi_testing_polynomials_ffi;
    uniffi_testing_types_ffi: typeof uniffi_testing_types_ffi;
};
export default _default;
