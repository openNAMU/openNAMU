"use strict";

let shortcut_key_list = {};

window.addEventListener("keyup", e => {
    delete shortcut_key_list[e.code];
});
window.addEventListener("blur", () => {
    shortcut_key_list = {};
});
document.addEventListener("visibilitychange", () => {
    if(document.hidden) {
        shortcut_key_list = {};
    }
});

window.addEventListener("keydown", e => {
    let shortcut_check = e.target.tagName.toLowerCase();
    if(shortcut_check === 'input' || shortcut_check === 'textarea') {
        return;
    } else if(e.repeat) {
        return;
    }

    shortcut_key_list[e.code] = 1;
    if(Object.keys(shortcut_key_list).length === 1) {
        let doc_shortcut = /^\/(w|w_from|history|edit|acl|topic|xref)\//i;

        if(shortcut_key_list['KeyF'] === 1) {
            window.location.href = '/';
        } else if(shortcut_key_list['KeyC'] === 1) {
            window.location.href = '/recent_change';
        } else if(shortcut_key_list['KeyD'] === 1) {
            window.location.href = '/recent_discuss';
        } else if(shortcut_key_list['KeyA'] === 1) {
            window.location.href = '/random';
        }

        if(window.location.pathname.match(doc_shortcut)) {
            let doc_href = window.location.pathname.replace(doc_shortcut, '');
            
            if(shortcut_key_list['KeyW'] === 1) {
                window.location.pathname = '/w/' + doc_href;
            } else if(shortcut_key_list['KeyE'] === 1) {
                window.location.pathname = '/edit/' + doc_href;
            } else if(shortcut_key_list['KeyH'] === 1) {
                window.location.pathname = '/history/' + doc_href;
            }
        }
    }
});