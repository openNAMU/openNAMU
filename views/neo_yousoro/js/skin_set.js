function get_post() {
    check = document.getElementById('strike');
    if(check.checked === true) {
        document.cookie = 'del_strike=true;';
    } else {
        document.cookie = 'del_strike=false;';
    }

    check = document.getElementById('include');
    if(check.checked === true) {
        document.cookie = 'include_link=true;';
    } else {
        document.cookie = 'include_link=false;';
    }

    window.location.reload(true);
}

function regex_data(data) {
    r_data = new RegExp('(?:^|; )' + data + '=([^;]*)')

    return r_data;
}

cookies = document.cookie;

function main_load() {
    head_data = document.querySelector('head');
    if(
        cookies.match(regex_data('del_strike')) &&
        cookies.match(regex_data('del_strike'))[1] === 'true'
    ) {
        head_data.innerHTML += '<style>s { display: none; }</style>';
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
        document.getElementById("main_top").innerHTML = '<h1>Skin setting</h1>';
        document.title = document.title.replace(/.*(\- .*)$/, "Skin setting $1");
        
        data = document.getElementById("main_data");
        set_data = {};

        if(
            cookies.match(regex_data('del_strike')) &&
            cookies.match(regex_data('del_strike'))[1] === 'true'
        ) {
            set_data["strike"] = "checked";
        } 
        
        if(
            cookies.match(regex_data('include_link')) &&
            cookies.match(regex_data('include_link'))[1] === 'true'
        ) {
            set_data["include"] = "checked";
        }

        data.innerHTML = ' \
            <input ' + set_data["strike"] + ' type="checkbox" id="strike" name="strike" value="strike"> Remove strikethrough \
            <hr class="main_hr"> \
            <input ' + set_data["include"] + ' type="checkbox" id="include" name="include" value="include"> Using include link \
            <hr class="main_hr"> \
            <button onclick="get_post();">Save</button> \
        ';
    }
}