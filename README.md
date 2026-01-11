<p align="center">
  <a href="https://www.renamed.to">
    <img src="https://www.renamed.to/logo.svg" alt="renamed.to" width="120" />
  </a>
</p>

<h1 align="center">renamed.to SDK</h1>

<p align="center">
  <strong>Official SDKs for AI-powered file renaming, PDF splitting, and data extraction</strong>
</p>

<p align="center">
  <a href="https://www.npmjs.com/package/@renamed-to/sdk"><img src="https://img.shields.io/npm/v/@renamed-to/sdk?style=flat-square&label=TypeScript&color=3178c6" alt="TypeScript" /></a>
  <a href="https://pypi.org/project/renamed/"><img src="https://img.shields.io/pypi/v/renamed?style=flat-square&label=Python&color=3776ab" alt="Python" /></a>
  <a href="https://pkg.go.dev/github.com/renamed-to/renamed-sdk/sdks/go"><img src="https://img.shields.io/github/v/tag/renamed-to/renamed-sdk?style=flat-square&label=Go&color=00add8" alt="Go" /></a>
  <a href="https://central.sonatype.com/artifact/to.renamed/renamed-sdk"><img src="https://img.shields.io/maven-central/v/to.renamed/renamed-sdk?style=flat-square&label=Java&color=ed8b00" alt="Java" /></a>
  <a href="https://www.nuget.org/packages/Renamed.Sdk"><img src="https://img.shields.io/nuget/v/Renamed.Sdk?style=flat-square&label=C%23&color=512bd4" alt="C#" /></a>
</p>

<p align="center">
  <a href="https://rubygems.org/gems/renamed"><img src="https://img.shields.io/gem/v/renamed?style=flat-square&label=Ruby&color=cc342d" alt="Ruby" /></a>
  <a href="https://crates.io/crates/renamed"><img src="https://img.shields.io/crates/v/renamed?style=flat-square&label=Rust&color=dea584" alt="Rust" /></a>
  <a href="https://github.com/renamed-to/renamed-sdk"><img src="https://img.shields.io/github/v/tag/renamed-to/renamed-sdk?style=flat-square&label=Swift&color=f05138" alt="Swift" /></a>
  <a href="https://packagist.org/packages/renamed-to/renamed-php"><img src="https://img.shields.io/packagist/v/renamed-to/renamed-php?style=flat-square&label=PHP&color=777bb4" alt="PHP" /></a>
</p>

<p align="center">
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-installation">Installation</a> •
  <a href="#-features">Features</a> •
  <a href="#-api-reference">API</a> •
  <a href="#-documentation">Docs</a>
</p>

---

## 🚀 Quick Start

Get your API key at [renamed.to/settings](https://www.renamed.to/settings).

<table>
<tr>
<td><strong>TypeScript</strong></td>
<td><strong>Python</strong></td>
</tr>
<tr>
<td>

```typescript
import { RenamedClient } from '@renamed/sdk';

const client = new RenamedClient({
  apiKey: 'rt_...'
});

const result = await client.rename('invoice.pdf');
console.log(result.suggestedFilename);
// → "2025-01-15_AcmeCorp_INV-12345.pdf"
```

</td>
<td>

```python
from renamed import RenamedClient

client = RenamedClient(api_key='rt_...')

result = client.rename('invoice.pdf')
print(result.suggested_filename)
# → "2025-01-15_AcmeCorp_INV-12345.pdf"
```

</td>
</tr>
</table>

<details>
<summary><strong>More languages</strong></summary>

### Go

```go
import "github.com/renamed-to/renamed-sdk/sdks/go/renamed"

client := renamed.NewClient("rt_...")

result, _ := client.Rename(ctx, "invoice.pdf", nil)
fmt.Println(result.SuggestedFilename)
// → "2025-01-15_AcmeCorp_INV-12345.pdf"
```

### Java

```java
import to.renamed.sdk.*;

RenamedClient client = new RenamedClient("rt_...");

RenameResult result = client.rename(Path.of("invoice.pdf"), null);
System.out.println(result.getSuggestedFilename());
// → "2025-01-15_AcmeCorp_INV-12345.pdf"
```

### C#

```csharp
using Renamed.Sdk;

using var client = new RenamedClient("rt_...");

var result = await client.RenameAsync("invoice.pdf");
Console.WriteLine(result.SuggestedFilename);
// → "2025-01-15_AcmeCorp_INV-12345.pdf"
```

### Ruby

```ruby
require 'renamed'

client = Renamed::Client.new(api_key: 'rt_...')

result = client.rename('invoice.pdf')
puts result.suggested_filename
# → "2025-01-15_AcmeCorp_INV-12345.pdf"
```

### Rust

```rust
use renamed::RenamedClient;

let client = RenamedClient::new("rt_...");

let result = client.rename("invoice.pdf", None).await?;
println!("{}", result.suggested_filename);
// → "2025-01-15_AcmeCorp_INV-12345.pdf"
```

### Swift

```swift
import Renamed

let client = try RenamedClient(apiKey: "rt_...")

let file = try FileInput(url: URL(fileURLWithPath: "invoice.pdf"))
let result = try await client.rename(file: file)
print(result.suggestedFilename)
// → "2025-01-15_AcmeCorp_INV-12345.pdf"
```

### PHP

```php
use Renamed\Client;

$client = new Client('rt_...');

$result = $client->rename('invoice.pdf');
echo $result->suggestedFilename;
// → "2025-01-15_AcmeCorp_INV-12345.pdf"
```

</details>

---

## 📦 Installation

<table>
<tr><th>Language</th><th>Package Manager</th></tr>
<tr>
<td><strong>TypeScript</strong></td>
<td>

```bash
npm install @renamed/sdk
```

</td>
</tr>
<tr>
<td><strong>Python</strong></td>
<td>

```bash
pip install renamed
```

</td>
</tr>
<tr>
<td><strong>Go</strong></td>
<td>

```bash
go get github.com/renamed-to/renamed-sdk/sdks/go
```

</td>
</tr>
<tr>
<td><strong>Java</strong></td>
<td>

```xml
<dependency>
    <groupId>to.renamed</groupId>
    <artifactId>renamed-sdk</artifactId>
    <version>0.1.0</version>
</dependency>
```

</td>
</tr>
<tr>
<td><strong>C# / .NET</strong></td>
<td>

```bash
dotnet add package Renamed.Sdk
```

</td>
</tr>
<tr>
<td><strong>Ruby</strong></td>
<td>

```bash
gem install renamed
```

</td>
</tr>
<tr>
<td><strong>Rust</strong></td>
<td>

```toml
[dependencies]
renamed = "0.1"
```

</td>
</tr>
<tr>
<td><strong>Swift</strong></td>
<td>

```swift
.package(url: "https://github.com/renamed-to/renamed-sdk", from: "0.1.0")
```

</td>
</tr>
<tr>
<td><strong>PHP</strong></td>
<td>

```bash
composer require renamed/sdk
```

</td>
</tr>
</table>

---

## ✨ Features

### 🤖 Rename Files

AI-powered file renaming with intelligent naming suggestions:

```typescript
const result = await client.rename('scan001.pdf');
// {
//   suggestedFilename: "2025-01-15_AcmeCorp_INV-12345.pdf",
//   folderPath: "2025/AcmeCorp/Invoices",
//   confidence: 0.95
// }
```

### ✂️ Split PDFs

Split multi-page PDFs into individual documents:

```typescript
const job = await client.pdfSplit('multi-page.pdf', { mode: 'auto' });
const result = await job.wait();

for (const doc of result.documents) {
  const buffer = await client.downloadFile(doc.downloadUrl);
  // Save doc.filename with buffer
}
```

### 📊 Extract Data

Extract structured data from documents:

```typescript
const result = await client.extract('invoice.pdf', {
  prompt: 'Extract invoice number, date, and total amount'
});
console.log(result.data);
// { invoiceNumber: "INV-12345", date: "2025-01-15", total: 1234.56 }
```

---

## 📖 API Reference

| Method | Description |
|--------|-------------|
| `rename(file)` | Rename a file using AI |
| `pdfSplit(file, options)` | Split PDF into documents |
| `extract(file, options)` | Extract structured data |
| `getUser()` | Get user profile & credits |
| `downloadFile(url)` | Download a split document |

---

## 📋 Supported Files

| Type | Formats |
|------|---------|
| 📄 Documents | PDF |
| 🖼️ Images | JPEG, PNG, TIFF |

---

## 📚 Documentation

<p align="center">
  <a href="https://www.renamed.to/docs/api-docs"><img src="https://img.shields.io/badge/API_Docs-blue?style=for-the-badge" alt="API Docs" /></a>
</p>

<p align="center">
  <a href="./sdks/typescript/README.md"><img src="https://img.shields.io/badge/TypeScript-3178c6?style=flat-square&logo=typescript&logoColor=white" alt="TypeScript" /></a>
  <a href="./sdks/python/README.md"><img src="https://img.shields.io/badge/Python-3776ab?style=flat-square&logo=python&logoColor=white" alt="Python" /></a>
  <a href="./sdks/go/README.md"><img src="https://img.shields.io/badge/Go-00add8?style=flat-square&logo=go&logoColor=white" alt="Go" /></a>
  <a href="./sdks/java/README.md"><img src="https://img.shields.io/badge/Java-ed8b00?style=flat-square&logo=openjdk&logoColor=white" alt="Java" /></a>
  <a href="./sdks/csharp/README.md"><img src="https://img.shields.io/badge/C%23-512bd4?style=flat-square&logo=csharp&logoColor=white" alt="C#" /></a>
  <a href="./sdks/ruby/README.md"><img src="https://img.shields.io/badge/Ruby-cc342d?style=flat-square&logo=ruby&logoColor=white" alt="Ruby" /></a>
  <a href="./sdks/rust/README.md"><img src="https://img.shields.io/badge/Rust-dea584?style=flat-square&logo=rust&logoColor=black" alt="Rust" /></a>
  <a href="./sdks/swift/README.md"><img src="https://img.shields.io/badge/Swift-f05138?style=flat-square&logo=swift&logoColor=white" alt="Swift" /></a>
  <a href="./sdks/php/README.md"><img src="https://img.shields.io/badge/PHP-777bb4?style=flat-square&logo=php&logoColor=white" alt="PHP" /></a>
</p>

---

<p align="center">
  <a href="https://www.renamed.to">Website</a> •
  <a href="https://www.renamed.to/docs/api-docs">API Docs</a> •
  <a href="https://github.com/renamed-to/renamed-sdk/issues">Issues</a>
</p>

<p align="center">
  <sub>Built with ❤️ by the <a href="https://www.renamed.to">renamed.to</a> team</sub>
</p>
