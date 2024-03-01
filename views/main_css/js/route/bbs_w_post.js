"use strict";

function opennamu_change_comment(get_id) {
    const input = document.querySelector('#opennamu_comment_select');
    if(input !== null) {
        input.value = get_id;
        document.getElementById('opennamu_edit_textarea')?.focus();
    }
}

function opennamu_return_comment() {
    const input = document.querySelector('#opennamu_comment_select');
    if(input !== null) {
        document.getElementById(input.value)?.focus();
    }
}