import { type UniffiByteArray, Cursor } from "@ubjs/core";
export declare function complexAdd(a: Complex, b: Complex): Complex;
export declare function complexArgument(c: Complex): number;
export declare function complexConjugate(c: Complex): Complex;
export declare function complexDisplay(c: Complex): string;
export declare function complexDiv(a: Complex, b: Complex): Complex;
export declare function complexFromPolar(r: number, theta: number): Complex;
export declare function complexMagnitude(c: Complex): number;
export declare function complexMul(a: Complex, b: Complex): Complex;
export declare function complexNew(re: number, im: number): Complex;
export declare function complexSub(a: Complex, b: Complex): Complex;
export declare function quadraticRootsComplex(a: number, b: number, c: number): ComplexPair;
export type Complex = {
    re: number;
    im: number;
};
/**
 * Generated factory for {@link Complex} record objects.
 */
export declare const Complex: Readonly<{
    create: (partial: Partial<Complex> & Required<Omit<Complex, never>>) => Complex;
    new: (partial: Partial<Complex> & Required<Omit<Complex, never>>) => Complex;
    defaults: () => Partial<Complex>;
}>;
export type ComplexPair = {
    first: Complex;
    second: Complex;
};
/**
 * Generated factory for {@link ComplexPair} record objects.
 */
export declare const ComplexPair: Readonly<{
    create: (partial: Partial<ComplexPair> & Required<Omit<ComplexPair, never>>) => ComplexPair;
    new: (partial: Partial<ComplexPair> & Required<Omit<ComplexPair, never>>) => ComplexPair;
    defaults: () => Partial<ComplexPair>;
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
        FfiConverterTypeComplex: {
            readFromCursor(c: Cursor): Complex;
            writeIntoCursor(value: Complex, c: Cursor): void;
            allocationSize(value: Complex): number;
            lift(value: UniffiByteArray): Complex;
            lower(value: Complex, alloc: import("@ubjs/core").RustBufferAllocator): UniffiByteArray;
        };
        FfiConverterTypeComplexPair: {
            readFromCursor(c: Cursor): ComplexPair;
            writeIntoCursor(value: ComplexPair, c: Cursor): void;
            allocationSize(value: ComplexPair): number;
            lift(value: UniffiByteArray): ComplexPair;
            lower(value: ComplexPair, alloc: import("@ubjs/core").RustBufferAllocator): UniffiByteArray;
        };
    };
}>;
export default _default;
