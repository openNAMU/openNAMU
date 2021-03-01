function do_insert_data(name, data) {
    // https://stackoverflow.com/questions/11076975/insert-text-into-textarea-at-cursor-position-javascript
    if(document.selection) {
        document.getElementById(name).focus();

        var sel = document.selection.createRange();
        sel.text = data;
    } else if(document.getElementById(name).selectionStart || document.getElementById(name).selectionStart == '0') {
        var startPos = document.getElementById(name).selectionStart;
        var endPos = document.getElementById(name).selectionEnd;
        var myPos = document.getElementById(name).value;

        document.getElementById(name).value = myPos.substring(0, startPos) + data + myPos.substring(endPos, myPos.length);
    } else {
        document.getElementById(name).value += data;
    }
}

function monaco_to_content() {
    try {
        document.getElementById('content').innerHTML = window.editor.getValue();
    } catch(e) {}
}

function do_not_out() {
    window.addEventListener('DOMContentLoaded', function() {
        window.onbeforeunload = function() {
            monaco_to_content();
            
            data = document.getElementById('content').value;
            origin = document.getElementById('origin').value;
            if(data !== origin) {
                return '';
            }
        }
    });
}

function save_stop_exit() {
    window.onbeforeunload = function () {}
}

function do_paste_image() {
    window.addEventListener('DOMContentLoaded', function() {
        if(
            document.cookie.match(main_css_regex_data('main_css_image_paste')) &&
            document.cookie.match(main_css_regex_data('main_css_image_paste'))[1] === '1'
        ) {
            const textarea = document.querySelector("textarea");
            if (textarea) {
                textarea.addEventListener("paste", pasteListener);
            }
        }
    });
}

function pasteListener(e) {
    // find file
    if(e.clipboardData && e.clipboardData.items) {
        const items = e.clipboardData.items;
        let haveImageInClipboard = false;
        const formData = new FormData();
        for(let i = 0; i < items.length; i++) {
            if (items[i].type.indexOf("image") !== -1) {
                const file = items[i].getAsFile();
                const customName = prompt("파일 이름 (확장자 제외)");
                
                if (!customName) {
                    return alert("파일 이름 없음");
                }
                
                const customFile = new File([file], customName + ".png", { type: file.type });
                formData.append("f_data[]", customFile);
                haveImageInClipboard = true;
                e.preventDefault();
                
                break;
            }
        }
        if(!haveImageInClipboard) {
            return;
        }

        // send to server
        fetch("/upload", {
            method: "POST",
            body: formData,
        }).then((res) => {
            if (res.status === 200 || res.status === 201) {
                const url = res.url;
                alert(
                    '업로드 완료 : ' +
                    '[[' + decodeURIComponent(url.replace(/.*\/w\/file/, "file")) + ']]'
                );
            } else {
                console.error("[ERROR] PasteUpload Fail :", res.statusText);
                if(res.status === 400) {
                    alert("파일 이름 중복");
                } else if(res.status === 401) {
                    alert("권한 부족");    
                } else {
                    alert("업로드 실패");        
                }
            }
        }).catch((err) => {
            console.error("오류 내역 :", JSON.stringify(err), err);
            alert("업로드 실패");
        });
    }
}

function load_preview(name) {
    var s_data = new FormData();
    s_data.append('data', document.getElementById('content').value);

    var url = "/api/w/" + name;
    var url_2 = "/api/markup";

    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.send(s_data);

    var xhr_2 = new XMLHttpRequest();
    xhr_2.open("GET", url_2, true);
    xhr_2.send(null);

    xhr.onreadystatechange = function() {
        if(xhr.readyState === 4 && xhr.status === 200) {
            var o_p_data = JSON.parse(xhr.responseText);
            document.getElementById('see_preview').innerHTML = o_p_data['data'];
            eval(o_p_data['js_data'])
        }
    }
}

function load_raw_preview(name_1, name_2) {
    document.getElementById(name_2).innerHTML = document.getElementById(name_1).value;
}