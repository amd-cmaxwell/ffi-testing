import { type Point2 } from "./uniffi_testing_types_ffi.ts";
export declare function circleArea(radius: number): number;
export declare function circleCircumference(radius: number): number;
export declare function degreesToRadians(deg: number): number;
export declare function distance2d(x1: number, y1: number, x2: number, y2: number): number;
export declare function distance3d(x1: number, y1: number, z1: number, x2: number, y2: number, z2: number): number;
export declare function lawOfCosinesAngle(a: number, b: number, c: number): number;
export declare function lawOfCosinesSide(a: number, b: number, angleC: number): number;
export declare function lawOfSinesSide(a: number, angleA: number, angleB: number): number;
export declare function midpoint2d(x1: number, y1: number, x2: number, y2: number): Point2;
export declare function radiansToDegrees(rad: number): number;
export declare function triangleAreaHeron(a: number, b: number, c: number): number;
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
