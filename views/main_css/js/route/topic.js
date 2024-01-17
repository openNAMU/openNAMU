function opennamu_do_remove_blind_thread() {
    var style = document.querySelector('#opennamu_remove_blind');
    if (style !== null) {
        if (style.innerHTML !== "") {
            style.innerHTML = '';
        }
        else {
            style.innerHTML = "\n                .opennamu_comment_blind_js {\n                    display: none;\n                }\n            ";
        }
    }
}
