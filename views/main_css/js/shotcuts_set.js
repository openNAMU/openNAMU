all_list = [
  ["F", "/"],
  ["C", "/recent_changes"],
  ["D", "/recent_discuss"],
  ["A", "/random"]
];

all_list.forEach(function(element) {
  shortcut.add(
    element[0],
    function() {
      window.location.href = element[1];
    },
    {
      disable_in_input: true
    }
  );
});

all_list_2 = [["W", "/w"], ["H", "/history"], ["E", "/edit"]];

all_list_2.forEach(function(element) {
  shortcut.add(
    element[0],
    function() {
      href_d = window.location.href.split("/");
      if (href_d[4]) {
        window.location.href = element[1] + "/" + href_d[4];
      }
    },
    {
      disable_in_input: true
    }
  );
});
