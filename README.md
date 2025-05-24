# 🌐 Netshoot

**Netshoot** is a powerful, flexible, and customizable host-checking tool designed to detect tunneling-capable hosts, analyze connection behavior, and provide reliable diagnostics with guaranteed accuracy.

---

## 🚀 Features

### ✅ Guaranteed Host Checking

- Precision engine with accurate TCP, TLS, and application-layer checks.
- Detects partial successes, uncertain states (maybe), and hard failures.
- Handles retries, timeouts, and smart error categorization.

### 🛠️ Highly Customizable

- Define custom payloads (e.g., raw HTTP methods, raw byte sequences).
- Configure TLS settings, timeouts, local bind addresses, and more.
- Supports payload–response validation logic tailored to your use case.

### 🌈 Protocol Flexibility

- Supports raw HTTP, HTTPS/TLS, SOCKS, and other protocols.
- Works alongside VPNs and custom networking setups.
- Allows resume-by-host — pick up from where you left off.

### 📊 Insightful Output

- Structured JSON output for automation and parsing.
- Includes per-payload speed, TLS info, and error summaries.
- Live progress tracking with real-time status files.

---

## 🧩 Use Cases

- Detecting **tunneling-capable hosts**
- Identifying **tunnelable methods**
- Measuring and comparing **maximum host speeds**

---

## 📦 Installation

To build from source:

```bash
git clone https://github.com/sadeepa24/netshoot.git
cd netshoot
go build -o netshoot
```

(Ensure you have Go 1.17+ installed.)

---

## 🛠 Example Usage

```bash
./netshoot run
```

---

## 🧾 DOCS

Read Docs to Get Better Understand about Netshoot

[Docs](https://sadeepa24.github.io/netshoot/)

---

## 🧾 License

MIT License. © 2025 sadeepa24

---

## 🙋‍♂️ Contributing

Issues and PRs are welcome! Please open an issue to discuss major changes before submitting a pull request.
