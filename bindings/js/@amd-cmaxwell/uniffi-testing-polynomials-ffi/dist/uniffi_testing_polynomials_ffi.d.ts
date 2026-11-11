import { type UniffiByteArray, Cursor } from "@ubjs/core";
export declare function linearRoot(m: number, b: number): number;
export declare function polynomialAdd(a: Array<number>, b: Array<number>): Array<number>;
export declare function polynomialDerivative(coeffs: Array<number>): Array<number>;
export declare function polynomialEvaluate(coeffs: Array<number>, x: number): number;
export declare function polynomialMultiply(a: Array<number>, b: Array<number>): Array<number>;
export declare function quadraticRoots(a: number, b: number, c: number): F64Pair;
export type F64Pair = {
    first: number;
    second: number;
};
/**
 * Generated factory for {@link F64Pair} record objects.
 */
export declare const F64Pair: Readonly<{
    create: (partial: Partial<F64Pair> & Required<Omit<F64Pair, never>>) => F64Pair;
    new: (partial: Partial<F64Pair> & Required<Omit<F64Pair, never>>) => F64Pair;
    defaults: () => Partial<F64Pair>;
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
        FfiConverterTypeF64Pair: {
            readFromCursor(c: Cursor): F64Pair;
            writeIntoCursor(value: F64Pair, c: Cursor): void;
            allocationSize(value: F64Pair): number;
            lift(value: UniffiByteArray): F64Pair;
            lower(value: F64Pair, alloc: import("@ubjs/core").RustBufferAllocator): UniffiByteArray;
        };
    };
}>;
export default _default;
