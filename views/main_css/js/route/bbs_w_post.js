function opennamu_change_comment(get_id) {
    var _a;
    var input = document.querySelector('#opennamu_comment_select');
    if (input !== null) {
        input.value = get_id;
        (_a = document.getElementById('opennamu_edit_textarea')) === null || _a === void 0 ? void 0 : _a.focus();
    }
}
function opennamu_return_comment() {
    var _a;
    var input = document.querySelector('#opennamu_comment_select');
    if (input !== null) {
        (_a = document.getElementById(input.value)) === null || _a === void 0 ? void 0 : _a.focus();
    }
}
