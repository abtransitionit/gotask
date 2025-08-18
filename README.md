# Gotask


The core orchestration engine for `abtransitionit`'s tools/CLIs. This library manages and runs task pipelines using primitives from `gocore` and `golinux` to execute high-level administrative tasks on Linux Platforms.

It's a library that's built to be more high-level and task-oriented.

---

[![Dev CI](https://github.com/abtransitionit/gotask/actions/workflows/ci-dev.yaml/badge.svg?branch=dev)](https://github.com/abtransitionit/gotask/actions/workflows/ci-dev.yaml)
[![Main CI](https://github.com/abtransitionit/gotask/actions/workflows/ci-main.yaml/badge.svg?branch=main)](https://github.com/abtransitionit/gotask/actions/workflows/ci-main.yaml)
[![LICENSE](https://img.shields.io/badge/license-Apache_2.0-blue.svg)](https://choosealicense.com/licenses/apache-2.0/)

---


# Features  
This project template includes the following components:  


|Component|Description|
|-|-|
|Licensing|Predefined open-source license (Apache 2.0) for legal compliance.|
|Code of Conduct| Ensures a welcoming and inclusive environment for all contributors.|  
|README|Structured documentation template for clear project onboarding.|  

---

## Installation

To use this library in your project, run:

```bash
go get [github.com/abtransitionit/gotask](https://github.com/abtransitionit/gotask)
```

---

# Roadmap


## Build the `gocore` Library

**Goal:** This is your most foundational library. It should contain universal utilities that have no external dependencies on your other repositories. We started this with the `errorx` package.

**Next Action:** Continue to build out `gocore` with other universally useful packages. For example:
* A `logx` package for structured logging.
* A `filex` package for common file system operations.
* A `slicex` package for generic slice helpers.

## Build `golinux` and Connect it to `gocore`

**Goal:** This library will hold your cross-distribution Linux primitives. It will need to use utilities from `gocore`.

**Next Action:** You will start writing the code for packages like `dnfapt` or `oservice`. When you do, you will need to add `gocore` as a dependency. The process will be:
* In the `golinux` repository, run `go get github.com/abtransitionit/gocore`.
* This will add `gocore` to `golinux`'s `go.mod` file, linking your repositories together.

## Build `gotask` and Connect to Both Libraries

**Goal:** This is your orchestration engine. It will need to use primitives from `golinux` and universal utilities from `gocore`.

**Next Action:** You will start writing the core logic of `gotask`. This repository will have two dependencies:
* Run `go get github.com/abtransitionit/gocore`.
* Run `go get github.com/abtransitionit/golinux`.

## Create the First CLI (`kbe`)

**Goal:** The final piece is to create the actual end-user tool. This repository will be a simple Cobra CLI that calls functions from your `gotask` orchestrator.

**Next Action:** Use your `gotplrepo` to create a new repository called `kbe`. It will have a single dependency on `gotask`.

This structured approach ensures you build from the ground up, with each layer depending only on the layers below it.

Where would you like to start? The most logical first step is to continue building out the `gocore` library.

# The goal of code
- The goal of code is not to be as short as possible. The goal is to be correct.
- The goal of code is to be easy for another developer to understand. 
- The number of lines of code is a poor measure of complexity.

# Open-Source Go Alternatives

* **Task:** This is the most popular choice and a direct replacement for tools like `Make`. You define your tasks in a simple `Taskfile.yml` file, and Task handles the execution. It has built-in support for dependencies, caching, and a clean YAML syntax, which is much easier to read and maintain than a long list of function pointers.
* **Mage:** Mage takes a different approach. Instead of using a `YAML` file, you write your build scripts in Go itself. This lets you use Go's full power for more complex tasks, but it also means a higher learning curve.
* **Dagger:** This is a more modern, cloud-native approach. It treats your build logic as a program that can run anywhere—locally, in a container, or on a CI/CD pipeline. It's great for complex, reproducible workflows that involve containers.

For your specific use case, **Task** is a great fit. It would allow you to keep your custom Go code for each **phase** while moving the orchestration logic into a declarative file that's easy to read and manage.

