module github.com/caarlos0/env/v11

// v11.0.1 accidentally introduced a breaking change regarding the behavior of nil pointers.
// You can now chose to auto-initialize them by setting the `init` tag option.
retract v11.0.1

// v11.2.0 accidentally introduced a breaking change regarding the behavior of nil slices of complex types.
retract v11.2.0

// v11.3.0 merges OS environment variables with environments set with Options instead of overriding them.
retract v11.3.0

go 1.18
