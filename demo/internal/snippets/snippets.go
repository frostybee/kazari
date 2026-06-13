package snippets

type Snippet struct {
	ID    string
	Label string
	Lang  string
	Code  string
}

var All = []Snippet{
	{ID: "go", Label: "Go", Lang: "go", Code: GoCode},
	{ID: "js", Label: "JavaScript", Lang: "javascript", Code: JSCode},
	{ID: "ts", Label: "TypeScript", Lang: "typescript", Code: TSCode},
	{ID: "py", Label: "Python", Lang: "python", Code: PyCode},
	{ID: "bash", Label: "Bash", Lang: "bash", Code: BashCode},
	{ID: "php", Label: "PHP", Lang: "php", Code: PHPCode},
	{ID: "css", Label: "CSS", Lang: "css", Code: CSSCode},
	{ID: "html", Label: "HTML", Lang: "html", Code: HTMLCode},
}

const GoCode = `package main

import "fmt"

func main() {
	name := "world"
	fmt.Printf("Hello, %s!\n", name)
	for i := 0; i < 3; i++ {
		fmt.Println(i)
	}
}`

const JSCode = `const greet = (name) => {
  console.log("Hello, " + name + "!");
  return { greeting: name, time: Date.now() };
};

const users = ["Alice", "Bob"];
users.forEach((u) => greet(u));`

const TSCode = `interface CacheEntry<T> {
  key: string;
  value: T;
  expiresAt: number;
}

function getOrSet<T>(cache: Map<string, CacheEntry<T>>, key: string, fn: () => T): T {
  const entry = cache.get(key);
  if (entry && entry.expiresAt > Date.now()) {
    return entry.value;
  }
  const value = fn();
  cache.set(key, { key, value, expiresAt: Date.now() + 3600_000 });
  return value;
}`

const PyCode = `from dataclasses import dataclass
from typing import Optional

@dataclass
class User:
    name: str
    email: str
    age: Optional[int] = None

    def greet(self) -> str:
        return f"Hello, {self.name}!"

users = [User("Alice", "alice@example.com", 30)]
for u in users:
    print(u.greet())`

const BashCode = `#!/bin/bash
set -euo pipefail

PROJECT_DIR="${1:-.}"
echo "Building project in $PROJECT_DIR..."

for file in "$PROJECT_DIR"/*.go; do
  if [[ -f "$file" ]]; then
    go build -o "bin/$(basename "${file%.go}")" "$file"
    echo "  Built: $file"
  fi
done

echo "Done. $(ls bin/ | wc -l) binaries built."`

const PHPCode = `<?php

namespace App\Http\Controllers;

class UserController extends Controller
{
    public function __construct(
        private UserRepository $users,
        private LoggerInterface $logger,
    ) {}

    public function show(int $id): Response
    {
        $user = $this->users->find($id);
        if ($user === null) {
            throw new NotFoundHttpException();
        }
        return $this->json($user);
    }
}`

const CSSCode = `:root {
  --primary: #5965d8;
  --bg: #f8f9fa;
  --text: #1a1a1a;
  --radius: 0.5rem;
}

.card {
  padding: 1.5rem;
  border-radius: var(--radius);
  background: var(--bg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.18);
}`

const HTMLCode = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>My Page</title>
  <link rel="stylesheet" href="/styles.css">
</head>
<body>
  <header>
    <nav>
      <a href="/">Home</a>
      <a href="/about">About</a>
    </nav>
  </header>
  <main>
    <h1>Welcome</h1>
    <p>This is a sample page.</p>
  </main>
  <script src="/app.js" defer></script>
</body>
</html>`
