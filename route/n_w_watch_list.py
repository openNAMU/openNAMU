from .tool.func import *

from .go_api_w_watch_list import api_w_watch_list

def w_watch_list(db_set, name, num = 1, do_type = 'watch_list'):
    with get_db_connect() as conn:
        data = '<a href="/doc_watch_list/1/' + url_pas(name) + '">(' + load_lang('watchlist') + ')</a> <a href="/doc_star_doc/1/' + url_pas(name) + '">(' + load_lang('star_doc') + ')</a>'
        data += '<ul class="opennamu_ul">'
        
        list_data = json.loads(api_w_watch_list(db_set, name, do_type, num).data)
        data += ''.join(['<li><span class="opennamu_render_ip">' + for_a + '</span></li>' for for_a in list_data])

        data += '</ul>'
        data += '<script>opennamu_do_ip_render();</script>'

        data += get_next_page_bottom('/doc_' + do_type + '/{}/' + url_pas(name), num, list_data)

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang(do_type if do_type == 'star_doc' else 'watchlist') + ')', 0])],
            data = data,
            menu = [['w/' + url_pas(name), load_lang('return')]]
        ))