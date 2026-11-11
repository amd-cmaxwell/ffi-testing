import { type Point2 } from "./uniffi_testing_types_ffi.ts";
import { type UniffiByteArray, Cursor } from "@ubjs/core";
export declare function crossProduct3d(a: Point3, b: Point3): Point3;
export declare function determinant2x2(m: Matrix2x2): number;
export declare function determinant3x3(m: Matrix3x3): number;
export declare function dotProduct(a: Array<number>, b: Array<number>): number;
export declare function matrixInverse2x2(m: Matrix2x2): Matrix2x2;
export declare function matrixMultiply2x2(a: Matrix2x2, b: Matrix2x2): Matrix2x2;
export declare function solve2x2(a: Matrix2x2, b: Point2): Point2;
export declare function solve3x3(a: Matrix3x3, b: Point3): Point3;
export declare function transpose2x2(m: Matrix2x2): Matrix2x2;
export type Matrix2x2 = {
    m00: number;
    m01: number;
    m10: number;
    m11: number;
};
/**
 * Generated factory for {@link Matrix2x2} record objects.
 */
export declare const Matrix2x2: Readonly<{
    create: (partial: Partial<Matrix2x2> & Required<Omit<Matrix2x2, never>>) => Matrix2x2;
    new: (partial: Partial<Matrix2x2> & Required<Omit<Matrix2x2, never>>) => Matrix2x2;
    defaults: () => Partial<Matrix2x2>;
}>;
export type Matrix3x3 = {
    m00: number;
    m01: number;
    m02: number;
    m10: number;
    m11: number;
    m12: number;
    m20: number;
    m21: number;
    m22: number;
};
/**
 * Generated factory for {@link Matrix3x3} record objects.
 */
export declare const Matrix3x3: Readonly<{
    create: (partial: Partial<Matrix3x3> & Required<Omit<Matrix3x3, never>>) => Matrix3x3;
    new: (partial: Partial<Matrix3x3> & Required<Omit<Matrix3x3, never>>) => Matrix3x3;
    defaults: () => Partial<Matrix3x3>;
}>;
export type Point3 = {
    x: number;
    y: number;
    z: number;
};
/**
 * Generated factory for {@link Point3} record objects.
 */
export declare const Point3: Readonly<{
    create: (partial: Partial<Point3> & Required<Omit<Point3, never>>) => Point3;
    new: (partial: Partial<Point3> & Required<Omit<Point3, never>>) => Point3;
    defaults: () => Partial<Point3>;
}>;
/**
 * This should be called before anything else.
 *
 * It is likely that this is being done for you by the library's `index.ts`.
 *
 * It checks versions of uniffi between when the Rust scaffolding was generated
 * and when the bindings were generated.
 *
 * It also initializes the machinery to enable Rust to talk back to Javascript.
 */
declare function uniffiEnsureInitialized(): void;
declare const _default: Readonly<{
    initialize: typeof uniffiEnsureInitialized;
    converters: {
        FfiConverterTypeMatrix2x2: {
            readFromCursor(c: Cursor): Matrix2x2;
            writeIntoCursor(value: Matrix2x2, c: Cursor): void;
            allocationSize(value: Matrix2x2): number;
            lift(value: UniffiByteArray): Matrix2x2;
            lower(value: Matrix2x2, alloc: import("@ubjs/core").RustBufferAllocator): UniffiByteArray;
        };
        FfiConverterTypeMatrix3x3: {
            readFromCursor(c: Cursor): Matrix3x3;
            writeIntoCursor(value: Matrix3x3, c: Cursor): void;
            allocationSize(value: Matrix3x3): number;
            lift(value: UniffiByteArray): Matrix3x3;
            lower(value: Matrix3x3, alloc: import("@ubjs/core").RustBufferAllocator): UniffiByteArray;
        };
        FfiConverterTypePoint3: {
            readFromCursor(c: Cursor): Point3;
            writeIntoCursor(value: Point3, c: Cursor): void;
            allocationSize(value: Point3): number;
            lift(value: UniffiByteArray): Point3;
            lower(value: Point3, alloc: import("@ubjs/core").RustBufferAllocator): UniffiByteArray;
        };
    };
}>;
export default _default;
