function ringo_do_regex_data(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function ringo_get_post() {
    const check = document.getElementById('invert');
    if(check.checked === true) {
        document.cookie = 'main_css_darkmode=1; path=/';
    } else {
        document.cookie = 'main_css_darkmode=0; path=/';
    }

    const check_2 = document.getElementById('use_sys_darkmode');
    if(check_2.checked === true) {
        document.cookie = 'main_css_use_sys_darkmode=1; path=/';
    } else {
        document.cookie = 'main_css_use_sys_darkmode=0; path=/';
    }

    history.go(0);
}

function ringo_do_skin_set() {
    let cookies = document.cookie;
    if(!cookies.match(ringo_do_regex_data('main_css_use_sys_darkmode')) || (cookies.match(ringo_do_regex_data('main_css_use_sys_darkmode')) && cookies.match(ringo_do_regex_data('main_css_use_sys_darkmode'))[1] === '1')) {
        if(window.matchMedia('(prefers-color-scheme: dark)').matches) {
            document.cookie = 'main_css_darkmode=1; path=/';
        } else {
            document.cookie = 'main_css_darkmode=0; path=/';
        }
    }
}

function ringo_load_skin_set() {
    let cookies = document.cookie;
    
    if(window.location.pathname === '/change/skin_set') {
        let set_language = {
            "en-US" : {
                "save" : "Save",
                "darkmode" : "Darkmode",
                "use_sys_darkmode" : "Use system darkmode set",
            }, "ko-KR" : {
                "save" : "저장",
                "darkmode" : "다크모드",
                "use_sys_darkmode" : "시스템 다크모드 설정 사용",
            }
        }

        let language = cookies.match(ringo_do_regex_data('language'))[1];
        let user_language = cookies.match(ringo_do_regex_data('user_language'))[1];
        if(user_language in set_language) {
            language = user_language;
        }

        if(!language in set_language) {
            language = "en-US";
        }

        let set_data = {};

        if(cookies.match(ringo_do_regex_data('main_css_darkmode')) && cookies.match(ringo_do_regex_data('main_css_darkmode'))[1] === '1') {
            set_data["invert"] = "checked";
        }

        if(!cookies.match(ringo_do_regex_data('main_css_use_sys_darkmode')) || (cookies.match(ringo_do_regex_data('main_css_use_sys_darkmode')) && cookies.match(ringo_do_regex_data('main_css_use_sys_darkmode'))[1] === '1')) {
            set_data["use_sys_darkmode"] = "checked";
        }

        document.getElementById("main_skin_set").innerHTML = ' \
            <input ' + set_data["use_sys_darkmode"] + ' type="checkbox" id="use_sys_darkmode" name="use_sys_darkmode" value="use_sys_darkmode"> ' + set_language[language]['use_sys_darkmode'] + ' \
            <hr class="main_hr"> \
            <input ' + set_data["invert"] + ' type="checkbox" id="invert" name="invert" value="invert"> ' + set_language[language]['darkmode'] + ' \
            <hr class="main_hr"> \
            <button onclick="ringo_get_post();">' + set_language[language]['save'] + '</button> \
        ';
    }
}

window.addEventListener('DOMContentLoaded', ringo_do_skin_set);
window.addEventListener('DOMContentLoaded', ringo_load_skin_set);