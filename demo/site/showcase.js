document.addEventListener('click', function(e) {
  var btn = e.target.closest('.kazari-code .kz-copy-btn');
  if (!btn) return;
  var code = btn.getAttribute('data-code');
  if (!code) return;
  code = code.replace(/\x7f/g, '\n');
  var iconEl = btn.querySelector('svg');
  function onSuccess() {
    btn.classList.add('copied');
    if (iconEl) {
      iconEl.innerHTML = '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>';
    }
    setTimeout(function() {
      btn.classList.remove('copied');
      if (iconEl) {
        iconEl.innerHTML = '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>';
      }
    }, 2000);
  }
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(code).then(onSuccess).catch(function() {
      fallbackCopy(code);
      onSuccess();
    });
  } else {
    fallbackCopy(code);
    onSuccess();
  }
});
function fallbackCopy(text) {
  var ta = document.createElement('textarea');
  ta.value = text;
  ta.style.position = 'fixed';
  ta.style.opacity = '0';
  document.body.appendChild(ta);
  ta.select();
  document.execCommand('copy');
  ta.remove();
}
(function() {
  var SCALE_KEY = 'kz-fs-font-scale';
  var STEP = 0.1;
  var MIN = 0.6;
  var MAX = 5.0;
  var DEFAULT = 1;

  function getScale(block) {
    if (block) {
      var current = parseFloat(block.style.getPropertyValue('--kz-fs-font-scale'));
      if (current) return current;
    }
    try { return parseFloat(localStorage.getItem(SCALE_KEY)) || DEFAULT; }
    catch(e) { return DEFAULT; }
  }

  function setScale(block, scale) {
    scale = Math.max(MIN, Math.min(MAX, Math.round(scale * 100) / 100));
    block.style.setProperty('--kz-fs-font-scale', scale);
    try { localStorage.setItem(SCALE_KEY, scale); } catch(e) {}
  }

  document.addEventListener('click', function(e) {
    var btn = e.target.closest('.kazari-code .kz-fs-btn');
    if (!btn) return;
    var block = btn.closest('.kazari-code');
    if (!block) return;
    if (document.fullscreenElement === block) {
      document.exitFullscreen();
    } else if (block.requestFullscreen) {
      block.requestFullscreen();
    } else if (block.webkitRequestFullscreen) {
      block.webkitRequestFullscreen();
    }
  });

  document.addEventListener('click', function(e) {
    var btn = e.target.closest('.kazari-code .kz-font-inc');
    if (!btn) return;
    var block = btn.closest('.kazari-code');
    if (block) setScale(block, getScale(block) + STEP);
  });

  document.addEventListener('click', function(e) {
    var btn = e.target.closest('.kazari-code .kz-font-dec');
    if (!btn) return;
    // The second click of a double click is handled by the dblclick reset.
    if (e.detail > 1) return;
    var block = btn.closest('.kazari-code');
    if (block) setScale(block, getScale(block) - STEP);
  });

  document.addEventListener('dblclick', function(e) {
    var btn = e.target.closest('.kazari-code .kz-font-dec');
    if (!btn) return;
    var block = btn.closest('.kazari-code');
    if (block) setScale(block, DEFAULT);
  });

  document.addEventListener('fullscreenchange', function() {
    var el = document.fullscreenElement;
    if (el && el.classList.contains('kazari-code')) {
      setScale(el, getScale(el));
      var btn = el.querySelector('.kz-fs-btn');
      if (btn) btn.setAttribute('aria-expanded', 'true');
    } else {
      var btns = document.querySelectorAll('.kz-fs-btn[aria-expanded="true"]');
      for (var i = 0; i < btns.length; i++) {
        btns[i].setAttribute('aria-expanded', 'false');
      }
    }
  });
})();
document.addEventListener('click', function(e) {
  var btn = e.target.closest('.kazari-code .kz-wrap-btn');
  if (!btn) return;
  var block = btn.closest('.kazari-code');
  if (!block) return;
  var pre = block.querySelector('pre');
  if (!pre) return;
  var wrapped = pre.classList.toggle('wrap');
  btn.setAttribute('aria-pressed', wrapped);
  var label = wrapped ? btn.getAttribute('data-disable') : btn.getAttribute('data-enable');
  btn.setAttribute('title', label);
});
document.addEventListener('click', function(e) {
    var btn = e.target.closest('.kazari-code .kz-collapse-btn, .kazari-code .kz-collapse-toggle');
    if (!btn) return;
    var block = btn.closest('.kazari-code');
    if (!block) return;
    var bottomBtn = block.querySelector('.kz-collapse-btn');
    var toggleBtn = block.querySelector('.kz-collapse-toggle');
    var announce = block.querySelector('.kz-sr-announce');
    var isCollapsed = block.classList.contains('kz-collapsed');
    if (isCollapsed) {
        block.classList.remove('kz-collapsed');
        if (bottomBtn) {
            bottomBtn.setAttribute('aria-expanded', 'true');
            bottomBtn.textContent = bottomBtn.getAttribute('data-collapse');
        }
        if (toggleBtn) toggleBtn.setAttribute('aria-expanded', 'true');
        if (announce && bottomBtn) announce.textContent = bottomBtn.getAttribute('data-expanded-msg') || '';
    } else {
        block.classList.add('kz-collapsed');
        if (bottomBtn) {
            bottomBtn.setAttribute('aria-expanded', 'false');
            bottomBtn.textContent = bottomBtn.getAttribute('data-expand');
        }
        if (toggleBtn) toggleBtn.setAttribute('aria-expanded', 'false');
        if (announce && bottomBtn) announce.textContent = bottomBtn.getAttribute('data-collapsed-msg') || '';
    }
});
(function () {
  function init() {
    var groups = document.querySelectorAll(".kz-group");
    var syncing = false;

    function findTabByLabel(tabs, label) {
      for (var i = 0; i < tabs.length; i++) {
        if (tabs[i].textContent === label) return i;
      }
      return -1;
    }

    function syncOtherGroups(syncKey, label, sourceGroup) {
      groups.forEach(function (otherGroup) {
        if (otherGroup === sourceGroup) return;
        if (otherGroup.getAttribute("data-sync") !== syncKey) return;
        var otherTabs = otherGroup.querySelectorAll('.kz-group-tabs button[role="tab"]');
        var otherPanels = otherGroup.querySelectorAll('.kz-group-panels > [role="tabpanel"]');
        var matchIndex = findTabByLabel(otherTabs, label);
        if (matchIndex < 0) return;
        otherTabs.forEach(function (t, i) {
          var selected = i === matchIndex;
          t.setAttribute("aria-selected", selected ? "true" : "false");
          t.setAttribute("tabindex", selected ? "0" : "-1");
          if (otherPanels[i]) {
            if (selected) {
              otherPanels[i].removeAttribute("hidden");
            } else {
              otherPanels[i].setAttribute("hidden", "");
            }
          }
        });
      });
    }

    groups.forEach(function (group, gi) {
      var tabs = group.querySelectorAll('.kz-group-tabs button[role="tab"]');
      var panels = group.querySelectorAll('.kz-group-panels > [role="tabpanel"]');
      if (tabs.length === 0) return;

      // Assign IDs for aria-labelledby linking.
      tabs.forEach(function (tab, i) {
        var tabId = "kz-g" + gi + "-tab-" + i;
        var panelId = "kz-g" + gi + "-panel-" + i;
        tab.id = tabId;
        if (panels[i]) {
          panels[i].id = panelId;
          panels[i].setAttribute("aria-labelledby", tabId);
          tab.setAttribute("aria-controls", panelId);
        }
      });

      function activate(index, shouldSync) {
        tabs.forEach(function (t, i) {
          var selected = i === index;
          t.setAttribute("aria-selected", selected ? "true" : "false");
          t.setAttribute("tabindex", selected ? "0" : "-1");
          if (panels[i]) {
            if (selected) {
              panels[i].removeAttribute("hidden");
            } else {
              panels[i].setAttribute("hidden", "");
            }
          }
        });
        if (shouldSync) tabs[index].focus();

        var syncKey = group.getAttribute("data-sync");
        if (syncKey && shouldSync && !syncing) {
          var label = tabs[index].textContent;
          syncing = true;
          try {
            localStorage.setItem("kz-tabs-" + syncKey, label);
          } catch (e) {}
          syncOtherGroups(syncKey, label, group);
          syncing = false;
        }
      }

      // Click handler.
      tabs.forEach(function (tab, i) {
        tab.addEventListener("click", function () {
          activate(i, true);
        });
      });

      // Keyboard navigation.
      group.querySelector(".kz-group-tabs").addEventListener("keydown", function (e) {
        var current = Array.prototype.indexOf.call(tabs, document.activeElement);
        if (current < 0) return;
        var next = -1;
        if (e.key === "ArrowRight") {
          next = (current + 1) % tabs.length;
        } else if (e.key === "ArrowLeft") {
          next = (current - 1 + tabs.length) % tabs.length;
        } else if (e.key === "Home") {
          next = 0;
        } else if (e.key === "End") {
          next = tabs.length - 1;
        }
        if (next >= 0) {
          e.preventDefault();
          activate(next, true);
        }
      });

      // Restore synced tab from localStorage on load.
      var syncKey = group.getAttribute("data-sync");
      if (syncKey) {
        try {
          var saved = localStorage.getItem("kz-tabs-" + syncKey);
          if (saved) {
            var matchIndex = findTabByLabel(tabs, saved);
            if (matchIndex >= 0) {
              activate(matchIndex, false);
            }
          }
        } catch (e) {}
      }
    });
  }

  // The DOM may already be ready when this script runs, for example when it
  // is injected dynamically or loaded async after the page has parsed.
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();

(function () {
  var root = document.documentElement;
  var announcement = document.getElementById('showcase-announcement');

  function announce(message) {
    if (!announcement) return;
    announcement.textContent = '';
    window.setTimeout(function () { announcement.textContent = message; }, 10);
  }

  function fallbackCopy(text) {
    var textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    textarea.remove();
  }

  function setCopiedState(button, message) {
    if (!button) return;
    var label = button.querySelector('span');
    var icon = button.querySelector('svg');
    if (button._showcaseCopyTimer) window.clearTimeout(button._showcaseCopyTimer);
    if (button._showcaseOriginalLabel === undefined) {
      button._showcaseOriginalLabel = label ? label.textContent : '';
    }
    if (button._showcaseOriginalIcon === undefined) {
      button._showcaseOriginalIcon = icon ? icon.innerHTML : '';
    }
    button.classList.add('copied');
    if (label) label.textContent = message;
    if (icon) icon.innerHTML = '<path d="m5 12 4 4L19 6"/>';
    button._showcaseCopyTimer = window.setTimeout(function () {
      button.classList.remove('copied');
      if (label) label.textContent = button._showcaseOriginalLabel;
      if (icon) icon.innerHTML = button._showcaseOriginalIcon;
      button._showcaseCopyTimer = null;
    }, 2000);
  }

  function copyText(text, message, button) {
    function onSuccess() {
      announce(message);
      setCopiedState(button, message);
    }
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(text).then(onSuccess).catch(function () {
        fallbackCopy(text);
        onSuccess();
      });
      return;
    }
    fallbackCopy(text);
    onSuccess();
  }

  function applyTheme(isDark) {
    root.classList.toggle('dark', isDark);
    document.body.classList.toggle('dark', isDark);
    document.querySelectorAll('.kazari-tinted, .kazari-scoped, .kazari-customizer').forEach(function (element) {
      element.classList.toggle('dark', isDark);
    });
    localStorage.setItem('kz-demo-theme', isDark ? 'dark' : 'light');
  }

  var saved = localStorage.getItem('kz-demo-theme');
  if (saved === 'dark') applyTheme(true);

  var themeToggle = document.getElementById('theme-toggle');
  if (themeToggle) {
    if (saved === 'dark') {
      themeToggle.setAttribute('aria-pressed', 'true');
      themeToggle.textContent = 'Light mode';
    }
    themeToggle.addEventListener('click', function () {
      var isDark = root.classList.toggle('dark');
      applyTheme(isDark);
      themeToggle.setAttribute('aria-pressed', String(isDark));
      themeToggle.textContent = isDark ? 'Light mode' : 'Dark mode';
      var navToggle = document.getElementById('dark-toggle');
      if (navToggle) navToggle.checked = isDark;
    });
  }

  var darkToggle = document.getElementById('dark-toggle');
  if (darkToggle) {
    if (saved === 'dark') darkToggle.checked = true;
    darkToggle.addEventListener('change', function () {
      var isDark = darkToggle.checked;
      applyTheme(isDark);
      if (themeToggle) {
        themeToggle.setAttribute('aria-pressed', String(isDark));
        themeToggle.textContent = isDark ? 'Light mode' : 'Dark mode';
      }
    });
  }

  document.querySelectorAll('[data-backend-link]').forEach(function (link) {
    link.addEventListener('click', function (event) {
      event.preventDefault();
      window.location.href = link.getAttribute('href').split('#')[0] + window.location.hash;
    });
  });

  var exampleJump = document.getElementById('example-jump');
  if (exampleJump) {
    exampleJump.addEventListener('change', function () {
      if (exampleJump.value) window.location.hash = exampleJump.value;
    });
  }

  document.querySelectorAll('[data-copy-link]').forEach(function (button) {
    button.addEventListener('click', function () {
      var id = button.getAttribute('data-copy-link');
      var url = window.location.href.split('#')[0] + '#' + id;
      copyText(url, 'Link copied', button);
    });
  });

  document.querySelectorAll('.copy-recipe').forEach(function (button) {
    button.addEventListener('click', function () {
      var code = button.parentElement.querySelector('code');
      var label = button.getAttribute('data-copy-label') || 'Recipe';
      if (code) copyText(code.textContent, label + ' copied', button);
    });
  });

  document.querySelectorAll('[data-recipe-group]').forEach(function (group, groupIndex) {
    var tabs = Array.prototype.slice.call(group.querySelectorAll('[data-recipe-tab]'));
    var panels = Array.prototype.slice.call(group.querySelectorAll('[data-recipe-panel]'));
    if (tabs.length < 2) return;

    function activate(index, focus) {
      tabs.forEach(function (tab, tabIndex) {
        var selected = tabIndex === index;
        tab.setAttribute('aria-selected', String(selected));
        tab.setAttribute('tabindex', selected ? '0' : '-1');
        panels[tabIndex].hidden = !selected;
      });
      if (focus) tabs[index].focus();
    }

    tabs.forEach(function (tab, index) {
      var tabID = 'showcase-recipe-' + groupIndex + '-tab-' + index;
      var panelID = 'showcase-recipe-' + groupIndex + '-panel-' + index;
      tab.id = tabID;
      tab.setAttribute('aria-controls', panelID);
      panels[index].id = panelID;
      panels[index].setAttribute('role', 'tabpanel');
      panels[index].setAttribute('aria-labelledby', tabID);
      tab.addEventListener('click', function () { activate(index, false); });
      tab.addEventListener('keydown', function (event) {
        var next = -1;
        if (event.key === 'ArrowRight') next = (index + 1) % tabs.length;
        if (event.key === 'ArrowLeft') next = (index - 1 + tabs.length) % tabs.length;
        if (event.key === 'Home') next = 0;
        if (event.key === 'End') next = tabs.length - 1;
        if (next >= 0) {
          event.preventDefault();
          activate(next, true);
        }
      });
    });
  });

  var searchInput = document.getElementById('example-search');
  var categorySelect = document.getElementById('category-filter');
  var clearButton = document.getElementById('clear-filters');
  var resultCount = document.getElementById('result-count');
  var noResults = document.getElementById('no-results');
  var examples = Array.prototype.slice.call(document.querySelectorAll('[data-example]'));
  var categories = Array.prototype.slice.call(document.querySelectorAll('[data-category]'));

  function applyFilters() {
    var query = searchInput ? searchInput.value.trim().toLowerCase() : '';
    var selectedCategory = categorySelect ? categorySelect.value : 'all';
    var visibleCount = 0;

    examples.forEach(function (example) {
      var category = example.closest('[data-category]').getAttribute('data-category');
      var matchesText = !query || example.getAttribute('data-search').indexOf(query) >= 0;
      var matchesCategory = selectedCategory === 'all' || category === selectedCategory;
      var visible = matchesText && matchesCategory;
      example.hidden = !visible;
      if (visible) visibleCount++;
      var navLink = document.querySelector('[data-example-link="' + example.id + '"]');
      if (navLink) navLink.hidden = !visible;
      var jumpOption = document.querySelector('[data-jump-example="' + example.id + '"]');
      if (jumpOption) jumpOption.hidden = !visible;
    });

    categories.forEach(function (category) {
      var categoryID = category.getAttribute('data-category');
      var hasVisibleExamples = category.querySelector('[data-example]:not([hidden])') !== null;
      category.hidden = !hasVisibleExamples;
      var navCategory = document.querySelector('[data-nav-category="' + categoryID + '"]');
      if (navCategory) navCategory.hidden = !hasVisibleExamples;
      var jumpCategory = document.querySelector('[data-jump-category="' + categoryID + '"]');
      if (jumpCategory) jumpCategory.hidden = !hasVisibleExamples;
    });

    if (resultCount) resultCount.textContent = visibleCount + (visibleCount === 1 ? ' example' : ' examples');
    if (noResults) noResults.hidden = visibleCount !== 0;
    updateScrollSpy();
  }

  if (searchInput) searchInput.addEventListener('input', applyFilters);
  if (categorySelect) categorySelect.addEventListener('change', applyFilters);
  if (clearButton) {
    clearButton.addEventListener('click', function () {
      searchInput.value = '';
      categorySelect.value = 'all';
      applyFilters();
      searchInput.focus();
    });
  }

  var demoNav = document.querySelector('.demo-nav');
  var scrollSpyLinks = new Map();
  document.querySelectorAll('[data-category-link], [data-example-link]').forEach(function (link) {
    scrollSpyLinks.set(link.getAttribute('data-category-link') || link.getAttribute('data-example-link'), link);
  });
  var scrollSpyTargets = Array.prototype.slice.call(document.querySelectorAll('[data-category], [data-example]'));
  var activeScrollSpyID = '';
  var scrollSpyScheduled = false;

  function targetIsVisible(target) {
    if (target.hidden) return false;
    var category = target.matches('[data-category]') ? target : target.closest('[data-category]');
    return !category || !category.hidden;
  }

  function keepNavLinkVisible(link) {
    if (!demoNav || !link || demoNav.scrollHeight <= demoNav.clientHeight) return;
    var linkTop = link.offsetTop;
    var linkBottom = linkTop + link.offsetHeight;
    var visibleTop = demoNav.scrollTop;
    var visibleBottom = visibleTop + demoNav.clientHeight;
    if (linkTop < visibleTop) demoNav.scrollTop = Math.max(0, linkTop - 12);
    if (linkBottom > visibleBottom) demoNav.scrollTop = linkBottom - demoNav.clientHeight + 12;
  }

  function setActiveScrollSpyTarget(target) {
    var targetID = target ? target.id : '';
    if (targetID === activeScrollSpyID) return;
    activeScrollSpyID = targetID;
    scrollSpyLinks.forEach(function (link) { link.removeAttribute('aria-current'); });
    document.querySelectorAll('[data-nav-category]').forEach(function (category) {
      category.classList.remove('is-active');
    });
    if (!target) {
      if (exampleJump) exampleJump.value = '';
      return;
    }

    var activeLink = scrollSpyLinks.get(targetID);
    if (activeLink) {
      activeLink.setAttribute('aria-current', 'location');
      keepNavLinkVisible(activeLink);
    }
    var category = target.matches('[data-category]') ? target : target.closest('[data-category]');
    if (category) {
      var navCategory = document.querySelector('[data-nav-category="' + category.id + '"]');
      if (navCategory) navCategory.classList.add('is-active');
    }
    if (exampleJump) {
      exampleJump.value = target.matches('[data-example]') ? '#' + targetID : '';
    }
  }

  function updateScrollSpy() {
    scrollSpyScheduled = false;
    var visibleTargets = scrollSpyTargets.filter(targetIsVisible);
    if (!visibleTargets.length) {
      setActiveScrollSpyTarget(null);
      return;
    }

    var header = document.querySelector('.demo-header');
    var headerHeight = header ? header.getBoundingClientRect().height : 0;
    var activationLine = headerHeight + Math.min(120, window.innerHeight * .2);
    var activeTarget = visibleTargets[0];
    visibleTargets.forEach(function (target) {
      if (target.getBoundingClientRect().top <= activationLine) activeTarget = target;
    });
    if (window.scrollY + window.innerHeight >= document.documentElement.scrollHeight - 2) {
      activeTarget = visibleTargets[visibleTargets.length - 1];
    }
    setActiveScrollSpyTarget(activeTarget);
  }

  function scheduleScrollSpyUpdate() {
    if (scrollSpyScheduled) return;
    scrollSpyScheduled = true;
    window.requestAnimationFrame(updateScrollSpy);
  }

  window.addEventListener('scroll', scheduleScrollSpyUpdate, { passive: true });
  window.addEventListener('resize', scheduleScrollSpyUpdate);
  updateScrollSpy();
})();

(function() {
  var btn = document.querySelector('.back-to-top');
  if (!btn) return;
  var footer = document.querySelector('.site-footer');
  var threshold = 300;
  var visible = false;
  window.addEventListener('scroll', function() {
    var pastTop = window.scrollY > threshold;
    var atFooter = footer && (window.innerHeight + window.scrollY >= footer.offsetTop);
    var show = pastTop && !atFooter;
    if (show !== visible) {
      visible = show;
      btn.classList.toggle('is-visible', show);
    }
  }, { passive: true });
  btn.addEventListener('click', function() {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  });
})();
