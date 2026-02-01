> **Role**
> You are a senior Go engineer and developer tooling expert.
> You are building a production-ready, cross-platform Go CLI tool called **DepGuard**.

---

### 🎯 Goal

Build **DepGuard**, a Go-based CLI tool that prevents developers from adding new dependencies to `package.json` unless they can explain at least **40% of the existing dependencies**.

The tool must be:

* Fast
* Deterministic
* Offline-first
* CI-friendly
* Usable via a single static binary

No Node.js runtime is allowed.

---

### 🧱 Core requirements

#### 1. CLI commands (mandatory)

Implement the following commands:

```
depguard init
depguard verify
depguard explain <dependency-name>
depguard audit
```

Use **cobra** for CLI structure.

---

#### 2. `depguard init`

* Read `package.json`
* Collect `dependencies` and `devDependencies`
* Create `.depguard/`
* Save a baseline snapshot in `.depguard/baseline.json`
* Create `.depguard/knowledge.json` (empty if not present)
* Print next-step instructions to add:

  ```json
  "preinstall": "depguard verify"
  ```

---

#### 3. Dependency parsing

Parse `package.json` using Go structs.

Merge `dependencies` and `devDependencies` into a single dependency set for evaluation.

Ignore:

* peerDependencies
* optionalDependencies

---

#### 4. Baseline comparison

On `verify`:

* Load `.depguard/baseline.json`
* Compare with current `package.json`
* Detect **new dependencies**
* If no new dependencies → exit success immediately

---

#### 5. 40% deterministic selection

* When new dependencies are detected:

  * Randomly select **40%** of *all* dependencies
  * Selection must be deterministic
  * Seed randomness using:

    * CI mode → commit SHA (from env)
    * Local mode → timestamp

---

#### 6. Knowledge enforcement

Each selected dependency must:

* Exist in `.depguard/knowledge.json`
* Have a non-empty explanation
* Not be expired

Knowledge entries must contain:

```json
{
  "summary": "string",
  "explained_by": "string",
  "expires_at": "RFC3339 timestamp"
}
```

Expiration default: **90 days**

---

#### 7. `depguard explain <dep>`

* Interactive terminal prompt
* Ask:

  * What does this dependency do?
  * Where is it used?
  * Could it be removed?
* Save explanation
* Set expiration to now + 90 days

Use `bubbletea` for TUI.

---

#### 8. Blocking behavior

* If verification fails:

  * Print clear, human-readable error
  * Exit with non-zero status
* Must be suitable for npm `preinstall` hook

---

#### 9. `depguard audit`

Print:

* Total dependencies
* Explained dependencies
* Unexplained dependencies
* List top unexplained dependencies

---

### 🗂️ Project structure (mandatory)

```
depguard/
├── cmd/
│   ├── root.go
│   ├── init.go
│   ├── verify.go
│   ├── explain.go
│   └── audit.go
├── internal/
│   ├── parser/
│   ├── baseline/
│   ├── knowledge/
│   ├── selector/
│   └── util/
├── main.go
```

No global state. Clear separation of concerns.

---

### ⚙️ Non-functional constraints

* No SaaS
* No telemetry
* No network calls by default
* Must run in <100ms on average projects
* Must compile into a single static binary

---

### 🧪 Output expectations

You must:

* Generate all Go source files
* Include example JSON files
* Include a README explaining usage
* Include clear error messages
* Avoid placeholder TODOs

---

### 🧠 Style & philosophy

* Be opinionated
* Prefer clarity over cleverness
* Fail loudly and helpfully
* Assume users are competent but lazy
* Do not shame in output text (humor is optional but subtle)

---

### 🚫 Explicitly forbidden

* Writing any JavaScript
* Requiring Node.js
* Auto-installing dependencies
* Cloud services
* Asking the user for design clarification

---

### ✅ Completion definition

The task is complete when:

* `depguard init`
* `depguard verify`
* `depguard explain`
* `depguard audit`

all work together locally on a real Node.js project without modification.
