export declare function arithmeticSum(a1: number, d: number, n: bigint): number;
export declare function arithmeticTerm(a1: number, d: number, n: bigint): number;
export declare function combinations(n: bigint, r: bigint): bigint;
export declare function factorial(n: bigint): bigint;
export declare function fibonacci(n: bigint): bigint;
export declare function gcd(a: bigint, b: bigint): bigint;
export declare function geometricSum(a1: number, r: number, n: bigint): number;
export declare function geometricSumInfinite(a1: number, r: number): number;
export declare function geometricTerm(a1: number, r: number, n: bigint): number;
export declare function isPrime(n: bigint): boolean;
export declare function lcm(a: bigint, b: bigint): bigint;
export declare function permutations(n: bigint, r: bigint): bigint;
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
}>;
export default _default;
