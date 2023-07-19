"use strict";

function opennamu_do_id_check(data) {
    if(data.match(/\.|\:/)) {
        return 0;
    } else {
        return 1;
    }
}

function opennamu_do_url_encode(data) {
    return encodeURIComponent(data);
}

function opennamu_cookie_split_regex(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function opennamu_get_main_skin_set(set_name) {
    return fetch("/api/setting/" + opennamu_do_url_encode(set_name)).then(function(res) {
        return res.json();
    }).then(function(text) {
        if(
            document.cookie.match(opennamu_cookie_split_regex(set_name)) &&
            document.cookie.match(opennamu_cookie_split_regex(set_name))[1] !== '' &&
            document.cookie.match(opennamu_cookie_split_regex(set_name))[1] !== 'default'
        ) {
            return document.cookie.match(opennamu_cookie_split_regex(set_name))[1];
        } else {
            if(text[set_name]) {
                return text[set_name][0][0];
            } else {
                return '';
            }
        }
    });
}