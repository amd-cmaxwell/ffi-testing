import { type UniffiRustCallStatus } from '@ubjs/core';
interface NativeModuleInterface {
    ffi_uniffi_testing_trig_ffi_rust_future_poll_u8(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_u8(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_u8(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_trig_ffi_rust_future_free_u8(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_i8(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_i8(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_i8(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_trig_ffi_rust_future_free_i8(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_u16(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_u16(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_u16(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_trig_ffi_rust_future_free_u16(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_i16(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_i16(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_i16(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_trig_ffi_rust_future_free_i16(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_u32(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_u32(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_u32(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_trig_ffi_rust_future_free_u32(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_i32(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_i32(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_i32(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_trig_ffi_rust_future_free_i32(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_u64(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_u64(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_u64(handle: bigint, uniffi_out_err: UniffiRustCallStatus): bigint;
    ffi_uniffi_testing_trig_ffi_rust_future_free_u64(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_i64(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_i64(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_i64(handle: bigint, uniffi_out_err: UniffiRustCallStatus): bigint;
    ffi_uniffi_testing_trig_ffi_rust_future_free_i64(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_f32(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_f32(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_f32(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_trig_ffi_rust_future_free_f32(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_f64(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_f64(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_f64(handle: bigint, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_trig_ffi_rust_future_free_f64(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_rust_buffer(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_rust_buffer(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_rust_buffer(handle: bigint, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    ffi_uniffi_testing_trig_ffi_rust_future_free_rust_buffer(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_poll_void(handle: bigint, callback: UniffiRustFutureContinuationCallback, callbackData: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_cancel_void(handle: bigint): void;
    ffi_uniffi_testing_trig_ffi_rust_future_complete_void(handle: bigint, uniffi_out_err: UniffiRustCallStatus): void;
    ffi_uniffi_testing_trig_ffi_rust_future_free_void(handle: bigint): void;
    uniffi_uniffi_testing_trig_ffi_fn_func_circle_area(radius: number, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_trig_ffi_fn_func_circle_circumference(radius: number, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_trig_ffi_fn_func_degrees_to_radians(deg: number, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_trig_ffi_fn_func_distance_2d(x1: number, y1: number, x2: number, y2: number, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_trig_ffi_fn_func_distance_3d(x1: number, y1: number, z1: number, x2: number, y2: number, z2: number, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_trig_ffi_fn_func_law_of_cosines_angle(a: number, b: number, c: number, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_trig_ffi_fn_func_law_of_cosines_side(a: number, b: number, angleC: number, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_trig_ffi_fn_func_law_of_sines_side(a: number, angleA: number, angleB: number, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_trig_ffi_fn_func_midpoint_2d(x1: number, y1: number, x2: number, y2: number, uniffi_out_err: UniffiRustCallStatus): Uint8Array;
    uniffi_uniffi_testing_trig_ffi_fn_func_radians_to_degrees(rad: number, uniffi_out_err: UniffiRustCallStatus): number;
    uniffi_uniffi_testing_trig_ffi_fn_func_triangle_area_heron(a: number, b: number, c: number, uniffi_out_err: UniffiRustCallStatus): number;
    ffi_uniffi_testing_trig_ffi_uniffi_contract_version(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_circle_area(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_circle_circumference(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_degrees_to_radians(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_distance_2d(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_distance_3d(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_cosines_angle(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_cosines_side(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_sines_side(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_midpoint_2d(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_radians_to_degrees(): number;
    uniffi_uniffi_testing_trig_ffi_checksum_func_triangle_area_heron(): number;
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
        ffi_uniffi_testing_trig_ffi_rust_future_poll_u8: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_u8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_u8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "UInt8";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_u8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_i8: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_i8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_i8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Int8";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_i8: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_u16: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_u16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_u16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_u16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_i16: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_i16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_i16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Int16";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_i16: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_u32: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_u32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_u32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "UInt32";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_u32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_i32: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_i32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_i32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Int32";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_i32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_u64: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_u64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_u64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Handle";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_u64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_i64: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_i64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_i64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Int64";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_i64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_f32: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_f32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_f32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Float32";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_f32: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_f64: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_f64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_f64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_f64: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_rust_buffer: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_rust_buffer: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_rust_buffer: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_rust_buffer: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_poll_void: {
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
        ffi_uniffi_testing_trig_ffi_rust_future_cancel_void: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_complete_void: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_rust_future_free_void: {
            args: {
                readonly tag: "Handle";
            }[];
            ret: {
                readonly tag: "Void";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_circle_area: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_circle_circumference: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_degrees_to_radians: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_distance_2d: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_distance_3d: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_law_of_cosines_angle: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_law_of_cosines_side: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_law_of_sines_side: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_midpoint_2d: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "RustBuffer";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_radians_to_degrees: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        uniffi_uniffi_testing_trig_ffi_fn_func_triangle_area_heron: {
            args: {
                readonly tag: "Float64";
            }[];
            ret: {
                readonly tag: "Float64";
            };
            hasRustCallStatus: true;
        };
        ffi_uniffi_testing_trig_ffi_uniffi_contract_version: {
            args: never[];
            ret: {
                readonly tag: "UInt32";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_circle_area: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_circle_circumference: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_degrees_to_radians: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_distance_2d: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_distance_3d: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_cosines_angle: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_cosines_side: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_law_of_sines_side: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_midpoint_2d: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_radians_to_degrees: {
            args: never[];
            ret: {
                readonly tag: "UInt16";
            };
            hasRustCallStatus: false;
        };
        uniffi_uniffi_testing_trig_ffi_checksum_func_triangle_area_heron: {
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
