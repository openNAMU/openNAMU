"use strict";

let shortcut_key_list = [];
document.onkeyup = function(e) {
    delete shortcut_key_list[e.key];
}

document.onkeypress = function(e) {
    let shortcut_check = event.target.tagName.toLowerCase();
    if(
        shortcut_check !== 'input' &&
        shortcut_check !== 'textarea'
    ) {
        let doc_shortcut = /^\/(w|w_rev\/[0-9]+|w_from|history|edit|acl|topic|xref)\//i;

        shortcut_key_list[e.key] = 1;
        if(shortcut_key_list['f'] === 1) {
            window.location.href = '/';
        } else if(shortcut_key_list['c'] === 1) {
            window.location.href = '/recent_change';
        } else if(shortcut_key_list['d'] === 1) {
            window.location.href = '/recent_discuss';
        } else if(shortcut_key_list['a'] === 1) {
            window.location.href = '/random';
        }

        if(window.location.pathname.match(doc_shortcut)) {
            let doc_href = window.location.pathname.replace(doc_shortcut, '');
            
            if(shortcut_key_list['w'] === 1) {
                window.location.pathname = '/w/' + doc_href;
            } else if(shortcut_key_list['e'] === 1) {
                window.location.pathname = '/edit/' + doc_href;
            } else if(shortcut_key_list['h'] === 1) {
                window.location.pathname = '/history/' + doc_href;
            }
        }
    }
}