function opennamu_file_preview() {
    const input = document.getElementById('file_input');
    const preview = document.getElementById('preview');
    if (!input || !preview) return;

    input.addEventListener('change', () => {
        preview.innerHTML = '';

        Array.from(input.files).forEach(file => {
            const mime = file.type;

            if (mime.startsWith('image/')) {
                const img = document.createElement('img');

                img.src = URL.createObjectURL(file);
                img.onload = () => URL.revokeObjectURL(img.src);
                img.style.maxWidth = '100%';

                preview.appendChild(img);
            } else if (mime === 'application/pdf') {
                const embed = document.createElement('embed');

                embed.src = URL.createObjectURL(file);
                embed.type = 'application/pdf';
                embed.style.width = '100px';
                embed.style.height = '120px';
                embed.style.marginRight = '0.5em';

                preview.appendChild(embed);
            } else if (mime.startsWith('text/') || mime === '') {
                const reader = new FileReader();

                reader.onload = () => {
                    const pre = document.createElement('pre');

                    pre.textContent = reader.result.slice(0, 200) + '...';
                    pre.style.border = '1px solid #ccc';
                    pre.style.padding = '0.5em';
                    pre.style.maxWidth = '200px';
                    pre.style.whiteSpace = 'pre-wrap';

                    preview.appendChild(pre);
                };

                reader.readAsText(file);
            }

            const p = document.createElement('p');

            p.textContent = `파일명: ${file.name}, 크기: ${Math.round(file.size / 1024)} KB`;

            preview.appendChild(p);
        });
    });
}