"use strict";

function ringo_do_regex_data(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function ringo_do_skin_set() {
    let cookies = document.cookie;

    if(!window.localStorage.getItem('main_css_use_sys_darkmode') || window.localStorage.getItem('main_css_use_sys_darkmode') === '1') {
        if(cookies.match(ringo_do_regex_data('main_css_darkmode')) && cookies.match(ringo_do_regex_data('main_css_darkmode'))[1] === '1') {
            if(!window.matchMedia('(prefers-color-scheme: dark)').matches) {
                document.cookie = 'main_css_darkmode=0; path=/';
                history.go(0);
            }
        }
        
        if(!cookies.match(ringo_do_regex_data('main_css_darkmode')) || cookies.match(ringo_do_regex_data('main_css_darkmode'))[1] === '0') {
            if(window.matchMedia('(prefers-color-scheme: dark)').matches) {
                document.cookie = 'main_css_darkmode=1; path=/';
                history.go(0);
            }
        }
    }

    if(window.localStorage.getItem('main_css_off_sidebar') && window.localStorage.getItem('main_css_off_sidebar') === '0') {
    } else {
        document.getElementById('ringo_add_style').innerHTML += `
            section {
                width: auto !important;
                display: block !important;
                margin: auto !important;
            }

            .do_fixed {
                display: none !important;
            }
        `;
    }

    if(window.localStorage.getItem('main_css_fixed_width') && window.localStorage.getItem('main_css_fixed_width') !== '') {
        let fixed_width_data = window.localStorage.getItem('main_css_fixed_width');
        document.getElementById('ringo_add_style').innerHTML += `
            article.main {
                max-width: ` + fixed_width_data + `px !important;
            }
        `;
    }

    if(window.localStorage.getItem('main_css_sidebar_right') && window.localStorage.getItem('main_css_sidebar_right') === '1') {
        document.getElementById('ringo_add_style').innerHTML += `
            .do_fixed {
                float: right !important;
            }
        `;
    }
}

ringo_do_skin_set();