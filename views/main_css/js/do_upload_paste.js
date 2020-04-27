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
                alert("클립보드의 이미지를 업로드하는데 실패했습니다. 파일 이름 중복일 수 있습니다.");
            }
        }).catch((err) => {
            console.error("[ERROR] PasteUpload Fail :", JSON.stringify(err), err);
            alert("클립보드의 이미지를 업로드하는데 실패했습니다. 파일 이름 중복일 수 있습니다.");
        });
    }
}
