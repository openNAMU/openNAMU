"use strict";

class opennamu_render_wiki {
    constructor(
        render_part_id_add
    ) {
        this.render_part_id_add = render_part_id_add;
    }
    
    do_part_image() {
        
    }
    
    do_part_link() {
        let render_part_id_add = this.render_part_id_add;
        
        let link_list = {};
        let link_list_sub = [];
        for(
            let for_a = 0;
            document.getElementsByClassName(this.render_part_id_add + 'opennamu_link')[for_a];
            for_a++
        ) {
            let link_data = document.getElementsByClassName(this.render_part_id_add + 'opennamu_link')[for_a];
            
            link_list_sub.push(link_data.title);
            
            if(!link_list[link_data.title]) {
                link_list[link_data.title] = [for_a];
            } else {
                link_list[link_data.title].push(for_a);
            }
        }
        
        let data_form = new FormData();
        data_form.append('title_list', JSON.stringify(link_list_sub));
        
        let xhr = new XMLHttpRequest();
        xhr.open("POST", "/api/w/test/doc_tool/exist");
        xhr.send(data_form);
        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                let xhr_data = JSON.parse(this.responseText);
                for(let for_a in link_list) {
                    if(!xhr_data[for_a]) {
                        for(let for_b in link_list[for_a]) {
                            document.getElementsByClassName(render_part_id_add + 'opennamu_link')[link_list[for_a][for_b]].id = "not_thing";
                        }
                    }
                }
            }
        }
    }
    
    do_part_toc() {
        if(document.getElementById('opennamu_TOC_content_1')) {
            for(
                let for_a = 1;
                document.getElementById('opennamu_heading_content_' + String(for_a));
                for_a++
            ) {
                let heading_data = document.getElementById('opennamu_heading_content_' + String(for_a));
                document.getElementById('opennamu_TOC_content_' + String(for_a)).innerHTML = heading_data.innerText;

                document.getElementById('opennamu_heading_content_' + String(for_a)).id = heading_data.innerText;
            }

            let toc_data_all = document.getElementsByClassName('opennamu_TOC');
            let toc_data = '';
            for(
                let for_a = 0;
                for_a < toc_data_all.length;
                for_a++
            ) {
                if(toc_data === '') {
                    toc_data = toc_data_all[0].innerHTML;
                }

                document.getElementsByClassName('opennamu_TOC')[for_a].innerHTML = toc_data;
            }
        }
    }
    
    do_main() {
        this.do_part_link();
        this.do_part_toc();
    }
}