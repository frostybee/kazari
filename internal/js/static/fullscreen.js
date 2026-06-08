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
