from .tool.func import *

def bbs_w_tool(bbs_num = '', post_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        data = ''
        
        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)
        
        data += '''
            <h2>''' + load_lang('tool') + '''</h2>
            <ul class="opennamu_ul">
                <li><a href="/bbs/raw/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '">' + load_lang('raw') + '''</a></li>
            </ul>
        '''

        if admin_check() == 1:
            curs.execute(db_change('select set_data from bbs_data where set_code = ? and set_id = ? and set_name = "pinned"'), [post_num_str, bbs_num_str])
            pinned = load_lang('pinned') if not curs.fetchall() else load_lang('pinned_release')

            data += '''
                <h3>''' + load_lang('admin') + '''</h3>
                <ul class="opennamu_ul">
                    <!-- <li><a href="/bbs/blind/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '">' + load_lang('hide') + '''</a></li> -->
                    <li><a href="/bbs/pinned/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '">' + pinned + '''</a></li>
                </ul>
            '''

            data += '''
                <h3>''' + load_lang('owner') + '''</h2>
                <ul class="opennamu_ul">
                    <li><a href="/bbs/delete/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '">' + load_lang('delete') + '''</a></li>
                </ul>
            '''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('bbs_post_tool'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = data,
            menu = [['bbs/w/' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str), load_lang('return')]]
        ))