function load_user_info(name) {
  var n_ver = document.getElementById("get_user_info");

  var url = "/api/user_info/" + encodeURI(name) + "?render=1";

  var xhr = new XMLHttpRequest();
  xhr.open("GET", url, true);
  xhr.send(null);

  xhr.onreadystatechange = function() {
    if (this.readyState === 4 && this.status === 200) {
      n_ver.innerHTML += JSON.parse(this.responseText)["data"];
    }
  };
}
