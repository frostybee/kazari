document.addEventListener('click', function(e) {
    var btn = e.target.closest('.kazari-code .kz-collapse-btn');
    if (!btn) return;
    var block = btn.closest('.kazari-code');
    if (!block) return;
    var announce = block.querySelector('.kz-sr-announce');
    var expanded = btn.getAttribute('aria-expanded') === 'true';
    if (expanded) {
        block.classList.add('kz-collapsed');
        btn.setAttribute('aria-expanded', 'false');
        btn.textContent = btn.getAttribute('data-expand');
        if (announce) announce.textContent = btn.getAttribute('data-collapsed-msg') || '';
    } else {
        block.classList.remove('kz-collapsed');
        btn.setAttribute('aria-expanded', 'true');
        btn.textContent = btn.getAttribute('data-collapse');
        if (announce) announce.textContent = btn.getAttribute('data-expanded-msg') || '';
    }
});
