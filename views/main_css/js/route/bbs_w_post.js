"use strict";

function opennamu_change_comment(get_id) {
    const input = document.querySelector('#opennamu_comment_select');
    if(input !== null) {
        input.value = get_id;
        document.getElementById('opennamu_comment_select')?.focus();
    }
}

function opennamu_return_comment() {
    const input = document.querySelector('#opennamu_comment_select');
    if(input !== null) {
        document.getElementById(input.value)?.focus();
    }
}

function opennamu_post_tabom(bbs_id, bbs_code) {
    fetch('/api/v2/bbs/w/tabom/' + bbs_id + '-' + bbs_code, {
        method : 'POST'
    }).then(function(res) {
        return res.json();
    }).then(function(data) {
        opennamu_load_tabom_count(bbs_id, bbs_code);
    });
}

function opennamu_load_tabom_count(bbs_id, bbs_code) {
    fetch('/api/v2/bbs/w/tabom/' + bbs_id + '-' + bbs_code).then(function(res) {
        return res.json();
    }).then(function(data) {
        if(data) {
            for(let for_a = 0; for_a < document.getElementsByClassName('opennamu_tabom_count').length; for_a++) {
                document.getElementsByClassName('opennamu_tabom_count')[for_a].innerHTML = data["data"];
            }
        }
    });
}