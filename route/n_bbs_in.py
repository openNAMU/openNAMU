from .tool.func import *

from .go_api_bbs import api_bbs

async def bbs_in(bbs_num = 1, page = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        bbs_num_str = str(bbs_num)

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect(conn, '/bbs/main')
    
        bbs_name = db_data[0][0]

        return easy_minify(flask.render_template(await skin_check(),
            imp = [bbs_name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + await get_lang('bbs') + ') (' + str(page) + ')', 0])],
            data = '' + \
                '<div id="opennamu_bbs_in"></div>' + \
                '<script defer src="/views/main_css/js/route/bbs_in.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_bbs_in(); });</script>' + \
            '',
            menu = [['bbs/main', await get_lang('return')], ['bbs/edit/' + bbs_num_str, await get_lang('add')], ['bbs/set/' + bbs_num_str, await get_lang('bbs_set')]]
        ))

'''
async def bbs_in(bbs_num = 1, page = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        bbs_num_str = str(bbs_num)

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect(conn, '/bbs/main')
    
        bbs_name = db_data[0][0]

    data = await api_bbs(bbs_num, page)

    data_html = ''
    for for_a in range(len(data)):
        data_html += '<div class="opennamu_recent_change">'
        data_html += '<a href="/bbs/w/' + data[for_a]['set_id'] + '/' + data[for_a]['set_code'] + '">' + html.escape(data[for_a]['title']) + '</a>'
        data_html += '<div style="float: right;">'
        data_html += '<span id="opennamu_bbs_comment_' + str(for_a) + '"></span>'
        data_html += data[for_a]['user_id_render'] + ' | '

        if data[for_a]['pinned'] == '1':
            data_html += '<span style="color: red;">' + data[for_a]['date'] + '</span>'
        else:
            data_html += data[for_a]['date']

            data_html += '</div>'
            data_html += '<div style="clear: both;"></div>'

            data_html += '</div>'
            data_html += '<hr class="main_hr">'

        get_next_page_bottom()
'''