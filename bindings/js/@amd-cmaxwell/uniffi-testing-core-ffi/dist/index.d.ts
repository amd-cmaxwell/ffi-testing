export * from './uniffi_testing_core_ffi.ts';
import * as uniffi_testing_core_ffi from './uniffi_testing_core_ffi.ts';
import { type WasmSource } from '@ubjs/wasm';
export declare function uniffiInitAsync(source: WasmSource): Promise<void>;
declare const _default: {
    uniffi_testing_core_ffi: typeof uniffi_testing_core_ffi;
};
export default _default;
