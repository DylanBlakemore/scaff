Below is a **revised PRD** that is **less prescriptive** and **starts with package creation as the first capability**. It focuses on the **product vision, capabilities, and architecture**, rather than implementation details or specific folder layouts.

---

# PRD: Scaff

## Overview

**Scaff** is a Go-native scaffolding and code-generation tool designed to help developers create and evolve Go projects with consistent structure and conventions.

The tool provides lightweight generators for bootstrapping projects and adding new components to them over time. It is intentionally minimal and Go-idiomatic, focusing on practical workflows rather than framework-style abstraction.

Scaff’s initial capabilities focus on **creating Go packages and CLI tools**, with the architecture designed to support additional project types in the future such as web applications, background workers, and multi-component repositories.

Unlike many scaffolding tools that only run once at project creation, Scaff is intended to remain useful throughout a project’s lifecycle by supporting **generators that add new components inside existing projects**, similar to Phoenix generators.

The goal is to reduce repetitive setup work while keeping projects simple, transparent, and easy to maintain.

---

# Problem

Starting a new Go project typically involves repeated manual steps:

* creating module structures
* establishing directory conventions
* wiring entrypoints
* creating test files
* setting up linting and formatting
* configuring CI and release pipelines

Teams often solve these problems informally through internal templates or documentation, which leads to inconsistent structures and duplicated effort.

Existing tooling partially addresses this problem but tends to fall into one of several categories:

* generic template engines that lack Go-specific conventions
* framework-specific generators that assume a particular stack
* simple project generators that do not support evolving the project afterward

There is currently no widely adopted tool that provides **general-purpose scaffolding for Go projects while remaining lightweight and extensible**.

---

# Product Vision

Scaff should become a small, reliable utility that developers can use whenever they start or expand a Go project.

The tool should:

* provide **sensible starting structures**
* support **optional project tooling integrations**
* allow projects to **grow through structured generators**
* remain **framework-neutral and minimally opinionated**

Scaff should not attempt to impose architecture patterns or replace developer judgment. Instead, it should offer a set of conventions that make common workflows faster and more consistent.

---

# Goals

### Initial Goals

The first version of Scaff should support:

* generating Go packages
* generating Go CLI applications
* optional integration of common project tooling
* the ability to add new components within supported project types

These capabilities establish the core scaffolding engine and generator model that future project types will rely on.

### Long-Term Goals

Over time, Scaff should expand to support additional types of projects and generators, including:

* HTTP applications
* background workers
* multi-component repositories
* project-level tooling generators
* reusable template packs

The architecture should allow these capabilities to be added without significant redesign.

---

# Non-Goals

Scaff is not intended to become a full framework or project management platform.

Specifically, the tool will not initially:

* enforce architectural patterns
* manage dependency upgrades
* automatically modify arbitrary projects
* provide a plugin marketplace
* integrate deeply with non-Go ecosystems

The emphasis should remain on **simplicity, clarity, and developer control**.

---

# Core Concepts

## Project Generators

Project generators create a new project or module from scratch.

In the first version of Scaff, two project types will be supported:

* Go packages
* CLI applications

Each project generator defines a basic structure and any optional tooling integrations that may be enabled at creation time.

Project generators establish the foundation of the project and record minimal metadata so Scaff can recognize the project type later.

---

## Component Generators

Component generators add new functionality within an existing project created by Scaff.

These generators are scoped to a specific project type.

For example, within a CLI project a component generator might add a new command. In future project types, generators could add routes, handlers, jobs, or other project-specific elements.

This capability allows Scaff to remain useful after the initial project creation, enabling incremental project growth through structured generators.

---

## Optional Project Features

When creating a project, users may optionally enable integrations that support common development workflows.

Examples include:

* a Makefile containing common development tasks
* CI integration using GitHub Actions
* automated release pipelines for CLI binaries

These features are intended to reduce the overhead of setting up development workflows while remaining optional and easily understandable.

The generated tooling should follow simple conventions that teams can adjust if needed.

---

# User Experience

Developers interact with Scaff through a command-line interface that exposes two primary categories of operations:

1. **Creating new projects**
2. **Adding new components within existing projects**

Project creation commands initialize a new module and generate the selected project type.

Component commands operate inside an existing project and generate new functionality consistent with the project’s structure.

The interface should remain simple and discoverable, with commands organized around the types of generators available.

---

# Optional Tooling Integrations

Scaff may optionally generate additional development tooling to improve the developer experience.

These integrations should be treated as modular features that can be enabled during project creation.

### Makefile

A Makefile may be generated to provide a consistent interface for common development tasks such as formatting, testing, linting, and running CI checks locally.

The Makefile should include a standard command that performs the same checks that the CI pipeline will run.

### Continuous Integration

CI integration may be generated using GitHub Actions.

The workflow should run standard checks for formatting, tests, and linting, either directly or via the Makefile if one is present.

### Release Automation

For CLI applications, an optional release pipeline may be generated to build and publish binaries when version tags are created.

This capability is primarily intended to simplify distribution of CLI tools.

---

# Architecture Principles

Scaff should be built around a **generator engine** that separates project types, component generators, and optional features.

This separation allows the system to evolve gradually as new project types and generators are added.

The architecture should support three primary concepts:

1. **Project generators**, which create new projects
2. **Component generators**, which extend existing projects
3. **Feature packs**, which add optional tooling and integrations

This model ensures that the addition of new generators or features does not require major changes to the core tool.

Generated projects should include minimal metadata that allows Scaff to identify the project type and available generators. This enables future commands to safely extend the project.

---

# Extensibility

A key design goal for Scaff is the ability to add new generators over time without rewriting the core engine.

The tool should therefore support:

* additional project types
* additional component generators
* reusable feature packs
* template-driven file generation

Examples of potential future generators include:

* web application routes and handlers
* worker jobs and queues
* repository-level tooling
* monorepo workspaces containing multiple binaries

Because these features rely on the same underlying generator model, they can be introduced incrementally.

---

# Success Criteria

The first version of Scaff will be successful if developers can:

* create a Go package with minimal setup
* create a CLI application with a clear starting structure
* optionally enable development tooling such as CI and Makefiles
* generate additional components inside CLI projects
* understand and modify generated code easily

The tool should reduce friction when starting new projects while keeping the generated code straightforward and transparent.

---

# Future Opportunities

If Scaff proves useful, future development may include:

* generators for web services and APIs
* background worker templates
* multi-component repository layouts
* richer component generators for existing project types
* template sharing and distribution mechanisms

The long-term value of Scaff lies in providing a flexible generator framework that grows alongside the Go ecosystem while maintaining a simple and pragmatic developer experience.
