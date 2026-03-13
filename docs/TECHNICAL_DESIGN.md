# Scaff – Technical Design (v1)

## Purpose

This document describes the **technical architecture** of Scaff and the major components required to support the product requirements defined in the PRD.

The intent is to:

* define system boundaries and responsibilities
* establish extension points for future generators
* outline data structures and interfaces at a conceptual level
* avoid over-prescribing implementation details

The document intentionally avoids describing specific algorithms or writing code-level implementations.

---

# Architectural Principles

Scaff should be built around a few core principles.

### 1. Separation of behavior and structure

Generators describe **what is created**, while architecture styles describe **where it lives**.

This separation ensures that adding new project types, generators, or architecture styles does not require modifying existing logic.

### 2. Generator-centric architecture

All functionality should be implemented through **generators**.

Generators may operate in two modes:

* **project generators**, which create new projects
* **component generators**, which add functionality to existing projects

This unified model allows the system to grow without introducing new architectural concepts.

### 3. Data-driven structure

Project structure should be represented through configuration rather than code wherever possible.

Architecture styles define mappings between logical components and project locations.

This prevents structure rules from being embedded in generator logic.

### 4. Minimal project metadata

Scaff projects should include a small metadata file that allows the tool to detect project type, architecture style, and enabled features.

This metadata allows generators to safely extend projects over time.

### 5. Composable features

Optional integrations such as CI, Makefiles, and release pipelines should be implemented as **feature packs** that can be applied to project types.

---

# High-Level System Architecture

The system can be viewed as four major layers:

```
CLI Interface
    ↓
Generator Engine
    ↓
Structure Resolver
    ↓
Filesystem + Templates
```

Each layer has a distinct responsibility.

---

# CLI Interface

The CLI layer handles user interaction and command parsing.

Its responsibilities include:

* parsing commands and flags
* loading project metadata when operating inside an existing project
* invoking the appropriate generator
* presenting output and errors

The CLI should remain thin and should not contain generation logic itself.

Instead, it acts as an entry point that delegates work to the generator engine.

---

# Generator Engine

The generator engine is the core orchestration component.

It is responsible for:

* discovering available generators
* validating generator compatibility with the current project
* collecting generator inputs
* rendering templates
* coordinating filesystem changes

Generators are the primary extension mechanism for the system.

Two generator categories exist:

### Project Generators

Project generators create new projects from scratch.

Examples:

* package generator
* CLI generator

These generators establish the project structure and create project metadata.

### Component Generators

Component generators operate inside existing projects.

Examples:

* CLI command generator
* route generator (future)
* worker job generator (future)

These generators rely on project metadata to determine how the project is structured.

---

# Structure Resolver

The structure resolver determines **where generated files should be placed** within a project.

Generators specify logical components (for example, command, service, or handler) rather than hardcoded file paths.

The structure resolver maps those logical components to physical locations based on the project’s architecture style.

For example, a command generator may request the location for a `command` component. The resolver then determines the correct path according to the project’s architecture style.

Architecture styles therefore define project layout rules but do not affect generator behavior.

---

# Template System

Scaff uses templates to generate files.

Templates define the structure and content of generated files while allowing variable substitution based on generator inputs.

Templates should support reuse and composability so that different project types and generators can share common template fragments where appropriate.

The template system should remain simple and transparent so that generated code is easy for developers to understand.

---

# Project Metadata

Scaff projects include a small metadata file stored at the root of the project.

This metadata allows the tool to detect that the project was created by Scaff and to determine:

* project type
* architecture style
* enabled features
* generator compatibility

The metadata format should be intentionally minimal to avoid coupling projects tightly to the tool.

Metadata exists primarily to ensure that generators operate safely.

---

# Feature Packs

Optional integrations such as development tooling should be implemented as feature packs.

Feature packs may include:

* Makefiles
* CI configuration
* release pipelines

These features should be applied independently of project generators so that they can be reused across project types.

Feature packs should operate through the same generator engine used by other generators.

---

# Generator Discovery

Generators should be discoverable by the generator engine without requiring manual registration throughout the codebase.

The system should be designed so that new generators can be added by defining them in a predictable location within the codebase.

This approach simplifies future expansion of the tool.

---

# Filesystem Layer

The filesystem layer is responsible for safely applying generator changes.

Its responsibilities include:

* creating directories
* writing files
* detecting conflicts

The filesystem layer should remain isolated from generator logic so that it can handle cross-cutting concerns such as overwrite protection and preview functionality.

---

# Project Detection

When running inside a project, Scaff should attempt to detect whether the project was generated by Scaff.

If metadata is present, the tool can:

* determine project type
* determine architecture style
* validate compatible generators

If metadata is absent, the tool should avoid modifying the project.

---

# Extension Model

The architecture should allow future additions in three main areas:

### New Project Types

Additional project generators may introduce new project types such as:

* web applications
* background workers
* monorepos

### New Component Generators

Project types may introduce additional generators that add functionality within a project.

Examples may include:

* service generators
* handler generators
* job generators

### New Feature Packs

Feature packs may introduce optional integrations such as:

* container configuration
* deployment pipelines
* development tooling

Because these concepts are independent, new capabilities can be introduced incrementally.

---

# Observability and Debugging

Scaff should provide clear output describing what actions generators perform.

Clear logging or reporting of generated files will help developers understand and trust the tool.

---

# Design Tradeoffs

Several deliberate design decisions shape the architecture.

### Avoid embedding architecture logic in generators

Generators should remain simple and portable. Structural decisions belong to architecture style configuration.

### Prefer configuration over code for structure

Architecture styles should be represented through configuration so that layouts can evolve without rewriting generators.

### Maintain transparency of generated code

Generated projects should remain simple and understandable. Scaff should not hide complexity behind heavy abstractions.

---

# Future Evolution

The architecture is designed to support future capabilities such as:

* additional project types
* richer component generators
* more sophisticated project structures
* team-specific template packs

The generator model and structure resolver should allow these features to be added without fundamental changes to the system.

---

# Summary

The technical architecture of Scaff is built around a generator engine that separates behavior, structure, and configuration.

Generators define actions, architecture styles define project layout, and templates define file contents.

This separation ensures that the tool remains flexible and extensible while keeping the implementation understandable.

The resulting system should allow Scaff to grow gradually while maintaining a clear and maintainable codebase.
