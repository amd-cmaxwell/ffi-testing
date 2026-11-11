import { type UniffiByteArray, Cursor, uniffiTypeNameSymbol } from "@ubjs/core";
export type Point2 = {
    x: number;
    y: number;
};
/**
 * Generated factory for {@link Point2} record objects.
 */
export declare const Point2: Readonly<{
    create: (partial: Partial<Point2> & Required<Omit<Point2, never>>) => Point2;
    new: (partial: Partial<Point2> & Required<Omit<Point2, never>>) => Point2;
    defaults: () => Partial<Point2>;
}>;
export declare enum MathError_Tags {
    DivisionByZero = "DivisionByZero",
    InvalidInput = "InvalidInput",
    NoRealSolution = "NoRealSolution",
    Overflow = "Overflow"
}
export declare const MathError: Readonly<{
    instanceOf: (obj: any) => obj is MathError;
    DivisionByZero: {
        new (): {
            readonly tag: MathError_Tags.DivisionByZero;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        "new"(): {
            readonly tag: MathError_Tags.DivisionByZero;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        instanceOf(obj: any): obj is {
            readonly tag: MathError_Tags.DivisionByZero;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        hasInner(obj: any): obj is {
            readonly tag: MathError_Tags.DivisionByZero;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
    };
    InvalidInput: {
        new (inner: {
            message: string;
        }): {
            readonly tag: MathError_Tags.InvalidInput;
            readonly inner: Readonly<{
                message: string;
            }>;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        "new"(inner: {
            message: string;
        }): {
            readonly tag: MathError_Tags.InvalidInput;
            readonly inner: Readonly<{
                message: string;
            }>;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        instanceOf(obj: any): obj is {
            readonly tag: MathError_Tags.InvalidInput;
            readonly inner: Readonly<{
                message: string;
            }>;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        hasInner(obj: any): obj is {
            readonly tag: MathError_Tags.InvalidInput;
            readonly inner: Readonly<{
                message: string;
            }>;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        getInner(obj: {
            readonly tag: MathError_Tags.InvalidInput;
            readonly inner: Readonly<{
                message: string;
            }>;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        }): Readonly<{
            message: string;
        }>;
    };
    NoRealSolution: {
        new (): {
            readonly tag: MathError_Tags.NoRealSolution;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        "new"(): {
            readonly tag: MathError_Tags.NoRealSolution;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        instanceOf(obj: any): obj is {
            readonly tag: MathError_Tags.NoRealSolution;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        hasInner(obj: any): obj is {
            readonly tag: MathError_Tags.NoRealSolution;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
    };
    Overflow: {
        new (): {
            readonly tag: MathError_Tags.Overflow;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        "new"(): {
            readonly tag: MathError_Tags.Overflow;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        instanceOf(obj: any): obj is {
            readonly tag: MathError_Tags.Overflow;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
        hasInner(obj: any): obj is {
            readonly tag: MathError_Tags.Overflow;
            /**
             * @private
             * This field is private and should not be used, use `tag` instead.
             */
            readonly [uniffiTypeNameSymbol]: "MathError";
            name: string;
            message: string;
            stack?: string;
            cause?: unknown;
        };
    };
}>;
export type MathError = InstanceType<typeof MathError['DivisionByZero' | 'InvalidInput' | 'NoRealSolution' | 'Overflow']>;
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
        FfiConverterTypeMathError: {
            readFromCursor(c: Cursor): MathError;
            writeIntoCursor(value: MathError, c: Cursor): void;
            allocationSize(value: MathError): number;
            lift(value: UniffiByteArray): MathError;
            lower(value: MathError, alloc: import("@ubjs/core").RustBufferAllocator): UniffiByteArray;
        };
        FfiConverterTypePoint2: {
            readFromCursor(c: Cursor): Point2;
            writeIntoCursor(value: Point2, c: Cursor): void;
            allocationSize(value: Point2): number;
            lift(value: UniffiByteArray): Point2;
            lower(value: Point2, alloc: import("@ubjs/core").RustBufferAllocator): UniffiByteArray;
        };
    };
}>;
export default _default;
