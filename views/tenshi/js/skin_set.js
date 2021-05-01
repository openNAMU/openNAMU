function get_post() {
    check = document.getElementById('invert');
    if(check.checked === true) {
        document.cookie = 'main_css_darkmode=1;';
    } else {
        document.cookie = 'main_css_darkmode=0;';
    }

    history.go(0);
}

function main_load() {
    var head_data = document.querySelector('head');
    if(
        cookies.match(regex_data('main_css_darkmode')) &&
        cookies.match(regex_data('main_css_darkmode'))[1] === '1'
    ) {
        head_data.innerHTML += '' +
            '<link rel="stylesheet" href="/views/tenshi/css/dark.css?ver=8">' +
        '';
    }
}

function regex_data(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

var cookies = document.cookie;

function skin_set() {
    if(window.location.pathname === '/skin_set') {
        var set_language = {
            "en-US" : {
                "save" : "Save",
                "darkmode" : "Darkmode"
            }, "ko-KR" : {
                "save" : "저장",
                "darkmode" : "다크모드"
            }
        }

        var language = cookies.match(regex_data('language'))[1];
        var user_language = cookies.match(regex_data('user_language'))[1];
        if(user_language in set_language) {
            language = user_language;
        }

        if(!language in set_language) {
            language = "en-US";
        }

        var set_data = {};

        if(
            cookies.match(regex_data('main_css_darkmode')) &&
            cookies.match(regex_data('main_css_darkmode'))[1] === '1'
        ) {
            set_data["invert"] = "checked";
        }

        document.getElementById("main_skin_set").innerHTML = ' \
            <input ' + set_data["invert"] + ' type="checkbox" id="invert" name="invert" value="invert"> ' + set_language[language]['darkmode'] + ' \
            <hr class="main_hr"> \
            <button onclick="get_post();">' + set_language[language]['save'] + '</button> \
        ';
    }
}