module github.com/caarlos0/env/v11

retract v11.0.1 // v11.0.1 accidentally introduced a breaking change regarding the behavior of nil pointers. You can now chose to auto-initialize them by setting the `init` tag option.

go 1.18
