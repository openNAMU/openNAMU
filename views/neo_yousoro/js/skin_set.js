function get_post() {
    check = document.getElementById('strike');
    if(check.value === 'normal') {
        document.cookie = 'del_strike=0;';
    } else if(check.value === 'change') {
        document.cookie = 'del_strike=1;';
    } else {
        document.cookie = 'del_strike=2;';
    }

    check = document.getElementById('bold');
    if(check.value === 'normal') {
        document.cookie = 'del_bold=0;';
    } else if(check.value === 'change') {
        document.cookie = 'del_bold=1;';
    } else {
        document.cookie = 'del_bold=2;';
    }

    check = document.getElementById('include');
    if(check.checked === true) {
        document.cookie = 'include_link=true;';
    } else {
        document.cookie = 'include_link=false;';
    }

    history.go(0);
}

function regex_data(data) {
    r_data = new RegExp('(?:^|; )' + data + '=([^;]*)')

    return r_data;
}

cookies = document.cookie;

function main_load() {
    head_data = document.querySelector('head');
    if(cookies.match(regex_data('del_strike'))) {
        if(cookies.match(regex_data('del_strike'))[1] === '1') {
            head_data.innerHTML += '<style>s { text-decoration: none; } s:hover { background-color: transparent; }</style>';
        } else if(cookies.match(regex_data('del_strike'))[1] === '2') {
            head_data.innerHTML += '<style>s { display: none; }</style>';
        }
    }

    if(cookies.match(regex_data('del_bold'))) {
        if(cookies.match(regex_data('del_bold'))[1] === '1') {
            head_data.innerHTML += '<style>b { font-weight: normal; }</style>';
        } else if(cookies.match(regex_data('del_bold'))[1] === '2') {
            head_data.innerHTML += '<style>b { display: none; }</style>';
        }
    }

    if(
        cookies.match(regex_data('include_link')) &&
        cookies.match(regex_data('include_link'))[1] === 'true'
    ) {
        head_data.innerHTML += '<style>#include_link { display: inline; }</style>';
    }
}

main_load();

window.onload = function () {
    if(window.location.pathname === '/skin_set') {
        set_language = {
            "en-US" : {
                "default" : "Default",
                "change_to_noraml" : "Change to normal text",
                "delete" : "Delete",
                "skin_setting" : "Skin settings",
                "include_link" : "Using include link",
                "save" : "Save",
                "strike" : "Strike",
                "bold" : "Bold",
                "other" : "Other"
            }, "ko-KR" : {
                "default" : "기본값",
                "change_to_noraml" : "일반 텍스트로 변경",
                "delete" : "삭제",
                "skin_setting" : "스킨 설정",
                "include_link" : "틀 링크 사용",
                "save" : "저장",
                "strike" : "취소선",
                "bold" : "볼드체",
                "other" : "기타"
            }
        }

        language = cookies.match(regex_data('language'))[1];
        user_language = cookies.match(regex_data('user_language'))[1];
        if(user_language in set_language) {
            language = user_language;
        }

        document.getElementById("main_top").innerHTML = '<h1>' + set_language[language]['skin_setting'] + '</h1>';
        document.title = document.title.replace(/.*(\- .*)$/, set_language[language]['skin_setting'] + " $1");
        
        data = document.getElementById("main_data");
        set_data = {};

        if(cookies.match(regex_data('del_strike'))) {
            if(cookies.match(regex_data('del_strike'))[1] === '0') {
                set_data["strike"] = ' \
                    <option value="normal">' + set_language[language]['default'] + '</option> \
                    <option value="change">' + set_language[language]['change_to_normal'] + '</option> \
                    <option value="delete">' + set_language[language]['delete'] + '</option> \
                ';
            } else if(cookies.match(regex_data('del_strike'))[1] === '1') {
                set_data["strike"] = ' \
                    <option value="change">' + set_language[language]['change_to_normal'] + '</option> \
                    <option value="normal">' + set_language[language]['default'] + '</option> \
                    <option value="delete">' + set_language[language]['delete'] + '</option> \
                ';
            } else {
                set_data["strike"] = ' \
                    <option value="delete">' + set_language[language]['delete'] + '</option> \
                    <option value="normal">' + set_language[language]['default'] + '</option> \
                    <option value="change">' + set_language[language]['change_to_normal'] + '</option> \
                ';
            }
        } else {
            set_data["strike"] = ' \
                <option value="normal">' + set_language[language]['default'] + '</option> \
                <option value="change">' + set_language[language]['change_to_normal'] + '</option> \
                <option value="delete">' + set_language[language]['delete'] + '</option> \
            ';
        }

        if(cookies.match(regex_data('del_bold'))) {
            if(cookies.match(regex_data('del_bold'))[1] === '0') {
                set_data["bold"] = ' \
                    <option value="normal">' + set_language[language]['default'] + '</option> \
                    <option value="change">' + set_language[language]['change_to_normal'] + '</option> \
                    <option value="delete">' + set_language[language]['delete'] + '</option> \
                ';
            } else if(cookies.match(regex_data('del_bold'))[1] === '1') {
                set_data["bold"] = ' \
                    <option value="change">' + set_language[language]['change_to_normal'] + '</option> \
                    <option value="normal">' + set_language[language]['default'] + '</option> \
                    <option value="delete">' + set_language[language]['delete'] + '</option> \
                ';
            } else {
                set_data["bold"] = ' \
                    <option value="delete">' + set_language[language]['delete'] + '</option> \
                    <option value="normal">' + set_language[language]['default'] + '</option> \
                    <option value="change">' + set_language[language]['change_to_normal'] + '</option> \
                ';
            }
        } else {
            set_data["bold"] = ' \
                <option value="normal">' + set_language[language]['default'] + '</option> \
                <option value="change">' + set_language[language]['change_to_normal'] + '</option> \
                <option value="delete">' + set_language[language]['delete'] + '</option> \
            ';
        }
        
        if(
            cookies.match(regex_data('include_link')) &&
            cookies.match(regex_data('include_link'))[1] === 'true'
        ) {
            set_data["include"] = "checked";
        }

        data.innerHTML = ' \
            <h2>' + set_language[language]['strike'] + '</h2> \
            <hr class="main_hr"> \
            <select id="strike" name="strike"> \
                ' + set_data["strike"] + ' \
            </select> \
            <h2>' + set_language[language]['bold'] + '</h2> \
            <select id="bold" name="bold"> \
                ' + set_data["bold"] + ' \
            </select> \
            <h2>' + set_language[language]['other'] + '</h2> \
            <input ' + set_data["include"] + ' type="checkbox" id="include" name="include" value="include"> ' + set_language[language]['include_link'] + ' \
            <hr class="main_hr"> \
            <button onclick="get_post();">' + set_language[language]['save'] + '</button> \
        ';
    }
}