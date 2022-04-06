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