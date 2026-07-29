---
title: "Client-Side Extensibility"
description: "Use Kazari's DOM structure, data attributes, and event delegation to extend code blocks with CSS and JavaScript."
tags: [client-side-js, accessibility, css-variables]
sidebar:
  order: 3
---

Kazari outputs plain HTML with consistent CSS classes, data attributes, and ARIA attributes. Custom JavaScript and CSS can interact with any code block without a Go-side API. The DOM is the contract.

## Event delegation pattern

All of Kazari's built-in JS uses event delegation on `document`. Follow the same pattern for consistency and to handle dynamically inserted blocks:

```javascript
document.addEventListener('click', function(e) {
  var btn = e.target.closest('.kazari-block .my-custom-btn');
  if (!btn) return;
  var block = btn.closest('.kazari-block');
  if (!block) return;
  // interact with the block
});
```

Use `.closest()` to find the nearest matching ancestor. This works regardless of how deeply nested the click target is (e.g., clicking an SVG path inside a button).

## DOM structure

Every code block follows this hierarchy:

```
.kazari-block                          Root container
  └─ .frame[data-lang="go"]          Frame wrapper (editor or terminal)
      ├─ .kz-toolbar                 Toolbar (editor frames)
      │   ├─ .kz-toolbar-left        Language badge + title
      │   │   ├─ .kz-lang            Language label
      │   │   ├─ .kz-file-icon       File icon (if enabled)
      │   │   └─ .kz-title           File title
      │   └─ .kz-toolbar-right       Action buttons
      │       ├─ .kz-copy-btn        Copy button
      │       ├─ .kz-wrap-btn        Wrap toggle
      │       ├─ .kz-font-controls   Font +/- (visible in fullscreen only)
      │       ├─ .kz-fs-btn          Fullscreen toggle
      │       └─ .kz-collapse-toggle Collapse toggle (if threshold enabled)
      │
      ├─ .kz-fs-hint                 "Press Esc to exit" hint
      │
      └─ pre[data-language="go"]     Code container
          └─ code                    Token wrapper
              └─ .kz-line            One per line
                  ├─ .kz-gutter > .kz-ln   Line number
                  └─ .kz-code           Line content (token spans)
```

Terminal frames use `.kz-terminal-header` instead of `.kz-toolbar`, with `.kz-terminal-actions` for buttons.

## Finding elements

### Find the code block from any child element

```javascript
var block = element.closest('.kazari-block');
```

### Get the language

```javascript
var frame = block.querySelector('.frame');
var lang = frame.getAttribute('data-lang'); // "go", "javascript", etc.
```

### Get the raw source code

```javascript
var copyBtn = block.querySelector('.kz-copy-btn');
var code = copyBtn.getAttribute('data-code').replace(/\x7f/g, '\n');
```

The `data-code` attribute encodes newlines as `\x7f` (DEL character). Replace them to get the original source.

### Get the line count

```javascript
var count = parseInt(block.getAttribute('data-lines'), 10);
```

### Get all lines

```javascript
var lines = block.querySelectorAll('.kz-line');
lines.forEach(function(line) {
  var content = line.querySelector('.kz-code').textContent;
});
```

## Reading state

### Is the block collapsed?

```javascript
var isCollapsed = block.classList.contains('kz-collapsed');
```

### Is word wrap enabled?

```javascript
var pre = block.querySelector('pre');
var isWrapped = pre.classList.contains('wrap');
```

### Is the block in fullscreen?

```javascript
var isFullscreen = document.fullscreenElement === block;
```

### Is a line marked/highlighted?

```javascript
var line = block.querySelector('.kz-line:nth-child(3)');
line.classList.contains('mark');     // highlighted
line.classList.contains('ins');      // inserted (diff)
line.classList.contains('del');      // deleted (diff)
line.classList.contains('focused');  // focus line
```

## Modifying state

### Toggle word wrap programmatically

```javascript
var pre = block.querySelector('pre');
pre.classList.toggle('wrap');
```

### Copy code programmatically

```javascript
var copyBtn = block.querySelector('.kz-copy-btn');
var code = copyBtn.getAttribute('data-code').replace(/\x7f/g, '\n');
navigator.clipboard.writeText(code);
```

### Enter/exit fullscreen programmatically

```javascript
// Enter
block.requestFullscreen();

// Exit
if (document.fullscreenElement) {
  document.exitFullscreen();
}
```

### Expand/collapse programmatically

```javascript
block.classList.toggle('kz-collapsed');
```

### Set font scale in fullscreen

```javascript
block.style.setProperty('--kz-fs-font-scale', 2.0);
localStorage.setItem('kz-fs-font-scale', 2.0);
```

## Listening for state changes

### Collapse state

```javascript
var observer = new MutationObserver(function(mutations) {
  mutations.forEach(function(m) {
    if (m.attributeName === 'class') {
      var collapsed = m.target.classList.contains('kz-collapsed');
      console.log('Collapsed:', collapsed);
    }
  });
});
observer.observe(block, { attributes: true, attributeFilter: ['class'] });
```

### Fullscreen state

```javascript
document.addEventListener('fullscreenchange', function() {
  if (document.fullscreenElement && document.fullscreenElement.classList.contains('kazari-block')) {
    console.log('Entered fullscreen');
  } else {
    console.log('Exited fullscreen');
  }
});
```

### Tab changes in code groups

```javascript
document.addEventListener('click', function(e) {
  var tab = e.target.closest('.kz-group-tabs button[role="tab"]');
  if (!tab) return;
  console.log('Tab switched to:', tab.textContent);
});
```

## Adding a custom toolbar button

Inject a button into the toolbar after Kazari renders:

```javascript
document.querySelectorAll('.kazari-block .kz-toolbar-right').forEach(function(toolbar) {
  var btn = document.createElement('button');
  btn.className = 'kz-copy-btn'; // reuse existing button styling
  btn.setAttribute('title', 'Download code');
  btn.innerHTML = '<svg>...</svg>';
  toolbar.appendChild(btn);
});
```

Then handle clicks via event delegation:

```javascript
document.addEventListener('click', function(e) {
  var btn = e.target.closest('.kazari-block .my-download-btn');
  if (!btn) return;
  var block = btn.closest('.kazari-block');
  var copyBtn = block.querySelector('.kz-copy-btn');
  var code = copyBtn.getAttribute('data-code').replace(/\x7f/g, '\n');
  var blob = new Blob([code], { type: 'text/plain' });
  var a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = 'code.txt';
  a.click();
});
```

## Data attributes reference

| Element | Attribute | Value |
|---|---|---|
| `.kazari-block` | `data-lines` | Line count (e.g., `"15"`) |
| `.frame` | `data-lang` | Language identifier (e.g., `"javascript"`) |
| `pre` | `data-language` | Language identifier |
| `.kz-copy-btn` | `data-code` | Source code (newlines as `\x7f`) |
| `.kz-copy-btn` | `data-copied` | Success message text |
| `.kz-wrap-btn` | `data-enable` | "Enable word wrap" label |
| `.kz-wrap-btn` | `data-disable` | "Disable word wrap" label |
| `.kz-collapse-btn` | `data-expand` | Expand button text |
| `.kz-collapse-btn` | `data-collapse` | Collapse button text |
| `.kz-collapse-btn` | `data-expanded-msg` | Screen reader announcement |
| `.kz-collapse-btn` | `data-collapsed-msg` | Screen reader announcement |
| `.kz-file-icon` | `data-ext` | File extension (e.g., `"js"`) |
| `.kz-group` | `data-sync` | Tab sync key for cross-group sync |
| `.kz-line .kz-code` | `data-label` | Marker label text |

## CSS custom properties

Override any `--kz-*` variable in a stylesheet to change Kazari's appearance without touching Go code. The full list is in the [CSS variables reference](/reference/css-variables/). Key ones for JS interaction:

| Variable | Purpose | Set by |
|---|---|---|
| `--kz-fs-font-scale` | Font scale multiplier in fullscreen | JS (fullscreen.js) |
| `--kz-indent` | Hanging indent per line | Go renderer (inline style) |
| `--kz-editor-bg` | Code block background | Theme system |
| `--kz-editor-fg` | Code block foreground | Theme system |

## localStorage keys

| Key | Purpose | Format |
|---|---|---|
| `kz-fs-font-scale` | Persisted font scale for fullscreen | Float (e.g., `"1.3"`) |
| `kz-tabs-{syncKey}` | Selected tab index for synced code groups | Integer (e.g., `"0"`) |

## The pattern

Kazari exposes structured data as DOM attributes (`data-lang`, `data-lines`, `data-ext`, CSS classes, ARIA states). CSS selectors replace AST queries. Event delegation replaces lifecycle hooks. `MutationObserver` replaces state-change callbacks. CSS custom properties override every visual default.

Per-language styling is a CSS selector, not an API call:

```css
.kazari-block .frame[data-lang="python"] .kz-toolbar {
  border-bottom-color: #ffd43b;
}
```

## Related

- [Post-Render Callbacks](/plugins/post-render-callbacks/) for Go-side HTML modification after rendering
- [Custom Highlighters](/plugins/custom-highlighters/) for swapping tokenization engines or accessing raw tokens
- [CSS Custom Properties](/styling/css-custom-properties/) for the full styling contract
