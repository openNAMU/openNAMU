"use strict";

function opennamu_xss_filter(str) {
    return str.replace(/[&<>"']/g, function(match) {
        switch(match) {
            case '&':
                return '&amp;';
            case '<':
                return '&lt;';
            case '>':
                return '&gt;';
            case "'":
                return '&#x27;';
            case '"':
                return '&quot;';
        }
    });
}

function opennamu_xss_filter_decode(str) {
    return str.replace(/&amp;|&lt;|&gt;|&#x27;|&quot;/g, function(match) {
        switch(match) {
            case '&amp;':
                return '&';
            case '&lt;':
                return '<';
            case '&gt;':
                return '>';
            case '&#x27;':
                return "'";
            case '&quot;':
                return '"';
        }
    });
}

function opennamu_do_id_check(data) {
    if(data.match(/\.|\:/)) {
        return 0;
    } else {
        return 1;
    }
}

function opennamu_do_ip_render() {
    for(let for_a = 0; for_a < document.getElementsByClassName('opennamu_render_ip').length; for_a++) {
        let ip = document.getElementsByClassName('opennamu_render_ip')[for_a].innerHTML.replace(/&amp;/g, '&');

        fetch('/api/ip/' + opennamu_do_url_encode(ip)).then(function(res) {
            return res.json();
        }).then(function(data) {
            if(document.getElementsByClassName('opennamu_render_ip')[for_a].id !== "opennamu_render_end") {
                document.getElementsByClassName('opennamu_render_ip')[for_a].innerHTML = data["data"];
                document.getElementsByClassName('opennamu_render_ip')[for_a].id = "opennamu_render_end";
            }
        });
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

function opennamu_send_render(data) {
    if(data == '&lt;br&gt;' || data == '' || data.match(/^ +$/)) {
        data = '<br>';
    } else {
        data = data.replace(/( |^)(https?:\/\/(?:[^ ]+))/g, function(m0, m1, m2) {
            let link_main = m2;
            link_main = link_main.replace('"', '&quot;');

            return m1 + '<a href="' + link_main + '">' + link_main + '</a>';
        });
        data = data.replace(/&lt;a(?:(?:(?!&gt;).)*)&gt;((?:(?!&lt;\/a&gt;).)+)&lt;\/a&gt;/g, function(m0, m1) {
            let data_unescape = opennamu_xss_filter_decode(m1)

            return '<a href="/w/' + opennamu_do_url_encode(data_unescape) + '">' + m1 + '</a>'
        })
    }

    return data;
}

function opennamu_insert_v(name, data) {
    document.getElementById(name).value = data;
}

function opennamu_do_trace_spread() {
    if(document.getElementsByClassName('opennamu_trace')) {
        document.getElementsByClassName('opennamu_trace')[0].innerHTML = '' +
            '<style>.opennamu_trace_button { display: none; } .opennamu_trace { white-space: pre-wrap; overflow-x: unset; text-overflow: unset; }</style>' +
        '' + document.getElementsByClassName('opennamu_trace')[0].innerHTML
    }
}

function opennamu_do_render(to_obj, data, name = '', do_type = '', option = '') {
    let url;
    if(do_type === '') {
        url = "/api/render";
    } else {
        url = "/api/render/" + do_type;
    }

    fetch(url, {
        method : 'POST',
        headers : { 'Content-Type': 'application/x-www-form-urlencoded' },
        body : new URLSearchParams({
            'name' : name,
            'data': data,
            'option' : option
        })
    }).then(function(res) {
        return res.json();
    }).then(function(text) {
        if(document.getElementById(to_obj)) {
            if(text["data"]) {
                document.getElementById(to_obj).innerHTML = text["data"];
                eval(text["js_data"]);
            } else {
                document.getElementById(to_obj).innerHTML = '';
            }
        }
    });
}

function opennamu_page_control(url, page, data_length, data_length_max = 50) {
    let next = function() {
        if(data_length_max === data_length) {
            return '<a href="' + url.replace('{}', String(page + 1)) + '">(+)</a>';
        } else {
            return '';
        }
    };

    let back = function() {
        if(page !== 1) {
            return '<a href="' + url.replace('{}', String(page - 1)) + '">(-)</a>';
        } else {
            return '';
        }
    };

    return (next() + ' ' + back()).replace(/^ /, '');
}