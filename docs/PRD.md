# PRD: Scaff

## Overview

**Scaff** is a Go-native scaffolding and code generation tool designed to simplify the creation and evolution of Go projects.

The tool helps developers quickly bootstrap new modules and add structured components to them over time, using lightweight conventions and generators.

Scaff focuses on reducing repetitive setup work while keeping generated projects transparent, idiomatic, and easy to modify.

The initial capabilities focus on generating **Go packages** and **CLI applications**, while the underlying architecture is designed to support additional project types in the future such as:

* web applications
* background workers
* multi-binary repositories
* monorepos

A key design principle is that Scaff should remain useful beyond the initial project creation. Inspired by Phoenix-style generators, Scaff should support adding new components within an existing project using structured generators.

---

# Problem

Creating new Go projects typically requires repeating many setup steps:

* initializing modules
* establishing directory conventions
* creating test scaffolding
* wiring entrypoints
* setting up formatting, linting, and testing workflows
* configuring CI and release automation

Teams often solve these problems through documentation or copy-paste templates, which leads to inconsistent structures and duplicated effort.

Existing tools only partially address this problem. Some are too generic and lack Go-specific conventions, while others assume specific frameworks or only generate projects once.

There is currently no widely adopted tool that provides **lightweight, extensible scaffolding for general Go development workflows**.

---

# Product Vision

Scaff should become a small, reliable utility that Go developers use whenever they start or expand a project.

The tool should:

* create consistent project structures
* support optional development tooling integrations
* allow projects to evolve through structured generators
* remain minimally opinionated and framework-neutral

Scaff should prioritize **clarity and simplicity** over abstraction. Generated code should be easy for developers to understand and modify.

---

# Goals

## Initial Goals

The first version of Scaff should support:

* generating Go packages
* generating Go CLI applications
* optional integration of development tooling such as Makefiles and CI
* optional release automation for CLI tools
* generators that add components inside supported projects

These capabilities establish the core generator engine and architectural model that future project types will rely on.

## Long-Term Goals

Future capabilities may include generators for:

* HTTP applications
* background workers
* service-oriented architectures
* monorepos containing multiple binaries
* project-level tooling integrations
* reusable template packs

The architecture should allow these capabilities to be introduced incrementally.

---

# Non-Goals

Scaff is not intended to become a framework or project management platform.

In particular, the tool will not initially:

* enforce architectural patterns
* manage dependency upgrades
* modify arbitrary non-Scaff projects
* integrate deeply with frontend ecosystems
* implement a plugin marketplace

The focus should remain on lightweight scaffolding and structured code generation.

---

# Core Concepts

## Project Generators

Project generators create a new project or module from scratch.

In the first version of Scaff, two project types will be supported:

* Go packages
* CLI applications

Each project generator defines an initial project structure and may optionally include additional tooling integrations.

Generated projects include minimal metadata that allows Scaff to detect the project type later and enable appropriate generators.

---

## Component Generators

Component generators add new functionality inside an existing Scaff-generated project.

These generators allow projects to evolve over time.

Examples include:

* adding a CLI command to a CLI project
* adding routes to a web application
* adding jobs to a worker system

Each project type defines the component generators that apply to it.

This design allows Scaff to remain useful after project creation rather than acting as a one-time generator.

---

## Optional Project Features

Scaff may optionally generate supporting tooling to simplify common development workflows.

These integrations may include:

* Makefiles for common development tasks
* CI integration using GitHub Actions
* automated release pipelines for CLI binaries

These features are optional and should remain easy to understand and modify.

---

# Architecture Styles

Some project types may support **architecture style presets** that influence how generated code is organized within the project.

Architecture styles determine directory layout and generator behavior but do not impose specific frameworks or implementation patterns.

These styles exist to provide sensible starting structures and reduce early project design friction.

Initial architecture styles may include:

* **minimal**, representing simple and shallow structures suitable for small projects
* **layered**, representing separation between transport, logic, and persistence layers
* **domain**, representing organization by domain capability rather than technical layer

Architecture styles influence where generated components are placed within a project but should not constrain how developers implement their code.

---

# Generator Model

Scaff is built around a generator engine that separates three key concepts:

1. **Project types**, which define the overall project structure and supported generators
2. **Generators**, which define actions that create or modify project components
3. **Architecture styles**, which determine where generated code is placed

This separation allows new project types, generators, and architecture styles to be introduced without significant changes to the core engine.

---

# Strategy and Data-Driven Design

The internal architecture of Scaff should distinguish between **behavioral logic** and **structural configuration**.

### Strategy-Based Behavior

Some aspects of the system represent distinct behaviors and should be implemented using a strategy-style model.

These include:

* project types
* generators
* feature integrations

Each of these represents an action or capability with specific logic.

### Data-Driven Structure

Other aspects represent structural configuration and should be modeled as data rather than code.

In particular, architecture styles should be defined through configuration that maps logical components to project locations.

Generators should describe **what they create**, while the project structure determines **where those components are placed**.

This separation ensures that architecture styles can evolve independently from generators.

---

# Optional Development Tooling

Scaff may optionally generate project tooling that supports common development workflows.

These features should be modular and reusable across project types.

### Makefile

A Makefile may be generated to provide a consistent interface for development tasks such as formatting, testing, and linting.

The Makefile should expose a unified command used by CI pipelines.

### Continuous Integration

CI integration may be generated using GitHub Actions.

The workflow should run formatting checks, tests, and linting, either directly or via the Makefile if present.

### Release Automation

CLI projects may optionally include a release pipeline that builds and publishes binaries when version tags are created.

This capability simplifies distributing CLI tools.

---

# Extensibility

A core design goal of Scaff is extensibility.

The generator engine should allow new capabilities to be added with minimal impact on existing functionality.

Potential extensions include:

* web application generators
* worker/job generators
* monorepo layouts containing multiple binaries
* additional development tooling integrations
* additional architecture styles

Because project types, generators, and architecture styles are independent concepts, new capabilities can be introduced incrementally.

---

# Developer Experience

Scaff should provide a simple and discoverable command-line interface that organizes operations around generators.

Developers should be able to:

* create new projects
* inspect available generators
* add components within existing projects
* enable optional tooling integrations

The CLI should remain intuitive and predictable so developers can easily understand what will be generated.

---

# Success Criteria

The first version of Scaff will be successful if developers can:

* quickly generate a Go package with a working module and test scaffold
* create a CLI application with a clean entrypoint
* optionally enable Makefile, CI, and release integrations
* add new commands to a CLI project using a generator
* understand and modify generated code without difficulty

The generated code should remain simple, idiomatic, and easy to maintain.

---

# Long-Term Opportunity

If Scaff proves useful, it could evolve into a standard scaffolding tool for Go development workflows.

Future expansions may include:

* richer generators for service-oriented applications
* generators for worker systems
* project structures for multi-service repositories
* reusable template packs for teams
* structured project upgrades and migrations

The long-term value of Scaff lies in providing a flexible generator system that grows with the needs of Go developers while maintaining a simple and pragmatic user experience.
