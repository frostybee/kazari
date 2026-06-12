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
