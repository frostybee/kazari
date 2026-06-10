document.addEventListener("DOMContentLoaded", function () {
  var groups = document.querySelectorAll(".kz-group");
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

    function activate(index) {
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
      tabs[index].focus();
    }

    // Click handler.
    tabs.forEach(function (tab, i) {
      tab.addEventListener("click", function () {
        activate(i);
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
        activate(next);
      }
    });
  });
});
