module github.com/caarlos0/env/v11

retract v11.0.1 // v11.0.1 accidentally introduced a breaking change regarding the behavior of uninitalized pointers. You can now chose to auto-innit nil pointers by setting the 'init' tag option.

go 1.18
