function ringo_do_regex_data(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function ringo_get_post() {
    check = document.getElementById('invert');
    if(check.checked === true) {
        document.cookie = 'main_css_darkmode=1; path=/';
    } else {
        document.cookie = 'main_css_darkmode=0; path=/';
    }

    history.go(0);
}

function ringo_do_skin_set() {
    let cookies = document.cookie;
    
    if(window.location.pathname === '/change/skin_set') {
        let set_language = {
            "en-US" : {
                "save" : "Save",
                "darkmode" : "Darkmode"
            }, "ko-KR" : {
                "save" : "저장",
                "darkmode" : "다크모드"
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

        if(
            cookies.match(ringo_do_regex_data('main_css_darkmode')) &&
            cookies.match(ringo_do_regex_data('main_css_darkmode'))[1] === '1'
        ) {
            set_data["invert"] = "checked";
        }

        document.getElementById("main_skin_set").innerHTML = ' \
            <input ' + set_data["invert"] + ' type="checkbox" id="invert" name="invert" value="invert"> ' + set_language[language]['darkmode'] + ' \
            <hr class="main_hr"> \
            <button onclick="ringo_get_post();">' + set_language[language]['save'] + '</button> \
        ';
    }
}

window.addEventListener('DOMContentLoaded', ringo_do_skin_set);