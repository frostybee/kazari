// Site wide dark mode. Kazari does not ship a page level switch, only the
// per block button in each toolbar, so a real site provides its own. This
// one flips a .dark class on <html>, which is the selector the Kazari
// config declares under darkMode.
(function () {
  var STORAGE_KEY = "kazari-example-theme";
  var root = document.documentElement;
  var button = document.getElementById("theme-toggle");
  if (!button) return;

  var label = button.querySelector(".theme-toggle-label");

  function sync() {
    var dark = root.classList.contains("dark");
    button.setAttribute("aria-pressed", dark ? "true" : "false");
    if (label) label.textContent = dark ? "Light" : "Dark";
  }

  button.addEventListener("click", function () {
    var dark = root.classList.toggle("dark");
    try {
      localStorage.setItem(STORAGE_KEY, dark ? "dark" : "light");
    } catch (e) {}
    sync();
  });

  sync();
})();
