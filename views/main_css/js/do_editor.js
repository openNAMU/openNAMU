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

function do_not_out() {
    window.addEventListener('DOMContentLoaded', function() {
        window.onbeforeunload = function() {
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
    if (e.clipboardData && e.clipboardData.items) {
        const items = e.clipboardData.items;
        let haveImageInClipboard = false;
        const formData = new FormData();
        for (let i = 0; i < items.length; i++) {
            if (items[i].type.indexOf("image") !== -1) {
                const file = items[i].getAsFile();
                const customName = prompt("파일 이름을 설정해주세요. (확장자는 생략)");
                
                if (!customName) {
                    return alert("취소되었습니다.");
                }
                
                const customFile = new File([file], customName + ".png", { type: file.type });
                formData.append("f_data[]", customFile);
                haveImageInClipboard = true;
                e.preventDefault();
                
                break;
            }
        }
        if (!haveImageInClipboard) {
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
                    '클립보드의 이미지를 성공적으로 업로드했습니다. 아래 텍스트로 본문에 삽입할 수 있습니다. ' +
                    '[[' + decodeURIComponent(url.replace(/.*\/w\/file/, "file")) + ']]'
                );
            } else {
                console.error("[ERROR] PasteUpload Fail :", res.statusText);
                if(res.status === 400) {
                    alert("클립보드의 이미지를 업로드하는데 실패했습니다. 파일 이름 중복일 수 있습니다.");
                } else if(res.status === 401) {
                    alert("클립보드의 이미지를 업로드하는데 실패했습니다. 권한 부족일 수 있습니다.");    
                } else {
                    alert("클립보드의 이미지를 업로드하는데 실패했습니다.");        
                }
            }
        }).catch((err) => {
            console.error("[ERROR] PasteUpload Fail :", JSON.stringify(err), err);
            alert("클립보드의 이미지를 업로드하는데 실패했습니다. 서버가 응답하지 않습니다.");
        });
    }
}