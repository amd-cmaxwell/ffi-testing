import { type UniffiRustCallStatus } from '@ubjs/core';
interface NativeModuleInterface {
    ffi_uniffi_testing_complex_ffi_rust_future_poll_u8(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_u8(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_u8(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_complex_ffi_rust_future_free_u8(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_i8(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_i8(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_i8(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_complex_ffi_rust_future_free_i8(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_u16(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_u16(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_u16(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_complex_ffi_rust_future_free_u16(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_i16(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_i16(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_i16(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_complex_ffi_rust_future_free_i16(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_u32(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_u32(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_u32(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_complex_ffi_rust_future_free_u32(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_i32(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_i32(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_i32(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_complex_ffi_rust_future_free_i32(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_u64(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_u64(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_u64(handle: bigint, uniffi_out_err: UniffiRustCallStatus): bigint;
    ffi_uniffi_testing_complex_ffi_rust_future_free_u64(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_i64(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_i64(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_i64(handle: bigint, uniffi_out_err: UniffiRustCallStatus): bigint;
    ffi_uniffi_testing_complex_ffi_rust_future_free_i64(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_f32(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_f32(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_f32(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_complex_ffi_rust_future_free_f32(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_f64(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_f64(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_f64(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_complex_ffi_rust_future_free_f64(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_rust_buffer(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_rust_buffer(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_rust_buffer(handle: bigint, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    ffi_uniffi_testing_complex_ffi_rust_future_free_rust_buffer(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_poll_void(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_cancel_void(handle: bigint): void;
    ffi_uniffi_testing_complex_ffi_rust_future_complete_void(handle: bigint, uniffi_out_err: UniffiRustCallStatus): void;
    ffi_uniffi_testing_complex_ffi_rust_future_free_void(handle: bigint): void;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_add(a: Uint8Array, b: Uint8Array, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_argument(c: Uint8Array, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_conjugate(c: Uint8Array, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_display(c: Uint8Array, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_div(a: Uint8Array, b: Uint8Array, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_from_polar(r: number, theta: number, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_magnitude(c: Uint8Array, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_mul(a: Uint8Array, b: Uint8Array, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_new(re: number, im: number, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    uniffi_uniffi_testing_complex_ffi_fn_func_complex_sub(a: Uint8Array, b: Uint8Array, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    uniffi_uniffi_testing_complex_ffi_fn_func_quadratic_roots_complex(a: number, b: number, c: number, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    ffi_uniffi_testing_complex_ffi_uniffi_contract_version(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_add(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_argument(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_conjugate(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_display(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_div(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_from_polar(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_magnitude(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_mul(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_new(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_complex_sub(): number;
    uniffi_uniffi_testing_complex_ffi_checksum_func_quadratic_roots_complex(): number;
    rustbuffer_alloc(n: number): Uint8Array;
    rustbuffer_free(view: Uint8Array): void;
}
declare const getter: () => NativeModuleInterface;
export declare const PLAYER_DEFINITIONS: {
    symbols: {
        rustbuffer_alloc: string;
        rustbuffer_free: string;
        rustbuffer_from_bytes: string;
    };
    functions: {
        ffi_uniffi_testing_complex_ffi_rust_future_poll_u8: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_u8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_u8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "UInt8";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_u8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_i8: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_i8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_i8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Int8";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_i8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_u16: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_u16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_u16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_u16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_i16: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_i16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_i16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Int16";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_i16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_u32: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_u32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_u32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "UInt32";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_u32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_i32: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_i32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_i32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Int32";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_i32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_u64: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_u64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_u64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Handle";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_u64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_i64: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_i64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_i64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Int64";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_i64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_f32: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_f32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_f32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Float32";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_f32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_f64: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_f64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_f64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_f64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_rust_buffer: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_rust_buffer: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_rust_buffer: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_rust_buffer: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_poll_void: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Callback";
                readonly name: string;
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_cancel_void: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_complete_void: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_rust_future_free_void: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_add: {
            args: {
                readonly tag: "RustBuffer";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_argument: {
            args: {
                readonly tag: "RustBuffer";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_conjugate: {
            args: {
                readonly tag: "RustBuffer";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_display: {
            args: {
                readonly tag: "RustBuffer";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_div: {
            args: {
                readonly tag: "RustBuffer";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_from_polar: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_magnitude: {
            args: {
                readonly tag: "RustBuffer";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_mul: {
            args: {
                readonly tag: "RustBuffer";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_new: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_complex_sub: {
            args: {
                readonly tag: "RustBuffer";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_complex_ffi_fn_func_quadratic_roots_complex: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_complex_ffi_uniffi_contract_version: {
            args: never[];
            ret: {
                readonly tag: "UInt32";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_add: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_argument: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_conjugate: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_display: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_div: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_from_polar: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_magnitude: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_mul: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_new: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_complex_sub: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_complex_ffi_checksum_func_quadratic_roots_complex: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
    };
    callbacks: {
        RustFutureContinuationCallback: {
            args: ({
                readonly tag: "Handle";
            } | {
                readonly tag: "Int8";
            })[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ForeignFutureDroppedCallback: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
    };
    structs: {
        ForeignFutureDroppedCallbackStruct: ({
            name: string;
            type: {
                readonly tag: "Handle";
            };
        } | {
            name: string;
            type: {
                readonly tag: "Callback";
                readonly name: string;
            };
        })[];
    };
};
export declare function setNativeModule(m: Record<string, (...args: any[]) => any>): void;
export default getter;
export type UniffiRustFutureContinuationCallback = (data: bigint, pollResult: number) => void;
export type UniffiForeignFutureDroppedCallback = (handle: bigint) => void;
export type UniffiForeignFutureDroppedCallbackStruct = {
    handle: bigint;
    free: UniffiForeignFutureDroppedCallback;
};
