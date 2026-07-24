---
title: "Localization"
description: "Configure the locale and override UI strings for non-English deployments."
sidebar:
  order: 18
---

Kazari ships built-in UI strings for toolbar tooltips, collapsible button labels, and screen reader announcements. Set the locale to use a bundled translation, or override individual strings for full control.

## Set the locale

Three locales are bundled: `en-US` (default), `fr-FR`, and `ja-JP`.

```go
kazari.WithLocale("fr-FR")
```

Or in the config file:

```yaml
locale: "fr-FR"
```

An unrecognized locale falls back to `en-US` silently.

## Override individual strings

Override any UI string key with `WithUIStrings`. Keys not present in the map keep their locale default:

```go
kazari.WithUIStrings(map[string]string{
    "copy.label":   "Copier",
    "copy.success": "Copie !",
})
```

Or in the config file:

```yaml
uiStrings:
  copy.label: "Copier"
  copy.success: "Copie !"
```

Overrides compose with `WithLocale`. Start from a locale and tweak specific strings:

```go
kazari.WithLocale("fr-FR"),
kazari.WithUIStrings(map[string]string{
    "copy.label": "Copier le code",
}),
```

→ All French strings are used except `copy.label`, which gets the custom value.

Multiple `WithUIStrings` calls merge additively. Later calls win on key conflicts. Unknown keys are silently ignored.

## Where strings appear

Kazari buttons are icon-only. Localized text appears in three places:

- **Tooltips** (`data-tooltip`): shown on hover via a CSS `::after` pseudo-element. Hidden on touch devices and in print.
- **Aria labels** (`aria-label`): read by screen readers for accessibility.
- **Screen reader announcements** (`aria-live="polite"`): announced when state changes (copy success, collapse toggle, theme toggle).

Some strings also appear as visible text: the expand/collapse bar button label and the collapsed line count badge.

## UI string keys

All available keys and their `en-US` defaults:

| Key | Default | Description |
|---|---|---|
| `copy.label` | Copy | Copy button tooltip and aria-label |
| `copy.title` | Copy to clipboard | Reserved (not yet rendered) |
| `copy.success` | Copied! | Tooltip flash and screen reader announcement after copying |
| `fullscreen.label` | Fullscreen | Fullscreen button tooltip |
| `fullscreen.font.increase` | Increase font size | Font increase button tooltip (fullscreen mode) |
| `fullscreen.font.decrease` | Decrease font size | Font decrease button tooltip (fullscreen mode) |
| `fullscreen.font.reset` | Double-click to reset | Reserved (not yet rendered) |
| `fullscreen.hint` | Press Esc to exit fullscreen | Reserved (not yet rendered) |
| `wrap.enable` | Enable word wrap | Wrap button tooltip when wrap is off |
| `wrap.disable` | Disable word wrap | Wrap button tooltip when wrap is on |
| `collapse.expand` | Show more | Expand button tooltip and visible text |
| `collapse.collapse` | Show less | Collapse button tooltip and visible text |
| `collapse.expanded` | Code block expanded | Screen reader announcement on expand |
| `collapse.collapsed` | Code block collapsed | Screen reader announcement on collapse |
| `collapse.summary.singular` | 1 collapsed line | Visible badge text (singular) |
| `collapse.summary.plural` | %d collapsed lines | Visible badge text (plural, `%d` is line count) |
| `codegroup.fallback` | Code | Reserved (not yet rendered) |
| `terminal.label` | Terminal window | Screen reader label for untitled terminal frames |
| `theme.toggle` | Toggle theme | Theme toggle button tooltip |
| `theme.toggle.announcement` | Theme toggled | Screen reader announcement on theme switch |
| `output.label` | Output | Output panel header when no per-block `outputLabel` is set |

## Configuration

| Option | Default | YAML key | Description |
|---|---|---|---|
| `WithLocale(string)` | `"en-US"` | `locale` | Built-in locale for UI strings |
| `WithUIStrings(map[string]string)` | `nil` | `uiStrings` | Override individual UI string keys |

## Edge cases

- Unknown locales fall back to `en-US` silently. No error is raised.
- Locale keys are case-sensitive and must match exactly (`"fr-FR"`, not `"fr-fr"` or `"fr"`).
- Resolution happens once at engine construction, not per-render.
- `CollapsibleConfig.ExpandButtonText` and `CollapseButtonText` override the locale strings for collapse buttons when non-empty. These take priority over both the locale default and `WithUIStrings` overrides.
- The `codegroup.fallback` key is accepted but currently unused. The code group tab fallback label is hardcoded to "Code".
- Reserved keys (`copy.title`, `fullscreen.font.reset`, `fullscreen.hint`) are defined in the locale tables but not rendered anywhere. They exist for future use.
