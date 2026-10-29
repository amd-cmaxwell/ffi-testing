> [!NOTE]
> This directory is intentionally left empty to be populated by the Maturin build backend while building the package.
> Once built, this directory will contain compiled Rust library code as a shared object (.so) and a UniFFI-framework specific python script which facilitates calling into that shared object. 

> [!CAUTION]
> Do not commit anything to this directory. Doing so breaks git-based pip installations as git detects the build system overwritting tracked files.