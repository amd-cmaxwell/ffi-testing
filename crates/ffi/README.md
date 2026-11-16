# Shared FFI Types Crates

When you have multiple FFI wrapper crates (e.g. `crates/ffi/foo`,
`crates/ffi/bar`), they sometimes need to share common types — errors,
records, enums — across binding boundaries. A shared types crate gives
those types a single definition that all FFI crates can depend on and
all language bindings can import.

## How to create a shared FFI types crate

1. **Create a new crate** under `crates/ffi/` (e.g. `crates/ffi/my-types`).

2. **Add `uniffi` as a dependency** and call `uniffi::setup_scaffolding!()` in
   `lib.rs`.

3. **Define your shared types** with `#[derive(uniffi::Record)]`,
   `#[derive(uniffi::Enum)]`, `#[derive(uniffi::Error)]`, etc.

4. **Add the crate to your workspace** `Cargo.toml` dependencies.

5. **Depend on it from your other FFI crates** — add it as a dependency in
   each consuming FFI crate's `Cargo.toml`.

6. **Declare the shared crate in each consumer's `uniffi.toml`.**
   Every FFI crate that uses types from a shared types crate **must** have a
   `uniffi.toml` that maps the shared crate name to its Python package name
   under `[bindings.python.external_packages]`. Without this, UniFFI's
   binding generator will try to re-define the shared types inline instead
   of importing them, causing duplicate-type errors at build or runtime.

   For example, if your shared crate is `my_types_ffi`:

   ```toml
   [bindings.python.external_packages]
   my_types_ffi = "my_types_ffi"
   ```

   Each consuming FFI crate needs this declaration. You have two options:

   - **Copy** — place a `uniffi.toml` directly in each FFI crate directory.
   - **Symlink** — maintain a single `uniffi.toml` (e.g. at the project root)
     and symlink to it from each FFI crate directory:

     ```sh
     cd crates/ffi/my-consumer
     ln -s ../../../uniffi.toml uniffi.toml
     ```

   A symlink keeps all crates in sync automatically and avoids having to
   remember to update every copy when you add a new shared types crate.

   > **Warning:** Do **not** add a shared types crate as an external package
   > in its own `uniffi.toml`. The `external_packages` declaration tells
   > UniFFI "these types are defined elsewhere — import them, don't generate
   > them." If the types crate marks itself as external, UniFFI will skip
   > generating the type definitions entirely, and no other crate will be
   > able to import them. Only consuming crates should list the shared types
   > crate in `external_packages`.
