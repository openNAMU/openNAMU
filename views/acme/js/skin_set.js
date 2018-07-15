// 쿠키 생성
function setCookie(name, value, expiredays) {
    var cookie = name + "=" + escape(value) + "; path=/;"
    if (typeof expiredays != 'undefined') {
        var todayDate = new Date();
        todayDate.setDate(todayDate.getDate() + expiredays);
        cookie += "expires=" + todayDate.toGMTString() + ";"
    }
    document.cookie = cookie;
}
 
// 쿠키 획득
function getCookie(name) {
    name += "=";
    var cookie = document.cookie;
    var startIdx = cookie.indexOf(name);
    if (startIdx != -1) {
        startIdx += name.length;
        var endIdx = cookie.indexOf(";", startIdx);
        if (endIdx == -1) {
            endIdx = cookie.length;
            return unescape(cookie.substring(startIdx, endIdx));
        }
    }
    return null;
}
 
// 쿠키 삭제
function deleteCookie(name) {
    setCookie(name, "", -1);
}

// http://vip00112.tistory.com/33

function get_post() {
    console.log("test");

    check = document.getElementById('dark');
    if(check.checked == true) {
        setCookie("set_dark", "1");
        console.log("check");
    } else {
        deleteCookie("set_dark");
        console.log("delete");
    }

    console.log(getCookie("set_dark"));
}

if(getCookie("set_dark") != "1") {
    document.getElementById('set_dark').disabled = true;
}

$(document).ready(function() {
    if(window.location.pathname == "/skin_set") {
        title = document.getElementById("fix_title");
        data = document.getElementById("fix_data");
        
        get_title = "Skin Setting";

        set_data = {};
        set_data["dark"] = "";        

        if(getCookie("set_dark") == "1") {
            set_data["dark"] = "checked";
        }

        get_data =  ' \
                        <input ' + set_data["dark"] + ' type="checkbox" id="dark" name="dark" value="dark"> Dark Mode \
                        <hr> \
                        <button onclick="get_post(); window.location.reload(true);">Save</button> \
                    ';

        document.title = document.title.replace(/.*(\- .*)$/, get_title + " $1");

        title.innerHTML = get_title;
        data.innerHTML = get_data;
    }


});